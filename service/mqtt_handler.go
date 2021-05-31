package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"bitbucket.org/cloud-platform/uip2sop_platform/config"
	"bitbucket.org/cloud-platform/uip2sop_platform/integratedkafka"
	"bitbucket.org/cloud-platform/uip2sop_platform/model"
	"bitbucket.org/cloud-platform/uip2sop_platform/repository"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"

	"bitbucket.org/cloud-platform/uip2sop_platform/transport"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func PushMessage(c echo.Context) error {
	pushMessage := new(model.PushMessage)
	err := c.Bind(pushMessage)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}
	//parse token ---> get user status
	authResponse, err := BasicAuth(pushMessage.Token)

	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "Connection refused", "data": map[string]interface{}{"info": nil}})
	}
	//fmt.Println(authResponse)
	var result map[string]interface{}
	json.Unmarshal([]byte(authResponse), &result)
	if result["error"] != nil {
		if result["message"] != nil {
			return c.JSON(200, map[string]interface{}{"code": "7", "message": (result["message"]).(string), "data": map[string]interface{}{"info": nil}})
		}
		return c.JSON(200, map[string]interface{}{"code": "5", "message": (result["error"]).(string), "data": map[string]interface{}{"info": nil}})
	}
	attributes := result["attributes"].(map[string]interface{})
	if attributes["userstatus"] != "ACTIVE" {
		return c.JSON(200, map[string]interface{}{"code": "8", "message": "USER DISABLED", "data": map[string]interface{}{"info": nil}})
	}
	//check map thing to chanel
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}
	profileRepository := repository.NewProfileRepositoryMongo(db, "things_map_chanels")
	_, err = profileRepository.FindMapThingChanel(pushMessage.Thingid, pushMessage.Chanelid)
	if err == mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "2", "message": "ThingId haven't mapped to Chanel", "data": map[string]interface{}{"info": nil}})
	}
	//
	//connect broker
	connect, err := ConnectMQTT(pushMessage.Thingid, pushMessage.Thingkey)
	/*fmt.Println("---------------------")
	fmt.Println(err)
	fmt.Println("---------------------")*/
	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "0", "message": err.Error(), "data": map[string]interface{}{"info": nil}})
	}
	if connect != nil {
		//connect.Publish("test/topic", 1, false, "Example Payload")
		connect.Publish(pushMessage.Chanelid, 0, false, pushMessage.Message)
		/*if token := connect.Publish(pushMessage.Chanelid, 1, false, pushMessage.Message); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}*/
	}
	session.Close()
	connect.Disconnect(250)
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Push Success", "data": map[string]interface{}{"info": nil}})
}
func ConnectMQTT(thingid, thingkey string) (c mqtt.Client, err error) {
	//config.ReadConfig()
	//link := viper.GetString(`mqttbroker.urllocal`)
	link := "tcp://18.141.164.223:1883"
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		fmt.Println("MongoDB connection refused")
		return nil, errors.New("MongoDB connection refused")
	}
	//check thing_id + thning_key ton tai
	profileRepositorya := repository.NewProfileRepositoryMongo(db, "things")
	_, err = profileRepositorya.FindThing(thingid, thingkey)
	if err == mgo.ErrNotFound {
		fmt.Println("ThingID and ThingKey not exist")
		return nil, errors.New("ThingID and ThingKey not exist")
	}
	//if userstatus == "ACTIVE" {
	opts := mqtt.NewClientOptions().AddBroker(link).SetClientID(thingid)

	c = mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	session.Close()
	return c, nil
	//}
	//return nil, nil
}

//Connect cho user dai diện
func ConnectForUserTo(id string) (c mqtt.Client, err error) {
	//config.ReadConfig()
	//link := viper.GetString(`mqttbroker.urlpublic`)
	link := "tcp://18.141.164.223:1883"
	opts := mqtt.NewClientOptions().AddBroker(link).SetClientID(id)
	//opts.SetUsername("username")
	//opts.SetPassword("password")
	opts.SetKeepAlive(0)
	c = mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		//panic(token.Error())
		return nil, token.Error()
	}
	return c, nil
}

