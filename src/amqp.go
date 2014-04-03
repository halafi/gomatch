package main

import (
	"github.com/streadway/amqp"
	"log"
	"strings"
)

var (
	amqpReceiveUri           string
	amqpReceiveQueueName     string // logs
	amqpReceiveExchange      string // logs
	amqpReceiveFormat        string // plain or json
	amqpMatchedSendUri       string
	amqpMatchedSendQueueName string // gomatch
)

// parseAmqpConfigFile fills up all the necessary variables from a file.
// The config file can contain single line # comments.
func parseAmqpConfigFile(filePath string) {
	m := make(map[string]string)
	inputReader := openFile(filePath)

	for {
		line, eof := readLine(inputReader)
		configLine := string(line)
		if len(configLine) > 0 && configLine[0] != '#' {
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

	// check for missing statements in config file
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

// send sends a single message using the given queue and a channel.
func send(msg []byte, ch *amqp.Channel, q amqp.Queue) {
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})

	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}
}

// emitLog sends one log message. Testing purposes.
/*func emitLog(msg string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {log.Fatal("emmit", err)}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {log.Fatal("emmit", err)}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",    // name
		"fanout",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {log.Fatal("emmit", err)}

	body := msg
	err = ch.Publish(
		"logs", // exchange
		"",     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:     "text/plain",
			Body:            []byte(body),
		})

	if err != nil {log.Fatal("emmit", err)}
	log.Println("Sent: ",msg)
}*/
