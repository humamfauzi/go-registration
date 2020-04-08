package main

import "testing"

func TestInterfaceToString(t *testing.T) {
	asd := make(map[string]string)
	asd["asd"] = "asd"

	result1 := interfaceToString(asd["asd"], "def")
	if result1 != "asd" {
		t.Errorf("Got: %s, Expected: %s", result1, "asd")
	}

	result2 := interfaceToString(2, "Two")
	if result2 != "Two" {
		t.Errorf("Got: %s, Expected: %s", result2, "Two")
	}

}
