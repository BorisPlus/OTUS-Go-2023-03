package main

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	interfaces "hw12_13_14_15_calendar/internal/interfaces"
	models "hw12_13_14_15_calendar/internal/models"
	rpcClient "hw12_13_14_15_calendar/internal/server/rpc/client"
)

type Sheduler struct {
	sourceRPCDsn             string
	sourceRPCClient          rpcClient.Client
	targetRabbitMQ           models.RabbitMQTarget
	targetRabbitMQConnection *amqp.Connection
	logger                   interfaces.Logger
	loopTimeoutSec           int64
}

func NewSheduler(
	sourceRPCDsn string,
	targetRabbitMQ models.RabbitMQTarget,
	logger interfaces.Logger,
	loopTimeoutSec int64,
) *Sheduler {
	return &Sheduler{
		sourceRPCDsn:   sourceRPCDsn,
		targetRabbitMQ: targetRabbitMQ,
		logger:         logger,
		loopTimeoutSec: loopTimeoutSec,
	}
}

func (s *Sheduler) Stop() error {
	err := s.sourceRPCClient.Close()
	if err != nil {
		s.logger.Warning(err.Error())
	}
	if s.targetRabbitMQConnection.IsClosed() {
		return nil
	}
	s.logger.Info("Sheduler.Stop()")
	err = s.targetRabbitMQConnection.Close()
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	return nil
}

func (s *Sheduler) Start(ctx context.Context) error {
	s.logger.Info("Sheduler.Start()")
	var err error
	s.targetRabbitMQConnection, err = amqp.Dial(s.targetRabbitMQ.DSN)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	defer func() {
		s.Stop()
	}()
	channel, err := s.targetRabbitMQConnection.Channel()
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	defer func() {
		s.logger.Info("Channel close.")
		_ = channel.Close()
	}()
	s.sourceRPCClient.Connect(s.sourceRPCDsn)
	for {
		events, err := s.sourceRPCClient.ListNotSheduledEvents(ctx)
		if err != nil {
			s.logger.Error(err.Error())
			return err
		}
		for _, event := range events {
			s.logger.Info("Sheduler.Put %+v", event)
			err = channel.ExchangeDeclarePassive(s.targetRabbitMQ.ExchangeName, "direct", true, false, false, false, nil)
			if err != nil {
				s.logger.Error(err.Error())
				return err
			}
			data, err := json.Marshal(event)
			if err != nil {
				s.logger.Error(err.Error())
				return err
			}
			message := amqp.Publishing{
				Body: data,
			}
			err = channel.PublishWithContext(ctx, s.targetRabbitMQ.ExchangeName, s.targetRabbitMQ.RoutingKey, false, false, message)
			if err != nil {
				s.logger.Error(err.Error())
				return err
			}
			event.Sheduled = true
			s.sourceRPCClient.UpdateEvent(ctx, event)
		}
		s.logger.Info("Sheduler.Sleep %+v", s.loopTimeoutSec)
		time.Sleep(time.Duration(s.loopTimeoutSec) * time.Second)
	}
}
