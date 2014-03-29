package main

import (
	"github.com/streadway/amqp"
	"log"
	"strings"
)

var (
	amqpReceiveUri       string
	amqpReceiveQueueName = "logs"
	amqpReceiveExchange  = "logs"
	amqpMatchedSendUri   string
	amqpSendQueueName    = "gomatch"
)

// parseAmqpConfigFile fills up all the necessary variables from a file.
// The config file can contain single line # comments.
func parseAmqpConfigFile(filePath string) {
	dataMap := make(map[string]string)
	inputReader := openFile(filePath)
	for {
		configLine, eof := readLine(inputReader)
		if len(configLine) > 0 && configLine[0] != '#' {
			configLineWithoutSpaces := strings.Replace(configLine, " ", "", -1)
			configData := strings.Split(configLineWithoutSpaces, "=")
			if len(configData) == 2 {
				dataMap[configData[0]] = configData[1]
			} else {
				log.Println("invalid config line: \"", configLine, "\" (will be ignored)")
			}
		}
		if eof {
			break
		}
	}
	// check for missing statements in config file
	if dataMap["amqp.receive.uri"] == "" {
		log.Fatal("missing amqp.receive.uri in AMQP config file")
	}
	if dataMap["amqp.matched.send.uri"] == "" {
		log.Fatal("missing amqp.matched.send.uri in AMQP config file")
	}
	// fill global variables from dataMap
	amqpReceiveUri = dataMap["amqp.receive.uri"]
	amqpMatchedSendUri = dataMap["amqp.matched.send.uri"]
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

// send sends a single message using the given queue and a channel.
func send(msg string, ch *amqp.Channel, q amqp.Queue) {
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}
}
