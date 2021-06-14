package utils

// import (
// 	"fmt"
// )

type DatabaseRecord struct {
	dataBody             map[string]interface{}
	selectedFieldsLength int
	selectedFields       []string
}

var (
	fieldNames = []string{"firstName", "lastName", "email"}
)

func CreateDatabaseRecord() *DatabaseRecord {
	return &DatabaseRecord{
		selectedFieldsLength: -1,
		dataBody:             make(map[string]interface{}),
	}

}

func (record *DatabaseRecord) GetDataBody() map[string]interface{} {
	return record.dataBody
}

func (record *DatabaseRecord) GetAnonymizedFieldsLength() int {
	if record.selectedFieldsLength == -1 {
		record.UpdateRecord()
	}
	return record.selectedFieldsLength
}

func (record *DatabaseRecord) UpdateRecord() map[string]string {
	var result = make(map[string]string)
	record.selectedFields = make([]string, 0)
	var totalResult int
	for _, fieldName := range fieldNames {
		_, ok := record.GetDataBody()[fieldName]
		if ok {
			record.selectedFields = append(record.selectedFields, fieldName)
			str, ok := record.GetDataBody()[fieldName].(string)
			if ok {
				result[fieldName] = str
			}
			totalResult += 1
		}
	}
	record.selectedFieldsLength = totalResult
	return result
}

func (record *DatabaseRecord) GetAnonymizedFieldName(index int) string {
	if index > record.selectedFieldsLength {
		return ""
	}
	return record.selectedFields[index]
}

func (record *DatabaseRecord) GetAnonymizedFieldValue(index int) string {
	if index > record.selectedFieldsLength {
		return ""
	}
	result, ok := record.GetDataBody()[record.selectedFields[index]].(string)
	if ok {
		return result
	}
	return ""
}

func (record *DatabaseRecord) SetAnonymizedFieldValue(index int, data string) {
	record.GetDataBody()[record.selectedFields[index]] = data
}

func (record *DatabaseRecord) GetSelectedFieldAt(index int) string {
	if index > len(record.selectedFields) {
		return ""
	}
	return record.selectedFields[index]
}

func (record *DatabaseRecord) AddValueToRecord(key, value string) {
	record.GetDataBody()[key] = value
	record.UpdateRecord()
}
