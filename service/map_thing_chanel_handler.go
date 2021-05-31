package service

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/cloud-platform/uip2sop_platform/config"
	"bitbucket.org/cloud-platform/uip2sop_platform/model"
	"bitbucket.org/cloud-platform/uip2sop_platform/repository"
	"gopkg.in/mgo.v2"

	"github.com/labstack/echo"
)

func CreateMapThingChanel(c echo.Context) error {
	mapThingChanelRequest := new(model.MapThingChanelRequest)
	err := c.Bind(mapThingChanelRequest)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}
	//parse token ---> get user status
	authResponse, err := BasicAuth(mapThingChanelRequest.Token)

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
	mapThingChanel := new(model.MapThingChanel)
	mapThingChanel.Thingid = mapThingChanelRequest.Thingid
	mapThingChanel.Chanelid = mapThingChanelRequest.Chanelid
	mapThingChanel.Mapstatus = mapThingChanelRequest.Mapstatus
	mapThingChanel.Userparentid = (attributes["userparentid"]).(string)
	//fmt.Println(ioThing)
	//Save Things
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}
	//check thing_id ton tai
	profileRepositorya := repository.NewProfileRepositoryMongo(db, "things")
	_, err = profileRepositorya.FindThingById(mapThingChanel.Thingid)
	if err == mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "4", "message": "ThingId Is Not Exists", "data": map[string]interface{}{"info": nil}})
	}
	//check chanel_id ton tai
	profileRepositoryb := repository.NewProfileRepositoryMongo(db, "chanels")
	_, err = profileRepositoryb.FindChanelById(mapThingChanel.Chanelid)
	if err == mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "4", "message": "ChanelId Is Not Exists", "data": map[string]interface{}{"info": nil}})
	}
	//insert map thing data
	profileRepository := repository.NewProfileRepositoryMongo(db, "things_map_chanels")
	_, err = profileRepository.FindMapThingChanel(mapThingChanel.Thingid, mapThingChanel.Chanelid)
	if err != mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "2", "message": "ThingId is mapped to Chanel", "data": map[string]interface{}{"info": nil}})
	}
	err = profileRepository.SaveMapThingChanel(mapThingChanel)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}
	fmt.Println("Map Thing to Chanel success")
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Success", "data": map[string]interface{}{"info": nil}})
}
