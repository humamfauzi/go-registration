package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/spf13/viper"
)

type StringArray []string

func (sa StringArray) Includes(checkString string) bool {
	for _, content := range sa {
		if content == checkString {
			return true
		}
	}
	return false
}

func unpackJson(request io.Reader) (interface{}, error) {
	var buffer interface{}
	reqBody, err := ioutil.ReadAll(request)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(reqBody, &buffer)
	return buffer, err
}

func GetEnv(key string) interface{} {
	env := os.Getenv("ENV")
	if env != "" {
		env = "local"
	}

	viper.SetConfigType("json")
	viper.SetConfigFile("../config/" + env + ".config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	result := viper.Get(key)
	return result
}

func InterpretInterfaceString(input interface{}, defaultValue string) string {
	switch input.(type) {
	case string:
		return input.(string)
	default:
		return defaultValue
	}

}

func GenerateUUID(module string, mod int) string {
	currentTime := strconv.Itoa(time.Now().Year())
	genUuid := uuid.New().String()
	parts := strings.Split(genUuid, "-")
	fmt.Println(parts)
	switch mod {
	case 1:
		genUuid = module + "/" + string(currentTime) + "/" + parts[0]
	case 2:
		genUuid = module + "/" + string(currentTime) + "/" + parts[0]
		genUuid += "/" + parts[1]
	case 3:
		genUuid = module + "/" + string(currentTime) + "/" + parts[0]
		genUuid += "/" + parts[1]
		genUuid += "/" + parts[2]
	case 4:
		genUuid = module + "/" + string(currentTime) + "/" + parts[0]
		genUuid += "/" + parts[1]
		genUuid += "/" + parts[2]
		genUuid += "/" + parts[3]
	default:
		genUuid = module + "/" + string(currentTime) + "/" + parts[0]
		genUuid += "/" + parts[1]
		genUuid += "/" + parts[2]
		genUuid += "/" + parts[3]
		genUuid += "/" + parts[4]
	}
	return genUuid
}

type databaseProfile struct{}

func databasePurgeMySQL(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	currentEnv := os.Getenv("ENV")
	allowedEnv := StringArray{"local", "test"}
	if !allowedEnv.Includes(currentEnv) {
		return errors.New("OPERATING IN FORBIDDEN ENVIRONMENT")
	}
	tables := []string{"users"}
	delForm, err := db.Prepare("DELETE FROM ?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(tables)
	return nil

}
