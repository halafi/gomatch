package main

import (
	"github.com/streadway/amqp"
	"log"
	"strings"
)

var (
	// AMQP configuration variables.
	amqpReceiveUri           string
	amqpReceiveQueueName     string
	amqpReceiveExchange      string
	amqpMatchedSendUri       string
	amqpMatchedSendQueueName string
	// AMQP input, either plain or json
	amqpReceiveFormat string
)

// parseAmqpConfigFile fills up all the necessary variables from a file.
func parseAmqpConfigFile(filePath string) {
	m := make(map[string]string)
	inputReader := openFile(filePath)

	for {
		line, eof := readLine(inputReader)
		configLine := string(line)
		if len(configLine) > 0 && configLine[0] != '#' { // ignore empty lines and comments
			configLineWithoutSpaces := strings.Replace(configLine, " ", "", -1)
			configData := strings.Split(configLineWithoutSpaces, "=")
			if len(configData) == 2 {
				m[configData[0]] = configData[1]
			} else {
				log.Println("invalid config line: \"", configLine, "\" (will be ignored)")
			}
		}
		if eof {
			break
		}
	}

	// check for missing statements
	if m["amqp.receive.uri"] == "" || m["amqp.receive.format"] == "" ||
		m["amqp.receive.queue"] == "" || m["amqp.receive.exchange"] == "" ||
		m["amqp.matched.send.uri"] == "" || m["amqp.matched.send.queue"] == "" {
		log.Fatal("missing argument in AMQP config file")
	}

	// fill global variables from dataMap
	amqpReceiveUri = m["amqp.receive.uri"]
	amqpReceiveFormat = m["amqp.receive.format"]
	amqpReceiveQueueName = m["amqp.receive.queue"]
	amqpReceiveExchange = m["amqp.receive.exchange"]
	amqpMatchedSendUri = m["amqp.matched.send.uri"]
	amqpMatchedSendQueueName = m["amqp.matched.send.queue"]
}

// openConnection opens up a RabbitMQ connection.
func openConnection(amqpUri string) *amqp.Connection {
	conn, err := amqp.Dial(amqpUri)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	return conn
}

// openChannel opens up a RabbitMQ channel.
func openChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	return ch
}

// declareQueue declares a queue - a buffer for messages.
func declareQueue(name string, ch *amqp.Channel) amqp.Queue {
	q, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)

	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	return q
}

// bindReceiveQueue binds the amqpReceiveExchange name to a queue.
func bindReceiveQueue(ch *amqp.Channel, q amqp.Queue) {
	err := ch.QueueBind(
		q.Name,              // queue name
		"",                  // routing key
		amqpReceiveExchange, // exchange
		false,
		nil)

	if err != nil {
		log.Fatalf("Failed to bind a queue: %s", err)
	}
}

// send sends a message using the given queue and a channel.
func send(msg []byte, rk string, ch *amqp.Channel, q amqp.Queue) {
	err := ch.Publish(
		"",    // exchange
		rk,    // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})

	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}
}
