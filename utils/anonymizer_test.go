package utils

import (
	"testing"
)

var (
	firstName = "Donald Duck"
	email     = "someone@email.com"
)

func TestEmailAddress(t *testing.T) {
	if !IsStringEmailAddress("someone@email.com") {
		t.Fatal("email address not found")
	}
}

func TestInvalidEmailAddress(t *testing.T) {
	if IsStringEmailAddress("Hello world!") {
		t.Fatal("bad email address found")
	}
}

func TestHash(t *testing.T) {
	expected := "4f7ad45c100c79b3d23ea63832475d31"
	if GetMD5Hash(firstName) != expected {
		t.Fatalf("bad hash value for %s.\nexpected %s, got %s.", firstName, expected, GetMD5Hash(firstName))
	}
}

func TestAnonimizerValue(t *testing.T) {
	expected := "4f7ad45c100c79b3d23ea63832475d31"
	anon := CreateDatabaseRecord()
	anon.GetDataBody()["firstName"] = firstName
	anon.UpdateRecord()
	if anon.GetAnonymizedFieldsLength() != 1 {
		t.Fatalf("got %d fields instead of 1", anon.GetAnonymizedFieldsLength())
	}

	AnonymizeRecord(anon)
	if anon.GetAnonymizedFieldName(0) != "firstName" {
		t.Fatalf("%s returned - bad field name.", anon.GetAnonymizedFieldName(0))
	}
	if anon.GetAnonymizedFieldValue(0) != expected {
		t.Fatalf("%s returned - bad md5 hash value for regular value.", anon.GetAnonymizedFieldValue(0))
	}
}

func TestAnonimizerEmail(t *testing.T) {
	expected := "5c3ec817df@848797312a01b133a0f79c.com"
	anon := CreateDatabaseRecord()
	anon.AddValueToRecord("email", "someone@email.com")
	AnonymizeRecord(anon)
	if anon.GetAnonymizedFieldValue(0) != "5c3ec817df@848797312a01b133a0f79c.com" {
		t.Fatalf("bad email hash: got %s,expected %s", anon.GetAnonymizedFieldValue(0), expected)
	}
}

func TestAnonimizerMultipleValues(t *testing.T) {
	anon := CreateDatabaseRecord()
	anon.AddValueToRecord("email", email)
	anon.AddValueToRecord("firstName", firstName)
	AnonymizeRecord(anon)
	if anon.GetAnonymizedFieldValue(0) != "4f7ad45c100c79b3d23ea63832475d31" {
		t.Fatal("bad md5 hash value name")
	}
	if anon.GetAnonymizedFieldValue(1) != "5c3ec817df@848797312a01b133a0f79c.com" {
		t.Fatal("bad md5 hash value for email.")
	}
}
