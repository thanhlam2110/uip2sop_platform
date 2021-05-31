package integratedkafka

import (
	"bitbucket.org/cloud-platform/uip2sop_platform/transport"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var kafkaProducerPushPublic *kafka.Producer

// KafkaInit - KafkaInit
func KafkaInit(kafkaURL string) {
	//UseKafka = true
	producer, err := transport.InitProducer(kafkaURL)
	if err != nil {
		panic(err)
	}
	kafkaProducerPushPublic = producer

}
