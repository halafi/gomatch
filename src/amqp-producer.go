package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func publish(amqpURI, exchange, producerExchangeType, producerRoutingKey, producerBody string, producerReliable bool) error {

	// This function dials, connects, declares, publishes, and tears down,
	// all in one go. In a real service, you probably want to maintain a
	// long-lived connection as state, and publish against that.

	log.Printf("dialing %q", amqpURI)
	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}
	defer connection.Close()

	log.Printf("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring %q Exchange (%q)", producerExchangeType, exchange)
	if err := channel.ExchangeDeclare(
		exchange,             // name
		producerExchangeType, // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // noWait
		nil,                  // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	// producerReliable publisher confirms require confirm.select support from the
	// connection.
	if producerReliable {
		log.Printf("enabling publishing confirms.")
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}

		ack, nack := channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

		defer confirmOne(ack, nack)
	}

	log.Printf("declared Exchange, publishing %dB producerBody (%q)", len(producerBody), producerBody)
	if err = channel.Publish(
		exchange,           // publish to an exchange
		producerRoutingKey, // routing to 0 or more queues
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(producerBody),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func confirmOne(ack, nack chan uint64) {
	log.Printf("waiting for confirmation of one publishing")

	select {
	case tag := <-ack:
		log.Printf("confirmed delivery with delivery tag: %d", tag)
	case tag := <-nack:
		log.Printf("failed delivery of delivery tag: %d", tag)
	}
}
