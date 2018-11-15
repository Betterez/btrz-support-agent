package betterez

import (
	"fmt"
	"testing"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

const (
	collection = "col1"
	database   = "test"
)

func getLocalDeployment() *DeploymentData {
	return &DeploymentData{
		DatabaseName:  "test",
		ServerAddress: "localhost",
	}
}

func TestConnection(t *testing.T) {
	deploymentData := getLocalDeployment()
	session, err := mgo.Dial(deploymentData.MakeDialString())
	if err != nil {
		t.Fatal(err)
	}
	if session == nil {
		t.Fatal("db sesison is nil")
	}
}

func TestCreateRecords(t *testing.T) {
	deploymentData := getLocalDeployment()
	session, err := mgo.Dial(deploymentData.MakeDialString())
	if err != nil {
		t.Fatal(err)
	}
	if session == nil {
		t.Fatal("db sesison is nil")
	}
	session.DB(database).C(collection).DropCollection()
	for i := 0; i < 10; i++ {
		session.DB(database).C(collection).Insert(bson.M{"username": fmt.Sprintf("test%03d", i)})
	}
	records, err := session.DB(database).C(collection).Count()
	if err != nil {
		t.Fatal(err)
	}
	if records != 10 {
		t.Fatalf("Suppose to have 10 records, got %d.", records)
	}
}

func TestUpdatingRecords(t *testing.T) {
	deploymentData := getLocalDeployment()
	session, err := mgo.Dial(deploymentData.MakeDialString())
	if err != nil {
		t.Fatal(err)
	}
	if session == nil {
		t.Fatal("db sesison is nil")
	}
	err = session.DB(database).C(collection).Update(bson.M{"username": `/test/`}, bson.M{"username": "update"})
	if err != nil {
		t.Fatal(err)
	}
}

// func TestGettingDoc(t *testing.T) {
// 	deploymentData := getLocalDeployment()
// 	session, err := mgo.Dial(deploymentData.MakeDialString())
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if session == nil {
// 		t.Fatal("db sesison is nil")
// 	}
// 	data := session.DB(database).C(collection).Find(bson.M{"username": `/test/`})
// 	iter := data.Iter()
// 	for iter.Next()
// }
