package service

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/uip2sop_platform/config"
	"bitbucket.org/cloud-platform/uip2sop_platform/model"
	"bitbucket.org/cloud-platform/uip2sop_platform/repository"
	"gopkg.in/mgo.v2"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

func CreateThings(c echo.Context) error {
	ioTCreateRequest := new(model.IoTCreateRequest)
	err := c.Bind(ioTCreateRequest)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}
	//parse token ---> get user status
	authResponse, err := BasicAuth(ioTCreateRequest.Token)

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
	//fmt.Println(attributes)
	//Check User status
	if attributes["userstatus"] != "ACTIVE" {
		return c.JSON(200, map[string]interface{}{"code": "8", "message": "USER DISABLED", "data": map[string]interface{}{"info": nil}})
	}
	ioThing := new(model.IoThing)
	ioThing.Thingname = ioTCreateRequest.Thingname
	ioThing.Thingid = uuid.New().String()
	ioThing.Thingkey = uuid.New().String()
	ioThing.Thingstatus = ioTCreateRequest.Thingstatus
	ioThing.Userparentid = (attributes["userparentid"]).(string)
	//fmt.Println(ioThing)
	//Save Things
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}
	profileRepository := repository.NewProfileRepositoryMongo(db, "things")
	_, err = profileRepository.FindThingById(ioThing.Thingid)
	if err != mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "2", "message": "ThingId is Exists", "data": map[string]interface{}{"info": nil}})
	}
	err = profileRepository.SaveThing(ioThing)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}
	fmt.Println("Saved Thing success")
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Success", "data": map[string]interface{}{"thing_id": ioThing.Thingid, "thing_key": ioThing.Thingkey}})
}
