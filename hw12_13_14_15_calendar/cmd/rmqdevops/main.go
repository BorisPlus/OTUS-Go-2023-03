package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	config "hw12_13_14_15_calendar/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	configFile string
	withDrop   bool
)

func init() {
	pflag.StringVar(&configFile, "config", "", "Path to configuration file")
	pflag.BoolVarP(&withDrop, "with-drop", "d", false, "Preverious drop exchanges and queues")
}

func main() {
	pflag.Parse()
	if pflag.Arg(0) == "version" {
		fmt.Printf("2023.08.13 v.1")
		return
	}
	if configFile == "" {
		fmt.Println("Please set: '--config=<Path to configuration file>'")
		return
	}
	viper.SetConfigType("yaml")
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	viper.ReadConfig(file)
	cfg := config.NewRMQOpsConfig()
	err = viper.Unmarshal(cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	conn, err := amqp.Dial(cfg.RabbitMQ.DSN)
	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}
	log.Printf("Connect to %q.\n", cfg.RabbitMQ.DSN)

	defer func() {
		log.Println("Connection close.")
		_ = conn.Close()
	}()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel. Error: %s", err)
	}
	log.Println("Channel open.")

	defer func() {
		log.Println("Channel close.")
		_ = channel.Close()
	}()
	fmt.Println(withDrop)
	for _, exch := range cfg.RabbitMQ.Exchanges {
		if withDrop {
			err = channel.ExchangeDelete(exch.Name, false, true)
			if err != nil {
				log.Printf("ERROR. Drop exchange: %q.\n%s\n", exch.Name, err.Error())
			} else {
				log.Printf("OK. Drop exchange: %q.\n", exch.Name)
			}
			for _, binding := range exch.Bindings {
				_, err = channel.QueueDelete(binding.BindQueue, false, false, true)
				if err != nil {
					log.Printf("ERROR. Drop queue: %q.\n%s\n", binding.BindQueue, err.Error())
				} else {
					log.Printf("OK. Drop queue: %q.\n", binding.BindQueue)
				}
			}
		}
		// TODO: почему Passive закрывает channel
		err = channel.ExchangeDeclare(exch.Name, "direct", true, false, false, false, nil)
		if err != nil {
			panic(err)
		}
		log.Printf("Declare exchange: %q.\n", exch.Name)
		for _, binding := range exch.Bindings {
			_, err = channel.QueueDeclare(binding.BindQueue, true, false, false, false, nil)
			if err != nil {
				panic("error declaring the queue: " + binding.BindQueue + " " + err.Error())
			}
			log.Printf("Declare queue: %q.\n", binding.BindQueue)
			err = channel.QueueBind(binding.BindQueue, binding.BindKey, exch.Name, false, nil)
			if err != nil {
				fmt.Printf("Error ind queue: %q to exchange %q with routing key %q.\n",
					binding.BindQueue,
					exch.Name,
					binding.BindKey,
				)
				panic("error: " + err.Error())
			}
			log.Printf("Bind queue: %q to exchange %q with routing key %q.\n",
				binding.BindQueue,
				exch.Name,
				binding.BindKey,
			)
		}
	}
	log.Println("RabbitMQ (re-)configure done.")
}
