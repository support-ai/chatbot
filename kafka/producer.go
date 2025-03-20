package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"chatbot/models"

	"github.com/Shopify/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer() (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	brokers := []string{os.Getenv("KAFKA_BROKERS")}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kafka producer: %v", err)
	}

	return &Producer{
		producer: producer,
		topic:    os.Getenv("KAFKA_TOPIC"),
	}, nil
}

func (p *Producer) LogMessage(chatLog *models.ChatLog) error {
	value, err := json.Marshal(chatLog)
	if err != nil {
		return fmt.Errorf("error marshaling chat log: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("error sending message to Kafka: %v", err)
	}

	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
