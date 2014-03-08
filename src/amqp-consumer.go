package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Consumer struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	consumerTag string
	done        chan error
}

func NewConsumer(amqpURI, exchange, consumerExchangeType, queueName, key, ctag string) (*Consumer, error) {
	c := &Consumer{
		conn:        nil,
		channel:     nil,
		consumerTag: ctag,
		done:        make(chan error),
	}

	var err error

	log.Printf("dialing %q", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	go func() {
		fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring Exchange (%q)", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,             // name of the exchange
		consumerExchangeType, // type
		true,                 // durable
		false,                // delete when complete
		false,                // internal
		false,                // noWait
		nil,                  // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	log.Printf("declared Exchange, declaring consumerQueue %q", queueName)
	consumerQueue, err := c.channel.QueueDeclare(
		queueName, // name of the consumerQueue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("consumerQueue Declare: %s", err)
	}

	log.Printf("declared consumerQueue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		consumerQueue.Name, consumerQueue.Messages, consumerQueue.Consumers, key)

	if err = c.channel.QueueBind(
		consumerQueue.Name, // name of the consumerQueue
		key,                // consumerBindingKey
		exchange,           // sourceExchange
		false,              // noWait
		nil,                // arguments
	); err != nil {
		return nil, fmt.Errorf("consumerQueue Bind: %s", err)
	}

	log.Printf("consumerQueue bound to Exchange, starting Consume (consumer consumerTag %q)", c.consumerTag)
	deliveries, err := c.channel.Consume(
		consumerQueue.Name, // name
		c.consumerTag,      // consumerTag,
		false,              // noAck
		false,              // exclusive
		false,              // noLocal
		false,              // noWait
		nil,                // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("consumerQueue Consume: %s", err)
	}

	go handle(deliveries, c.done)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.consumerTag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}
