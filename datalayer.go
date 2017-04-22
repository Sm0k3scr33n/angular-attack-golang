package main

import (
  "log"
  "fmt"
  "time"
  "gopkg.in/mgo.v2"
  // "gopkg.in/mgo.v2/bson"
)

// mongo connection string
const MONGO_HOST = "localhost"
const MONGO_DB = "test"

type Challenge struct {
  Id int
  Latitude string
  Longitude string
  ChallengeStr string
  Score int
  CreatedAt time.Time
  UpdatedAt time.Time
}

type DataLayerInterface interface {
  Open(connStr string) error
  Close() error
  SaveStruct(obj *Challenge) error
}

type DataLayerObject struct {
  session *mgo.Session // reference to mgo connection object
}

// NewDataLayer() factory returns an object that implements
// DataLayerInterface. Internal state is held by DataLayerObject
func NewDataLayer() DataLayerInterface {
  obj := DataLayerObject{}
  obj.Open(MONGO_HOST)
  return &obj
}

// Open a mongo connection within the DataLayerObject
// which implements the DataLayerInterface.
func (dbo *DataLayerObject) Open(connStr string) error {
  session, err := mgo.Dial(connStr)
  if err != nil {
    panic(err)
  }
  log.Println("successfully opened dbo connection")
  dbo.session = session
  return nil
}

// Close the mongo connection (if exists) within the DataLayerObject.
func (dbo *DataLayerObject) Close() error {
  log.Println("closed dbo session")
  dbo.session.Close();
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
  }

  return nil
}

func main() {
  dal := NewDataLayer()
  defer dal.Close()
  obj := Challenge{Id: 2,
                   Latitude: "39.734114",
                   Longitude: "-104.9797755,15",
                   ChallengeStr: "Drink a beer with Rob Stone",
                   Score: 10,
                   CreatedAt: time.Now(),
                   UpdatedAt: time.Now()}
  dal.SaveStruct(&obj)


  // fetching data from the collection...
  // bson.M is a map[string]interface{} type
  // var result bson.M
  // err := c.Find(nil).One(&result)
  // if err != nil {
  //   log.Fatal(err)
  // }
  // fmt.Println(result)
}

