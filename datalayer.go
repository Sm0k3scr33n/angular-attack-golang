package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	// "time"
	// "gopkg.in/mgo.v2/bson"
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
	// c = collection
	c := dbo.session.DB(MONGO_DB).C("challenges")

	err := c.Insert(&obj)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

/*
func main() {
	dal := NewDataLayer()
	defer dal.Close()
	obj := Challenge{Id: 2,
		Latitude:     "39.734114",
		Longitude:    "-104.9797755,15",
		ChallengeStr: "Drink a beer with Rob Stone",
		Score:        10,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now()}
	// dal.SaveStruct(&obj)
	dal.SaveGeneric(&obj)

	// fetching data from the collection...
	// bson.M is a map[string]interface{} type
	// var result bson.M
	// err := c.Find(nil).One(&result)
	// if err != nil {
	//   log.Fatal(err)
	// }
	// fmt.Println(result)
}
*/
