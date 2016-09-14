package store

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"knowledge-base/conf"
)

//MongoConnection ...
type MongoConnection struct {
	originalSession *mgo.Session
}

var db *MongoConnection

func init() {
	db = NewDBConnection()
}

//NewDBConnection create new DB connections
func NewDBConnection() (conn *MongoConnection) {
	conn = new(MongoConnection)
	conn.createConnection()
	return
}

func (c *MongoConnection) createConnection() (err error) {
	fmt.Println("Connecting to local mongo server....")
	c.originalSession, err = mgo.Dial(conf.Config.DB.Host)
	if err == nil {
		fmt.Println("Connection established to mongo server")
	} else {
		fmt.Printf("Error occured while creating mongodb connection: %s", err.Error())
	}
	index := mgo.Index{
		Key: []string{"$text:subject", "$text:tags"},
	}
	err = c.originalSession.DB(conf.Config.DB.Name).C("Articles").EnsureIndex(index)

	return
}

//CloseConnection close DB connection
func (c *MongoConnection) CloseConnection() {
	if c.originalSession != nil {
		c.originalSession.Close()
	}
}

func (c *MongoConnection) getSessionAndCollection(collectionName string) (session *mgo.Session, collection *mgo.Collection, err error) {
	if c.originalSession != nil {
		session = c.originalSession.Copy()
		collection = session.DB(conf.Config.DB.Name).C(collectionName)
	} else {
		err = errors.New("No original session found")
	}
	return
}
