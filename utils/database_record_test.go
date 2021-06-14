package utils

import (
	"testing"
)

var (
	firstName = "Donald Duck"
	email     = "someone@email.com"
)

func TestNewRecord(t *testing.T) {
	data := CreateDatabaseRecord()
	if len(data.GetDataBody()) != 0 {
		t.Fatal("body was badly initialized")
	}
}

func TestFieldsCount(t *testing.T) {
	data := CreateDatabaseRecord()
	data.AddValueToRecord("email", email)
	if data.GetAnonymizedFieldsLength() != 1 {
		t.Fatal("Bad fields count")
	}
}

func TestFieldsCount2(t *testing.T) {
	data := CreateDatabaseRecord()
	data.AddValueToRecord("email1", email)
	if data.GetAnonymizedFieldsLength() != 0 {
		t.Fatal("Bad fields count")
	}
}

func TestFieldsCount3(t *testing.T) {
	data := CreateDatabaseRecord()
	data.AddValueToRecord("email", email)
	data.AddValueToRecord("firstName", "Jesus")
	if data.GetAnonymizedFieldsLength() != 2 {
		t.Fatal("Bad fields count")
	}
}

func TestFieldAssignment(t *testing.T) {
	data := CreateDatabaseRecord()
	data.GetDataBody()["firstName"] = firstName
	data.UpdateRecord()
	if data.GetAnonymizedFieldValue(0) != firstName {
		t.Fatalf("failed to get correct value, got %s", data.GetAnonymizedFieldValue(0))
	}
}

func TestReturnedMap(t *testing.T) {
	data := CreateDatabaseRecord()
	data.GetDataBody()["firstName"] = firstName
	value := data.UpdateRecord()
	if value["firstName"] != firstName {
		t.Fatalf("expected %s, got %s", firstName, value["firstName"])
	}
}

func TestReturnedMultipleMap(t *testing.T) {
	data := CreateDatabaseRecord()
	data.GetDataBody()["firstName"] = firstName
	data.GetDataBody()["email"] = email
	value := data.UpdateRecord()
	if value["firstName"] != firstName {
		t.Fatalf("expected %s, got %s", firstName, value["firstName"])
	}
	if value["email"] != email {
		t.Fatalf("expected %s, got %s", email, value["firstName"])
	}
}
