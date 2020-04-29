package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/humamfauzi/go-registration/exconn"
	"github.com/stretchr/testify/assert"
)

func TestRegistrationHandler(t *testing.T) {
	db = exconn.ConnectToMySQL()
	defer db.Close()

	db.AutoMigrate(&User{})
	exampleJson := `{"email": "example@user.com", "password": "heya!", "name":"user_registration_test"}`
	req := httptest.NewRequest("POST", "http://example.com/users/register", bytes.NewReader([]byte(exampleJson)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	RegisterHandler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	var someBody Reply

	json.Unmarshal(body, &someBody)
	operationStatus := someBody.Operation

	oprationsAssert := assert.New(t)
	oprationsAssert.Equal(operationStatus.Name, "OP_USER_REGISTRATION", "should be the same")
	oprationsAssert.Equal(operationStatus.Success, true, "should be the same")
	return
}

func TestLoginHandler(t *testing.T) {
	db = exconn.ConnectToMySQL()
	defer db.Close()

	db.AutoMigrate(&User{})
	db.Delete(&User{})

	exampleJson := `{"email": "new@user.com", "password": "Pass123!", "name":"user_login_test"}`
	reqRegistration := httptest.NewRequest("POST", "http://example.com/users/register", bytes.NewReader([]byte(exampleJson)))
	reqRegistration.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	RegisterHandler(w, reqRegistration)

	// Choosing wrong password
	exampleJson = `{"email": "new@user.com", "password":"Paass123!"}`
	reqLogin := httptest.NewRequest("POST", "http://example.com/users/login", bytes.NewReader([]byte(exampleJson)))
	reqRegistration.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	LoginHandler(w, reqLogin)

	body, _ := ioutil.ReadAll(w.Result().Body)
	var someBody Reply

	json.Unmarshal(body, &someBody)
	operationStatus := someBody.Operation
	errorStatus := someBody.Error

	oprationsAssert := assert.New(t)
	oprationsAssert.Equal(operationStatus.Name, "OP_USER_LOGIN", "should be the same")
	oprationsAssert.Equal(operationStatus.Success, false, "should be the same")

	oprationsAssert.Equal(errorStatus.Code, "ERR_EMAIL_PASS_NOT_MATCH", "should be the same")
	oprationsAssert.Equal(errorStatus.Message, "Combination of Email and Password not found", "should be the same")

	// Choosing the right password
	exampleJson = `{"email": "new@user.com", "password":"Pass123!"}`
	reqLogin = httptest.NewRequest("POST", "http://example.com/users/login", bytes.NewReader([]byte(exampleJson)))
	reqRegistration.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	LoginHandler(w, reqLogin)

	body, _ = ioutil.ReadAll(w.Result().Body)

	json.Unmarshal(body, &someBody)
	operationStatus = someBody.Operation

	oprationsAssert.Equal(operationStatus.Name, "OP_USER_LOGIN", "should be the same")
	oprationsAssert.Equal(operationStatus.Success, true, "should be the same")
}
