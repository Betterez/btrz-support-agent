package utils

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	collectionNames = []string{
		"tickets",
		"transactions",
	}
	masterCollectionName = "customers"
	usersCollection      = "users"
)

// AnonymizeDB anonymize all records in the DB
func AnonymizeDB(deploymentData *DeploymentData) (bool, error) {
	session, err := mgo.Dial(deploymentData.MakeDialString())
	if err != nil {
		return false, err
	}
	defer session.Close()
	masterCollection := session.DB(deploymentData.DatabaseName).C(masterCollectionName)
	allCustomers := masterCollection.Find(bson.M{})
	record := CreateDatabaseRecord()
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
			log.Printf("%d records processed ", recordsProcessed)
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

	log.Printf("Done! %d records processed.", recordsProcessed)
	log.Println("updating users password")
	session.DB(deploymentData.DatabaseName).C(usersCollection).UpdateAll(bson.M{},
		// sha 256 12341234
		bson.M{"$set": bson.M{"password": "1718c24b10aeb8099e3fc44960ab6949ab76a267352459f203ea1036bec382c2"}})
	removeCreditData(session, deploymentData)
	return true, nil
}

func removeCreditData(session *mgo.Session, deploymentData *DeploymentData) {
	session.DB(deploymentData.DatabaseName).C("accounts").UpdateAll(bson.M{}, bson.M{"$unset": bson.M{"preferences.paymentProviders.online_credit.params": "1"}})
}
