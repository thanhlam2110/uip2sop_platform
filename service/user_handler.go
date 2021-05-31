package service

import (
	"fmt"

	"bitbucket.org/cloud-platform/uip2sop_platform/config"
	"bitbucket.org/cloud-platform/uip2sop_platform/model"
	"bitbucket.org/cloud-platform/uip2sop_platform/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
)

func CreateUser(c echo.Context) error {
	userSSO := new(model.UserSSO)
	err := c.Bind(userSSO)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}
	if len(userSSO.Password) < 8 {
		return c.JSON(200, map[string]interface{}{"code": "6", "message": "Password less 8 character", "data": map[string]interface{}{"info": nil}})
	}
	db, session, err := config.GetMongoDataBase()
	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}

	profileRepository := repository.NewProfileRepositoryMongo(db, "users")
	_, err = profileRepository.FindByUser(userSSO.Username)
	if err != mgo.ErrNotFound {
		return c.JSON(200, map[string]interface{}{"code": "2", "message": "User is Exists", "data": map[string]interface{}{"info": nil}})
	}
	userSSO.Userid = uuid.New().String()
	err = profileRepository.SaveUser(userSSO)
	if err != nil {
		fmt.Println(err)
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "MongoDB connection refused", "data": map[string]interface{}{"info": nil}})
	}
	fmt.Println("Saved User success")
	session.Close()
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "Success", "data": map[string]interface{}{"info": nil}})
}
