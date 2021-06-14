package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// Anonymizable any type with fields that can be anonymize
type Anonymizable interface {
	GetAnonymizedFieldsLength() int
	GetAnonymizedFieldName(int) string
	GetAnonymizedFieldValue(int) string
	SetAnonymizedFieldValue(int, string)
}

// GetMD5Hash return md5 hash for a given string
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// IsStringEmailAddress return true if a string is an email address
func IsStringEmailAddress(text string) bool {
	return strings.Contains(text, "@")
}

// AnonymizeEmailAddress special md5 hashing for email address
func AnonymizeEmailAddress(clearAddress string) string {
	hashedEmail := GetMD5Hash(clearAddress)
	var result string
	for index, v := range hashedEmail {
		if index == 10 {
			result = result + string("@")
		}
		result = result + string(v)
	}
	return result + ".com"
}

// AnonymizeRecord traverse through all anonymizable fields
func AnonymizeRecord(record Anonymizable) {
	fieldsNumber := record.GetAnonymizedFieldsLength()
	for i := 0; i < fieldsNumber; i++ {
		if IsStringEmailAddress(record.GetAnonymizedFieldValue(i)) {
			record.SetAnonymizedFieldValue(i, AnonymizeEmailAddress(record.GetAnonymizedFieldValue(i)))
		} else {
			record.SetAnonymizedFieldValue(i, GetMD5Hash(record.GetAnonymizedFieldValue(i)))
		}
	}
}
