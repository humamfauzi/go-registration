package utils

import "testing"

func TestInitError(t *testing.T) {
	errorList := InitError("./error_test.json")
	testError := errorList["ERR_EMAIL_ALREADY_TAKEN"]
	expect := "Email already taken, please use another email"
	if testError != expect {
		t.Errorf("Should have %s but got %s", testError, expect)
	}
	t.Log(testError, " --- ", expect)
}
