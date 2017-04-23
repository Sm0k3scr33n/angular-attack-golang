package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	// "time"
	"gopkg.in/mgo.v2/bson"
)

// mongo connection string
const MONGO_HOST = "localhost"
const MONGO_DB = "test"

// type Challenge struct {
// 	Id           int
// 	Latitude     string
// 	Longitude    string
// 	ChallengeStr string
// 	Score        int
// 	CreatedAt    time.Time
// 	UpdatedAt    time.Time
// }

type DataLayerObject struct {
	session *mgo.Session // reference to mgo connection object
}

type DataLayerInterface interface {
	Open(connStr string) error
	Close() error
	SaveStruct(obj *Challenge) error
	SaveGeneric(obj interface{}) error

	SaveChallenge(obj interface{}) error
	LoadChallenge() (bson.M, error)
	GetChallengeTable() (bson.M, error)
}

// NewDataLayer() factory returns an object that implements
// DataLayerInterface. Internal state is held by DataLayerObject
func NewDataLayer() DataLayerInterface {
	obj := DataLayerObject{}
	err := obj.Open(MONGO_HOST)
	if err != nil {
		panic(err)
	}
	return &obj
}

// Open a mongo connection within the DataLayerObject
// which implements the DataLayerInterface.
func (dbo *DataLayerObject) Open(connStr string) error {
	if dbo.session != nil {
		dbo.Close()
	}
	session, err := mgo.Dial(connStr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("successfully opened dbo connection")
	dbo.session = session
	return nil
}

// Close the mongo connection (if exists) within the DataLayerObject.
func (dbo *DataLayerObject) Close() error {
	if dbo.session == nil {
		return nil
	}
	dbo.session.Close()
	dbo.session = nil
	log.Println("closed dbo session")
	return nil
}

func (dbo *DataLayerObject) SaveStruct(obj *Challenge) error {
	// c = collection
	c := dbo.session.DB(MONGO_DB).C("challenges")

	// writing data into the collection...
	fmt.Println(obj)
	err := c.Insert(&obj)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (dbo *DataLayerObject) SaveGeneric(obj interface{}) error {
	// c = collection
	c := dbo.session.DB(MONGO_DB).C("challenges")

	err := c.Insert(&obj)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (dbo *DataLayerObject) SaveChallenge(obj interface{}) error {
	log.Println("saving challenge to mongo")
	// c = collection
	c := dbo.session.DB(MONGO_DB).C("challenges")

	err := c.Insert(&obj)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (dbo *DataLayerObject) LoadChallenge() (bson.M, error) {
	var result bson.M
	c := dbo.session.DB(MONGO_DB).C("challenges")
	err := c.Find(nil).One(&result)
	if err != nil {
		return nil, err
	}
	fmt.Println(result)
	return result, nil
}

func (dbo *DataLayerObject) GetChallengeTable() (bson.M, error) {
	var results []bson.D
	c := dbo.session.DB(MONGO_DB).C("challenges")
	iter := c.Find(nil).Limit(10).Iter()
	err := iter.All(&results)
	if err != nil {
		log.Println(err)
		return nil, err
	}
        var out = bson.M{"results": results}
	return out, nil
}

