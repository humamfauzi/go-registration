package utils

import (
	"encoding/json"
	"io/ioutil"
)

func InitError(errorListDirectory string) map[string]string {
	newError := make(map[string]string)
	errorList, err := ioutil.ReadFile(errorListDirectory)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(errorList, &newError)
	return newError
}
