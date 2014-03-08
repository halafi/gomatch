package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	consumerUri          string
	consumerExchangeName string
	consumerExchangeType string
	consumerQueue        string
	consumerBindingKey   string
	consumerTag          string
	consumerLifetime     time.Duration
	producerUri          string
	producerExchangeName string
	producerExchangeType string
	producerRoutingKey   string
	producerBody         string
	producerReliable     bool
)

// parseAmqpConfigFile fills up all the necessary variables from a file.
func parseAmqpConfigFile(filePath string) {
	// map for storing data read
	dataMap := make(map[string]string)
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	lines := lineSplit(string(b))
	// fill map
	for j := range lines {
		if len(lines[j]) != 0 && lines[j][0] != '#' {
			data := strings.Split(lines[j], "=")
			dataMap[data[0]] = data[1]
		}
	}
	// fill variables from data in map
	consumerUri = dataMap["amqp.consumer.uri"]
	consumerExchangeName = dataMap["amqp.consumer.exchangeName"]
	consumerExchangeType = dataMap["amqp.consumer.exchangeType"]
	consumerQueue = dataMap["amqp.consumer.queue"]
	consumerBindingKey = dataMap["amqp.consumer.bindingKey"]
	consumerTag = dataMap["amqp.consumer.tag"]
	tmpConsumerLifetime, err := strconv.Atoi(dataMap["amqp.consumer.lifetimesec"])
	if err != nil {
		log.Fatal(err)
	}
	consumerLifetime = time.Duration(tmpConsumerLifetime) * time.Second
	producerUri = dataMap["amqp.producer.uri"]
	producerExchangeName = dataMap["amqp.producer.exchangeName"]
	producerExchangeType = dataMap["amqp.producer.exchangeType"]
	producerRoutingKey = dataMap["amqp.producer.routingKey"]
	producerBody = dataMap["amqp.producer.body"]
	tmpProducerReliable := dataMap["amqp.producer.reliable"]
	if tmpProducerReliable == "true" {
		producerReliable = true
	} else if tmpProducerReliable == "false" {
		producerReliable = false
	} else {
		log.Fatal("invalid AMQP config file amqp.producer.reliable cannot be: ", tmpProducerReliable)
	}
}
