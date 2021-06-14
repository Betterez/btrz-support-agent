package utils

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

type composeData struct {
	Version int64 `bson:"version"`
	Minor   int64 `bson:"minor"`
}

type dataTest struct {
	Username string      `bson:"Username"`
	Data     string      `bson:"Data"`
	Password string      `bson:"Password"`
	Compose  composeData `bson:"compose"`
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
		session.DB(database).C(collection).Insert(
			bson.M{"Username": fmt.Sprintf("test%03d", i),
				"Data":     "hello",
				"Password": "secret",
				"compose":  bson.M{"version": 11, "minor": 22},
			})
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
		bson.M{"$set": bson.M{"Password": "1718c24b10aeb8099e3fc44960ab6949ab76a267352459f203ea1036bec382c2"}})
	if err != nil {
		t.Fatal(err)
	}
	if ok.Updated != 10 {
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

func TestDeletingSpecificKey(t *testing.T) {
	deploymentData := getLocalDeployment()
	session, err := mgo.Dial(deploymentData.MakeDialString())
	if err != nil {
		t.SkipNow()
	}
	values, err := session.DB(database).C(collection).UpdateAll(bson.M{},
		bson.M{"$unset": bson.M{"compose.minor": "1"}})
	if err != nil {
		t.Fatal(err, "when running a command")
	}
	if values.Updated != 10 {
		t.Errorf("number of values effected is %d and not 10", values.Updated)
	}
	query := session.DB(database).C(collection).Find(bson.M{"Username": "test000"})
	records, _ := query.Count()
	if records < 1 {
		t.Fatal("no records to check")
	}
	iter := query.Iter()
	var info dataTest
	for iter.Next(&info) {
		if info.Compose.Minor != 0 {
			t.Error("data suppose to be cleared")
		}
	}
}
