package kafka

import (
	"log"

	cli "gopkg.in/urfave/cli.v2"

	"github.com/Shopify/sarama"
)

var (
	kAsyncProducer sarama.AsyncProducer
	kClient        sarama.Client
	ChatTopic      string
)

func initKafka(c *cli.Context) {
	addrs := c.StringSlice("kafka-brokers")
	ChatTopic = c.String("chat-topic")
	config := sarama.NewConfig()
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = false
	producer, err := sarama.NewAsyncProducer(addrs, config)
	if err != nil {
		log.Fatalln(err)
	}

	kAsyncProducer = producer
	cli, err := sarama.NewClient(addrs, nil)
	if err != nil {
		log.Fatalln(err)
	}
	kClient = cli
}

func Init(c *cli.Context) {
	initKafka(c)
}
func NewConsumer() (sarama.Consumer, error) {
	return sarama.NewConsumerFromClient(kClient)
}
