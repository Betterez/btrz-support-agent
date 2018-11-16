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

type dataTest struct {
	Username string `bson:"Username"`
	Data     string `bson:"Data"`
	Password string `bson:"Password"`
}

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
		t.SkipNow()
	}
	if session == nil {
		t.Fatal("db sesison is nil")
	}
}

func TestCreateRecords(t *testing.T) {
	deploymentData := getLocalDeployment()
	session, err := mgo.Dial(deploymentData.MakeDialString())
	if err != nil {
		t.SkipNow()
	}
	if session == nil {
		t.Fatal("db sesison is nil")
	}
	session.DB(database).C(collection).DropCollection()
	for i := 0; i < 10; i++ {
		session.DB(database).C(collection).Insert(bson.M{"Username": fmt.Sprintf("test%03d", i), "Data": "hello", "Password": "secret"})
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
		t.SkipNow()
	}
	if session == nil {
		t.Fatal("db sesison is nil")
	}
	ok, err := session.DB(database).C(collection).UpdateAll(bson.M{},
		bson.M{"$set": bson.M{"Password": "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"}})
	if err != nil {
		t.Fatal(err)
	}
	if ok.Updated < 1 {
		t.Fatalf("bad number of records updated: %d", ok.Updated)
	}
}

func TestGettingDoc(t *testing.T) {
	deploymentData := getLocalDeployment()
	session, err := mgo.Dial(deploymentData.MakeDialString())
	if err != nil {
		t.SkipNow()
	}
	if session == nil {
		t.Fatal("db sesison is nil")
	}
	query := session.DB(database).C(collection).Find(bson.M{"Username": "test000"})
	records, _ := query.Count()
	if records < 1 {
		t.Fatal("no records to check")
	}
	const dataName = "hello"
	iter := query.Iter()
	var info dataTest
	for iter.Next(&info) {
		if info.Data != dataName {
			t.Fatalf("data should be %s, but we got %s", dataName, info.Data)
		}
	}
}
