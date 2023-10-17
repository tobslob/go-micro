package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tobslob/listener/event"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	// connect to rabbitmq server
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	defer rabbitConn.Close()
	// start listening for messages
	log.Println("Listening for and consuming Rabbitmq messages...")

	// connect to consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}
	// watch the queue and consume messages
	err = consumer.Listen([]string{"log.INFO", "log.WARN", "log.ERROR"})
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ connection not established")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Print("Backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
