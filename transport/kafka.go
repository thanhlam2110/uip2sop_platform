package transport

import (
	"fmt"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// InitProducer - InitProducer
func InitProducer(kafkaURL string) (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaURL})
	if err != nil {
		return nil, err
	}
	return p, err
}

// CloseProducer - CloseProducer
func CloseProducer(producer *kafka.Producer) {
	fmt.Println("Closed Producer")
	defer producer.Close()
}

// PublishMessage - PublishMessage
func PublishMessage(p *kafka.Producer, topic string, partition int32, value string) error {
	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Value:          []byte(value),
	}, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	p.Flush(10000)
	return nil
}
