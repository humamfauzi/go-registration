package utils

import (
	"testing"
)

func TestGetEnv(t *testing.T) {
	database := GetEnv("database.mysql")
	switch database.(type) {
	case nil:
		t.Errorf("DATABASE SHOULD NOT EMPTY")
	case map[string]string:
		mapDatabase := database.(map[string]string)
		t.Logf(mapDatabase["database"])
	}
}

func TestStringArrayIncludes(t *testing.T) {
	exampleArray := StringArray([]string{"asdf", "hjkl"})
	if !exampleArray.Includes("asdf") {
		t.Errorf("SHOULD GAVE TRUE VALUE")
	}
	if exampleArray.Includes("uiop") {
		t.Error("SHOULD GAVE FALSE ANSWER")
	}

}