//danh cho user đại diện
func SubscribeMQTT(topic string) {
	//var msg mqtt.Message
	//subscribe with organization1 userid
	client, _ := ConnectForUserTo("76579b0e-e899-4fc9-a839-607c371b9585")
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
		//parse string payload to json
		var result map[string]interface{}
		json.Unmarshal([]byte(string(msg.Payload())), &result)
		//fmt.Println(result["destination"])
		err, organizationUserId := FindRootOfUser((result["destination"]).(string))
		if err != nil {
			panic(err)
		}
		//Khoi tao Kafka
		//config.ReadConfig()
		//url := viper.GetString(`kafka.url`)
		url := "13.212.48.40"
		integratedkafka.KafkaInit(url)
		producer, _ := transport.InitProducer(url)
		//
		fmt.Println(organizationUserId)
		//lay PublichanelID theo organizationUserId
		db, session, err := config.GetMongoDataBase()
		if err != nil {
			panic(err)
		}
		profileRepository := repository.NewProfileRepositoryMongo(db, "public_chanels")
		organizationPublicChanel, err := profileRepository.FindPublicChanelID(organizationUserId)
		//fmt.Println(err)
		//fmt.Println(organizationPublicChanel.Publichanelid)
		if err != mgo.ErrNotFound {

			transport.PublishMessage(producer, organizationPublicChanel.Publichanelid, kafka.PartitionAny, string(msg.Payload()))
		} else {
			transport.PublishMessage(producer, "UNKNOW_TOPIC", kafka.PartitionAny, string(msg.Payload()))
		}
		session.Close()
	})
}

//push message to public broker
func PushMessageToPublic(c echo.Context) error {
	pushMessagePublic := new(model.PushMessagePublic)
	err := c.Bind(pushMessagePublic)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}
	//parse token ---> get user status
	authResponse, err := BasicAuth(pushMessagePublic.Token)

	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "Connection refused", "data": map[string]interface{}{"info": nil}})
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(authResponse), &result)
	if result["error"] != nil {
		if result["message"] != nil {
			return c.JSON(200, map[string]interface{}{"code": "7", "message": (result["message"]).(string), "data": map[string]interface{}{"info": nil}})
		}
		return c.JSON(200, map[string]interface{}{"code": "5", "message": (result["error"]).(string), "data": map[string]interface{}{"info": nil}})
	}
	attributes := result["attributes"].(map[string]interface{})
	if attributes["userstatus"] != "ACTIVE" {
		return c.JSON(200, map[string]interface{}{"code": "8", "message": "USER DISABLED", "data": map[string]interface{}{"info": nil}})
	}
	//connect broker
	connect, err := ConnectForUserTo((attributes["userid"]).(string))
	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "9", "message": err.Error(), "data": map[string]interface{}{"info": nil}})
	}
	dataPublic := new(model.DataPublic)
	dataPublic.Source = pushMessagePublic.Source
	dataPublic.Destination = pushMessagePublic.Destination
	dataPublic.Message = pushMessagePublic.Message
	stringData, err := json.Marshal(dataPublic)
	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "9", "message": err.Error(), "data": map[string]interface{}{"info": nil}})
	}
	if connect != nil {
		//connect.Publish("test/topic", 1, false, "Example Payload")
		connect.Publish("public", 0, false, string(stringData))
		/*if token := connect.Publish("public/", 1, false, string(stringData)); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}*/
	}
	//session.Close()
	connect.Disconnect(250)
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Push Public Broker Success", "data": map[string]interface{}{"info": nil}})
}

//HAM NEW CONFLUENT
func PushMessageToPublicBroker(c echo.Context) error {
	pushMessagePublic := new(model.PushMessagePublic)
	err := c.Bind(pushMessagePublic)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}
	//parse token ---> get user status
	authResponse, err := BasicAuth(pushMessagePublic.Token)

	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "Connection refused", "data": map[string]interface{}{"info": nil}})
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(authResponse), &result)
	if result["error"] != nil {
		if result["message"] != nil {
			return c.JSON(200, map[string]interface{}{"code": "7", "message": (result["message"]).(string), "data": map[string]interface{}{"info": nil}})
		}
		return c.JSON(200, map[string]interface{}{"code": "5", "message": (result["error"]).(string), "data": map[string]interface{}{"info": nil}})
	}
	attributes := result["attributes"].(map[string]interface{})
	if attributes["userstatus"] != "ACTIVE" {
		return c.JSON(200, map[string]interface{}{"code": "8", "message": "USER DISABLED", "data": map[string]interface{}{"info": nil}})
	}
	//connect broker
	connect, err := ConnectForUserTo((attributes["userid"]).(string))
	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "9", "message": err.Error(), "data": map[string]interface{}{"info": nil}})
	}
	dataPublic := new(model.DataPublic)
	dataPublic.Source = pushMessagePublic.Source
	dataPublic.Destination = pushMessagePublic.Destination
	dataPublic.Message = pushMessagePublic.Message
	stringData, err := json.Marshal(dataPublic)
	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "9", "message": err.Error(), "data": map[string]interface{}{"info": nil}})
	}
	if connect != nil {
		//connect.Publish("test/topic", 1, false, "Example Payload")
		connect.Publish("public/lam", 0, false, string(stringData))
		/*if token := connect.Publish("public/", 1, false, string(stringData)); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}*/
	}
	//session.Close()
	connect.Disconnect(250)
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Push Public Broker Success", "data": map[string]interface{}{"info": nil}})
}

//HAM NEW CONFLUENT
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
