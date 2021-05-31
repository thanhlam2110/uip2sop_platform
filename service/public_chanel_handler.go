package service

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/uip2sop_platform/config"
	"bitbucket.org/cloud-platform/uip2sop_platform/integratedkafka"
	"bitbucket.org/cloud-platform/uip2sop_platform/model"
	"bitbucket.org/cloud-platform/uip2sop_platform/repository"
	"bitbucket.org/cloud-platform/uip2sop_platform/transport"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"gopkg.in/mgo.v2"
)

func CreateOrganizationPublicChanel(c echo.Context) error {
	organizationPublicChanelRequest := new(model.OrganizationPublicChanelRequest)
	err := c.Bind(organizationPublicChanelRequest)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}
	//Parse token
	//parse token ---> get user status
	authResponse, err := BasicAuth(organizationPublicChanelRequest.Token)

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
	//Check User status
	if attributes["userstatus"] != "ACTIVE" {
		return c.JSON(200, map[string]interface{}{"code": "8", "message": "USER DISABLED", "data": map[string]interface{}{"info": nil}})
	}
	//
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}
	//check organization exists
	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	_, err = profileRepository.FindUserByUserID((attributes["userid"]).(string))
	if err == mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "2", "message": "Organization is Not-Exists", "data": map[string]interface{}{"info": nil}})
	}
	//check organization created public chanel
	profileRepository2 := repository.NewProfileRepositoryMongo(db, "public_chanels")
	_, err = profileRepository2.FindPublicChanel((attributes["userid"]).(string))
	if err != mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "2", "message": "Organization have created Public chanel", "data": map[string]interface{}{"info": nil}})
	}
	//
	organizationPublicChanel := new(model.OrganizationPublicChanel)
	organizationPublicChanel.Publichanelid = uuid.New().String()
	organizationPublicChanel.Name = organizationPublicChanelRequest.Name
	organizationPublicChanel.OrganizatioName = (attributes["usermail"]).(string)
	organizationPublicChanel.OrganizationId = (attributes["userid"]).(string)
	err = profileRepository2.SavePublicChanel(organizationPublicChanel)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}
	//create kafka topic
	config.ReadConfig()
	url := viper.GetString(`kafka.url`)
	integratedkafka.KafkaInit(url)
	producer, err := transport.InitProducer(url)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "Kafka connection refused", "data": map[string]interface{}{"info": nil}})
	}
	//kiem tra kafka topic = public channel nào mới push message
	transport.PublishMessage(producer, organizationPublicChanel.Publichanelid, kafka.PartitionAny, "INIT PUBLIC CHANEL "+organizationPublicChanel.OrganizatioName)
	//
	fmt.Println("Create Public Chanel success")
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Success Create Public Chanel", "data": map[string]interface{}{"info": nil}})
}
