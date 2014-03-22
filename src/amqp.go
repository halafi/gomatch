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
	// fill the amqp global variables from dataMap
	amqpReceiveUri = dataMap["amqp.receive.uri"]
	amqpMatchedSendUri = dataMap["amqp.matched.send.uri"]
}

// receiveLogs reads all of the messages.
func receiveLogs() []string {
	logs := make([]string, 0)

	conn, err := amqp.Dial(amqpReceiveUri)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		amqpReceiveQueueName, // name
		false,                // durable
		false,                // delete when usused
		false,                // exclusive
		false,                // noWait
		nil,                  // arguments
	)
	failOnError(err, "Failed to declare a queue")
	err = ch.QueueBind(
		q.Name,              // queue name
		"",                  // routing key
		amqpReceiveExchange, // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	for d := range msgs {
		logs = append(logs, string(d.Body))
		ch.Close()
	}
	return logs

}

// send sends a single string message using RabbitMQ (queue: gomatch).
func send(msg string) {
	conn, err := amqp.Dial(amqpMatchedSendUri)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	defer ch.Close()

	q, err := ch.QueueDeclare(
		amqpSendQueueName, // name
		false,             // durable
		false,             // delete when usused
		false,             // exclusive
		false,             // noWait
		nil,               // arguments
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

// emitLog sends one log message. Testing purposes.
/*func emitLog(msg string) {
	conn, err := amqp.Dial(amqpMatchedSendUri)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		amqpReceiveQueueName,    // name
		"fanout",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // noWait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	body := msg
	err = ch.Publish(
		amqpReceiveExchange, // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:     "text/plain",
			Body:            []byte(body),
		})

	failOnError(err, "Failed to publish a message")
	log.Println("Sent: ",msg)
}*/
