package config

import (
	"github.com/spf13/viper"
	mgo "gopkg.in/mgo.v2"
)

func ReadConfig() {
	// Get file config
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
func ReadClientConfig() {
	// Get file config
	viper.SetConfigFile(`organization_config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
func GetMongoDataBase() (*mgo.Database, *mgo.Session, error) {
	//
	ReadConfig()
	link := viper.GetString(`mongo.url`)
	session, err := mgo.Dial(link)
	if err != nil {
		return nil, nil, err
	}
	db := session.DB("users")
	return db, session, nil
}
