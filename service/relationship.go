package service

import (
	"context"
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserSSO struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Userid       string `json:"userid"`
	Usermail     string `json:"usermail"`
	Userstatus   string `json:"userstatus"`
	Userparentid string `json:"userparentid"`
	Usertype     string `json:"usertype"`
}
type UserSSOs []UserSSO

func FindRootOfUser(node string) (err error, info string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://casuser:Mellon@13.212.137.148:27017/cas"))
	if err != nil {
		return err, ""
	}
	defer client.Disconnect(ctx)
	database := client.Database("users")
	users := database.Collection("users")
	matchStage := bson.D{{"$match", bson.D{{"username", node}}}}
	graphStage := bson.D{{"$graphLookup", bson.D{{"from", "users"}, {"startWith", "$username"}, {"connectFromField", "userparentid"}, {"connectToField", "username"}, {"as", "ancestors"}}}}
	unWind := bson.D{{"$unwind", "$ancestors"}}
	replaceRoot := bson.D{{"$replaceRoot", bson.D{{"newRoot", "$ancestors"}}}}
	proJect := bson.D{{"$project", bson.D{{"ancestors", 0}}}}
	showInfoCursor, err := users.Aggregate(ctx, mongo.Pipeline{matchStage, graphStage, unWind, replaceRoot, proJect})
	if err != nil {
		return err, ""
	}
	var showsWithInfo []bson.M
	/*if err = showInfoCursor.All(ctx, &showsWithInfo); err != nil {
		return err, ""
	}*/
	if err = showInfoCursor.All(ctx, &showsWithInfo); err != nil {
		return err, ""
	}
	bs, err := json.Marshal(showsWithInfo)
	if err != nil {
		panic(err)
	}

	var userSSOs UserSSOs
	if err := json.Unmarshal(bs, &userSSOs); err != nil {
		panic(err)
	}
	var user_rootid string
	for i := 0; i < len(userSSOs); i++ {
		if userSSOs[i].Usertype == "representation" {
			user_rootid = userSSOs[i].Userid
			//fmt.Println(userSSOs[i].Username)
			break
		}
	}
	return nil, user_rootid
}
