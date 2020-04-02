package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/spf13/viper"
)

const (
	letters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphaNumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	stringBytes := make([]byte, length)
	var randomInteger int
	for i := 0; i < length; i++ {
		randomInteger = rand.Intn(len(letters))
		stringBytes[i] = letters[randomInteger]
	}
	return string(stringBytes)
}

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
	if env == "" {
		env = "local"
	}
	viper.AddConfigPath("./config/")
	// viper.AddConfigPath("../config")
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
	if currentEnv == "" {
		currentEnv = "local"
	}
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
