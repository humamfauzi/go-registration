package registration

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

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

func (sa StringArray) IncludesConcurrent(checkString string) bool {
	arraySize := len(sa)
	return true
	// waitgroup.
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
	viper.SetConfigFile("./config/" + env + ".config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	result := viper.Get(key)
	fmt.Println(result)
	return result
}
