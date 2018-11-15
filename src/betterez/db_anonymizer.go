package betterez

import (
	"betterez/data"
	"os"

	"github.com/bsphere/le_go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	collectionNames = []string{
		"tickets",
		"transactions",
	}
	masterCollectionName = "customers"
)

// AnonymizeDB anonymize all records in the DB
func AnonymizeDB(deploymentData *DeploymentData) (bool, error) {
	leToken := os.Getenv("LE_TOKEN")
	le, _ := le_go.Connect(leToken)
	session, err := mgo.Dial(deploymentData.MakeDialString())
	if err != nil {
		return false, err
	}
	defer session.Close()
	masterCollection := session.DB(deploymentData.DatabaseName).C(masterCollectionName)
	allCustomers := masterCollection.Find(bson.M{})
	record := data.CreateDatabaseRecord()
	recordsProcessed := 0
	iterator := allCustomers.Iter()
	doneUpdates := make(map[string]int)
	for iterator.Next(record.GetDataBody()) {
		anonymizableFields := record.UpdateRecord()
		AnonymizeRecord(record)
		updateKey := anonymizableFields["firstName"] + anonymizableFields["lastName"] + anonymizableFields["email"]
		_, ok := doneUpdates[updateKey]
		recordsProcessed++
		if recordsProcessed%100000 == 0 && recordsProcessed > 0 {
			le, _ = le_go.Connect(leToken)
			le.Printf("%d records processed ", recordsProcessed)
			le.Close()
		}
		if ok {
			doneUpdates[updateKey] = 1 + doneUpdates[updateKey]
		} else {
			for _, collectionName := range collectionNames {
				currentCollection := session.DB(deploymentData.DatabaseName).C(collectionName)
				currentCollection.UpdateAll(
					bson.M{"customerNumber": record.GetDataBody()["customerNumber"]},
					bson.M{
						"$set": bson.M{
							"firstName": record.GetAnonymizedFieldValue(0),
							"lastName":  record.GetAnonymizedFieldValue(1),
							"email":     record.GetAnonymizedFieldValue(2)},
					})
			}
		}
		masterCollection.Update(bson.M{"_id": record.GetDataBody()["_id"]}, record.GetDataBody())
	}
	le, _ = le_go.Connect(leToken)
	le.Printf("Done! %d records processed.", recordsProcessed)
	le.Close()
	return true, nil
}
