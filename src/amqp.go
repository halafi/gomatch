package main

import (
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"strings"
)

var (
	amqpReceiveUri   string
	amqpMatchSendUri string
)

// parseAmqpConfigFile fills up all the necessary variables from a file.
func parseAmqpConfigFile(filePath string) {
	dataMap := make(map[string]string)
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	lines := lineSplit(string(b))
	for j := range lines {
		if len(lines[j]) != 0 && lines[j][0] != '#' {
			lines[j] = strings.Replace(lines[j], " ", "", -1)
			data := strings.Split(lines[j], "=")
			dataMap[data[0]] = data[1]
		}
	}
	// fill global variables from data in map
	if dataMap["amqp.receive.uri"] == "" || dataMap["amqp.match.send.uri"] == "" {
		log.Fatal("missing amqp.receive.uri or amqp.match.send.uri in AMQP config file")
	}
	amqpReceiveUri = dataMap["amqp.receive.uri"]
	amqpMatchSendUri = dataMap["amqp.match.send.uri"]
}

// receiveLogs reads all logs (exchange logs) from AMQP.
func receiveLogs() []string {
	data := ""
	conn, err := amqp.Dial(amqpReceiveUri)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")
	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	for d := range msgs {
		data = data + string(d.Body)
	}

	return lineSplit(data)
}

// Function send sends string to AMQP.
func send(msg string) {
	conn, err := amqp.Dial(amqpMatchSendUri)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"gomatch", // name
		false,     // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := msg
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}

// failOnError logs error and fails the program execution.
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
