package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bitbucket.org/cloud-platform/uip2sop_platform/model"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RequestSSOTokenv2(c echo.Context) error {
	tokenRequestBodyv2 := new(model.TokenRequestBodyv2)
	err := c.Bind(tokenRequestBodyv2)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"token": nil, "refresh_token": nil}})
	}
	username := tokenRequestBodyv2.Username
	password := tokenRequestBodyv2.Password
	clientID := tokenRequestBodyv2.Clientid
	clientSecret := tokenRequestBodyv2.Clientsecret
	_, userstatus := CheckUserStatus(username)
	if userstatus != "ACTIVE" {
		return c.JSON(200, map[string]interface{}{"code": "8", "message": "USER IS DISABLE", "data": map[string]interface{}{"token": nil, "refresh_token": nil}})
	}
	token, refreshtoken := RequestTokenv2(clientID, clientSecret, username, password)
	if token == "" {
		return c.JSON(200, map[string]interface{}{"code": "8", "message": "Authen Failed", "data": map[string]interface{}{"token": nil, "refresh_token": nil}})
	}
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "success", "data": map[string]interface{}{"token": token, "refresh_token": refreshtoken}})
}
func ParseSSOToken(c echo.Context) error {
	authenRequestBody := new(model.AuthenRequestBody)
	err := c.Bind(authenRequestBody)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, map[string]interface{}{"code": "6", "message": "Body is Invalid", "data": map[string]interface{}{"info": nil}})
	}
	token := authenRequestBody.Token
	authResponse, err := BasicAuth(token)
	//
	if err != nil {
		return c.JSON(200, map[string]interface{}{"code": "10", "message": "Connection refused", "data": map[string]interface{}{"info": nil}})
	}
	//
	var result map[string]interface{}
	json.Unmarshal([]byte(authResponse), &result)

	if result["error"] != nil {
		if result["message"] != nil {
			return c.JSON(200, map[string]interface{}{"code": "7", "message": (result["message"]).(string), "data": map[string]interface{}{"info": nil}})
		}
		return c.JSON(200, map[string]interface{}{"code": "5", "message": (result["error"]).(string), "data": map[string]interface{}{"info": nil}})
	}
	//test
	attributes := result["attributes"].(map[string]interface{})
	attributes["username"] = result["id"]
	return c.JSON(200, map[string]interface{}{"code": "0", "message": "success", "data": map[string]interface{}{"info": attributes}})
}
func CheckUserStatus(username string) (err error, info string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:Vnpt123@54.169.219.90:27017/cas?authSource=admin"))
	if err != nil {
		return err, ""
	}
	defer client.Disconnect(ctx)
	database := client.Database("users")
	users := database.Collection("users")
	var result bson.M
	if err = users.FindOne(ctx, bson.M{"username": username}).Decode(&result); err != nil {

	}
	return nil, (result["userstatus"]).(string)
}
