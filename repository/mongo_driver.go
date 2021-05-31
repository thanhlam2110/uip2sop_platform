package repository

import (
	"bitbucket.org/cloud-platform/uip2sop_platform/model"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/mgo.v2"
)

// ProfileRepositoryMongo - ProfileRepositoryMongo
type ProfileRepositoryMongo struct {
	db         *mgo.Database
	collection string
}

//NewProfileRepositoryMongo - NewProfileRepositoryMongo
func NewProfileRepositoryMongo(db *mgo.Database, collection string) *ProfileRepositoryMongo {
	return &ProfileRepositoryMongo{
		db:         db,
		collection: collection,
	}
}

//<------------------USER------------------>
//SaveUser
func (r *ProfileRepositoryMongo) SaveUser(userSSO *model.UserSSO) error {
	err := r.db.C(r.collection).Insert(userSSO)
	return err
}

//FindByUser
func (r *ProfileRepositoryMongo) FindByUser(username string) (*model.UserSSO, error) {
	var userSSO model.UserSSO
	err := r.db.C(r.collection).Find(bson.M{"username": username}).One(&userSSO)
	if err != nil {
		return nil, err
	}
	return &userSSO, nil
}

//<------------------THINGS------------------>
//Save Things
func (r *ProfileRepositoryMongo) SaveThing(ioThing *model.IoThing) error {
	err := r.db.C(r.collection).Insert(ioThing)
	return err
}

//FindThings
func (r *ProfileRepositoryMongo) FindThingById(thingid string) (*model.IoThing, error) {
	var ioThing model.IoThing
	err := r.db.C(r.collection).Find(bson.M{"thingid": thingid}).One(&ioThing)
	if err != nil {
		return nil, err
	}
	return &ioThing, nil
}

//Check ThingID and Thing Key
func (r *ProfileRepositoryMongo) FindThing(thingid, thingkey string) (*model.IoThing, error) {
	var ioThing model.IoThing
	err := r.db.C(r.collection).Find(bson.M{"thingid": thingid, "thingkey": thingkey}).One(&ioThing)
	if err != nil {
		return nil, err
	}
	return &ioThing, nil
}

//<------------------CHANELS------------------>
//Save Chanels
func (r *ProfileRepositoryMongo) SaveChanel(ioTChanel *model.IoTChanel) error {
	err := r.db.C(r.collection).Insert(ioTChanel)
	return err
}

//FindChanel
func (r *ProfileRepositoryMongo) FindChanelById(chanelid string) (*model.IoTChanel, error) {
	var ioTChanel model.IoTChanel
	err := r.db.C(r.collection).Find(bson.M{"chanelid": chanelid}).One(&ioTChanel)
	if err != nil {
		return nil, err
	}
	return &ioTChanel, nil
}

//<------------------MAP THINGS TO CHANEL------------------>
//Save Chanels
func (r *ProfileRepositoryMongo) SaveMapThingChanel(mapThingChanel *model.MapThingChanel) error {
	err := r.db.C(r.collection).Insert(mapThingChanel)
	return err
}

//FindThings
func (r *ProfileRepositoryMongo) FindMapThingChanel(thingid, chanelid string) (*model.MapThingChanel, error) {
	var mapThingChanel model.MapThingChanel
	err := r.db.C(r.collection).Find(bson.M{"chanelid": chanelid, "thingid": thingid}).One(&mapThingChanel)
	if err != nil {
		return nil, err
	}
	return &mapThingChanel, nil
}

//<------------------ORGANIZATION PUBLIC CHANEL------------------>

func (r *ProfileRepositoryMongo) FindUserByUserID(userid string) (*model.OrganizationPublicChanelRequest, error) {
	var organizationPublicChanelRequest model.OrganizationPublicChanelRequest
	err := r.db.C(r.collection).Find(bson.M{"userid": userid}).One(&organizationPublicChanelRequest)
	if err != nil {
		return nil, err
	}
	return &organizationPublicChanelRequest, nil
}
func (r *ProfileRepositoryMongo) FindPublicChanel(organizationid string) (*model.OrganizationPublicChanelRequest, error) {
	var organizationPublicChanelRequest model.OrganizationPublicChanelRequest
	err := r.db.C(r.collection).Find(bson.M{"organizationid": organizationid}).One(&organizationPublicChanelRequest)
	if err != nil {
		return nil, err
	}
	return &organizationPublicChanelRequest, nil
}

//SavePublicChanel
func (r *ProfileRepositoryMongo) SavePublicChanel(organizationPublicChanel *model.OrganizationPublicChanel) error {
	err := r.db.C(r.collection).Insert(organizationPublicChanel)
	return err
}

//Find public chanel ID
func (r *ProfileRepositoryMongo) FindPublicChanelID(organizationid string) (*model.OrganizationPublicChanel, error) {
	var organizationPublicChanel model.OrganizationPublicChanel
	err := r.db.C(r.collection).Find(bson.M{"organizationid": organizationid}).One(&organizationPublicChanel)
	if err != nil {
		return nil, err
	}
	return &organizationPublicChanel, nil
}
