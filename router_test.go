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

	operationsAssert := assert.New(t)
	operationsAssert.Equal(operationStatus.Name, "OP_USER_REGISTRATION", "should be the same")
	operationsAssert.Equal(operationStatus.Success, true, "should be the same")
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

	// Choosing nonexist user
	exampleJson = `{"email": "nonexist@user.com", "password":"Paass123!"}`
	reqLogin := httptest.NewRequest("POST", "http://example.com/users/login", bytes.NewReader([]byte(exampleJson)))
	reqLogin.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	LoginHandler(w, reqLogin)

	body, _ := ioutil.ReadAll(w.Result().Body)
	var someBody Reply

	json.Unmarshal(body, &someBody)
	operationStatus := someBody.Operation
	errorStatus := someBody.Error

	operationsAssert := assert.New(t)
	operationsAssert.Equal(operationStatus.Name, "OP_USER_LOGIN", "should be the same")
	operationsAssert.Equal(operationStatus.Success, false, "should be the same")

	operationsAssert.Equal(errorStatus.Code, "ERR_EMAIL_PASS_NOT_MATCH", "should be the same")
	operationsAssert.Equal(errorStatus.Message, "Combination of Email and Password not found", "should be the same")

	// Choosing wrong password
	exampleJson = `{"email": "new@user.com", "password":"Paass123!"}`
	reqLogin = httptest.NewRequest("POST", "http://example.com/users/login", bytes.NewReader([]byte(exampleJson)))
	reqLogin.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	LoginHandler(w, reqLogin)

	body, _ = ioutil.ReadAll(w.Result().Body)

	json.Unmarshal(body, &someBody)
	operationStatus = someBody.Operation
	errorStatus = someBody.Error

	operationsAssert.Equal(operationStatus.Name, "OP_USER_LOGIN", "should be the same")
	operationsAssert.Equal(operationStatus.Success, false, "should be the same")

	operationsAssert.Equal(errorStatus.Code, "ERR_EMAIL_PASS_NOT_MATCH", "should be the same")
	operationsAssert.Equal(errorStatus.Message, "Combination of Email and Password not found", "should be the same")

	// Choosing the right password
	exampleJson = `{"email": "new@user.com", "password":"Pass123!"}`
	reqLogin = httptest.NewRequest("POST", "http://example.com/users/login", bytes.NewReader([]byte(exampleJson)))
	reqLogin.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	LoginHandler(w, reqLogin)

	body, _ = ioutil.ReadAll(w.Result().Body)

	json.Unmarshal(body, &someBody)
	operationStatus = someBody.Operation

	operationsAssert.Equal(operationStatus.Name, "OP_USER_LOGIN", "should be the same")
	operationsAssert.Equal(operationStatus.Success, true, "should be the same")

	payload := someBody.Body.(map[string]interface{})
	operationsAssert.Contains(payload, "token")
	return
}

func TestLogoutHandler(t *testing.T) {
	db = exconn.ConnectToMySQL()
	defer db.Close()

	db.AutoMigrate(&User{})
	db.Delete(&User{})

	// Creating user
	exampleJson := `{"email": "new@user.com", "password": "Pass123!", "name":"user_logout_test"}`
	reqRegistration := httptest.NewRequest("POST", "http://example.com/users/register", bytes.NewReader([]byte(exampleJson)))
	reqRegistration.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	RegisterHandler(w, reqRegistration)

	// login user
	exampleJson = `{"email": "new@user.com", "password": "Pass123!"}`
	reqLogin := httptest.NewRequest("POST", "http://example.com/users/login", bytes.NewReader([]byte(exampleJson)))
	reqLogin.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	LoginHandler(w, reqLogin)

	body, _ := ioutil.ReadAll(w.Result().Body)
	var someBody Reply

	json.Unmarshal(body, &someBody)
	payload := someBody.Body.(map[string]interface{})
	payloadToken := payload["token"].(string)

	reqLogout := httptest.NewRequest("POST", "http://example.com/users/logout", bytes.NewReader([]byte(exampleJson)))
	reqLogout.Header.Set("Content-Type", "application/json")
	reqLogout.Header.Set("Authorization", "Bearer "+payloadToken)
	w = httptest.NewRecorder()
	LogoutHandler(w, reqLogout)

	body, _ = ioutil.ReadAll(w.Result().Body)
	json.Unmarshal(body, &someBody)
	operationStatus := someBody.Operation
	t.Log(someBody)
	operationsAssert := assert.New(t)
	operationsAssert.Equal(operationStatus.Name, "OP_USER_LOGOUT", "should be the same")
	operationsAssert.Equal(operationStatus.Success, true, "should be the same")

}
