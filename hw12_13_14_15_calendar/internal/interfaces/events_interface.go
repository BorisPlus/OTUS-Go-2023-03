package interfaces

import (
	"hw12_13_14_15_calendar/internal/models"
	calendarrpcapi "hw12_13_14_15_calendar/internal/server/rpc/rpcapi"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Event interface {
	*models.Event | *calendarrpcapi.Event | *amqp.Delivery | amqp.Delivery
}
