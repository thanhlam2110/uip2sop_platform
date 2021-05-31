package integratedkafka

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"bitbucket.org/cloud-platform/uip2sop_platform/transport"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type MessageData struct {
	Source      string
	Destination string
	Message     string
}

func HandlerConsumer() {
	fmt.Println("-------------------Init Consumer---------------------")
	//config.ReadConfig()
	//url := viper.GetString(`kafka.url`)
	url := "13.212.48.40"
	//groupID := viper.GetString(`kafka.groupID`)
	groupID := "MQTT"
	//autoOffsetReset := viper.GetString(`kafka.autoOffsetReset`)
	autoOffsetReset := "earliest"
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": url,
		"group.id":          groupID,
		"auto.offset.reset": autoOffsetReset,
	})

	if err != nil {
		panic(err)
	}
	c.SubscribeTopics([]string{"95ce1a32-2136-417e-85b4-46b432f1c9ad"}, nil)
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			//fmt.Printf("Nhận Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			ReveiveMessageFromKafka(string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
	c.Close()
}
func HandlerConsumer2() {
	fmt.Println("-------------------Init Consumer---------------------")
	//config.ReadConfig()
	//url := viper.GetString(`kafka.url`)
	//groupID := viper.GetString(`kafka.groupID`)
	//autoOffsetReset := viper.GetString(`kafka.autoOffsetReset`)
	url := "13.212.48.40"
	groupID := "MQTT"
	autoOffsetReset := "earliest"
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": url,
		"group.id":          groupID,
		"auto.offset.reset": autoOffsetReset,
	})

	if err != nil {
		panic(err)
	}
	c.SubscribeTopics([]string{"95ce1a32-2136-417e-85b4-46b432f1c9ad"}, nil)
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			//fmt.Printf("Nhận Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			fmt.Println(string(msg.Value))
			ReveiveMessageFromKafka(string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
	c.Close()
}
func ReveiveMessageFromKafka(message string) {
	fmt.Println("Nhận message " + message)
	//parse string to json
	/*var result map[string]interface{}
	json.Unmarshal([]byte(message), &result)
	destination := (result["destination"]).(string)
	data := (result["message"]).(string)
	PushMessageReceiveFromKafka2(data, destination)*/
	//fmt.Println(message)
	clean := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, message)
	var messageData MessageData
	clean = trimFirstRune(clean)
	Data := []byte(clean)
	err := json.Unmarshal(Data, &messageData)

	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Struct is:", messageData.Message)
	PushMessageReceiveFromKafka2(messageData.Message, messageData.Destination)
}
func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func ReveiveMessageFromKafka2(message string) {
	fmt.Println("Nhận message " + message)
	var result map[string]interface{}
	json.Unmarshal([]byte(message), &result)
	destination := (result["destination"]).(string)
	data := (result["message"]).(string)
	PushMessageReceiveFromKafka2(data, destination)
}
func PushMessageReceiveFromKafka(message, destination string) {
	//khởi tạo connect với id của user đại diện
	client, err := ConnectForUserTo("76579b0e-e899-4fc9-a839-607c371b9585")
	if err != nil {
		panic(err)
	}
	if client != nil {
		client.Publish("public/"+destination, 2, false, string(message))
	}
}
func PushMessageReceiveFromKafka2(data, destination string) {
	//config.ReadConfig()
	//url := viper.GetString(`kafka.url`)
	url := "13.212.48.40"
	//groupID := "MQTT"
	//autoOffsetReset := "earliest"
	KafkaInitPublic(url)
	producerPublic, _ := transport.InitProducer(url)
	fmt.Println("--------")
	fmt.Println(data)
	transport.PublishMessage(producerPublic, "public-nhuy", kafka.PartitionAny, data)
	defer producerPublic.Close()
}

//Connect cho user dai diện
func ConnectForUserTo(id string) (c mqtt.Client, err error) {
	//config.ReadConfig()
	//link := viper.GetString(`mqttbroker.urlpublic`)
	link := "tcp://13.212.194.253:1883"
	opts := mqtt.NewClientOptions().AddBroker(link).SetClientID(id)
	opts.SetKeepAlive(0)
	c = mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return c, nil
}
func KafkaInitPublic(kafkaURL string) {
	//UseKafka = true
	producer, err := transport.InitProducer(kafkaURL)
	if err != nil {
		panic(err)
	}
	kafkaProducerPushPublic = producer

}
