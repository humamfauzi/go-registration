package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/humamfauzi/go-registration/utils"
)

type OperationReply struct {
	Name    string
	Success bool
}

func (or *OperationReply) SetFail() {
	or.Success = false
}

type Reply struct {
	Operation OperationReply `json:"operation"`
	Error     ErrorReply     `json:"error"`
	Body      interface{}    `json:"body,omitempty"`
}

type ErrorReply struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Meta    string `json:"meta,omitempty"`
}

func CreateReply(opProfile OperationReply, erProfile ErrorReply, body interface{}) ([]byte, error) {
	newReply := Reply{
		Operation: opProfile,
		Error:     erProfile,
		Body:      body,
	}
	reply, err := json.Marshal(newReply)
	if err != nil {
		return []byte{}, err
	}
	return reply, nil
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	opReply := OperationReply{
		"OP_USER_REGISTRATION",
		true,
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_CANNOT_READ_REQUEST",
			Message: "Cannot read incoming buffer",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}

	var newUser User
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(newUser)
	var findUser User
	db.Debug().Where("email = ?", newUser.Email).Find(&findUser)
	if findUser.Id != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorMap["ERR_EMAIL_ALREADY_TAKEN"]))
		errReply := ErrorReply{
			Code:    "ERR_EMAIL_ALREADY_TAKEN",
			Message: "Email already taken please use another email",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}
	passwordHash, err := GeneratePasswordHash(newUser.Email, newUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_INTERNAL_SERVER_ERROR",
			Message: "There is something wrong, please try some moment",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}
	newUser.SetPassword(passwordHash)
	newUser.CreateUser()

	w.WriteHeader(http.StatusOK)
	errReply := ErrorReply{}
	result, _ := CreateReply(opReply, errReply, []byte{})
	w.Write(result)
	return
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	opReply := OperationReply{
		"OP_USER_LOGIN",
		true,
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_CANNOT_READ_REQUEST",
			Message: "Cannot read incoming buffer",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}

	var loginUser User
	err = json.Unmarshal(body, &loginUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_CANNOT_READ_REQUEST",
			Message: "Cannot read incoming buffer",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}
	var findUser User
	err = findUser.FindUserLoginByEmail(loginUser.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_EMAIL_PASS_NOT_MATCH",
			Message: "Combination of Email and Password not found",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}
	combined := loginUser.Email + ":" + loginUser.Password + ":" + PASSWORD_SALT
	ok := ValidatePasswordHash(combined, findUser.Password)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_EMAIL_PASS_NOT_MATCH",
			Message: "Combination of Email and Password not found",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}
	token := utils.GenerateUUID("token", 4)
	findUser.UpdateUser(User{AccessToken: &token})

	payload, err := GenerateWebToken(findUser.Id, token, PURP_ACCESS_TOKEN)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_INTERNAL_SERVER_ERROR",
			Message: "There is something wrong, please try some moment",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}

	tokenStruct := struct {
		Token []byte `json:"token"`
	}{payload}

	w.WriteHeader(http.StatusOK)
	errReply := ErrorReply{}

	// Token inside result come in form of base64, any incoming token should converter
	// from base64 to normal byte befire getting parsed
	result, _ := CreateReply(opReply, errReply, tokenStruct)
	w.Write(result)
	return
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	opReply := OperationReply{
		"OP_USER_LOGOUT",
		true,
	}
	user, err := GetWebToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_UNAUTHORIZED_OPERATION",
			Message: "User unable do this operation",
			Meta:    err.Error(),
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}
	var token *string
	user.UpdateUser(User{AccessToken: token})
	errReply := ErrorReply{}

	result, _ := CreateReply(opReply, errReply, []byte{})
	w.Write(result)
	return
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	opReply := OperationReply{
		"OP_USER_FORGOT_PASS_REQUEST",
		true,
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_CANNOT_READ_REQUEST",
			Message: "Cannot read incoming buffer",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}

	var loginUser User
	err = json.Unmarshal(body, &loginUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_CANNOT_READ_REQUEST",
			Message: "Cannot read incoming buffer",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}
	var findUser User
	err = findUser.FindUserLoginByEmail(loginUser.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_EMAIL_NOT_REGISTERED",
			Message: "Email not found",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}
	passToken := utils.GenerateUUID("token", 4)
	findUser.UpdateUser(User{PassToken: &passToken})

	forgotEmailMeta := make(map[string]interface{})
	forgotEmailMeta["token"] = passToken
	forgotEmailMeta["time"] = time.Now().Format("2020-11-12T09:00:33")
	forgotEmailMeta["name"] = findUser.Name

	forgotEmailNotif := Notification{
		Code:    "NOTIF_FORGOT_PASSWORD_EMAIL",
		Message: "A request to renew password",
		Meta:    forgotEmailMeta,
	}
	err = NotificationRouter(forgotEmailNotif)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errReply := ErrorReply{
			Code:    "ERR_INTERNAL_SERVER_ERROR",
			Message: "There is something wrong, please try some moment",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}

	errReply := ErrorReply{}
	result, _ := CreateReply(opReply, errReply, []byte{})
	w.Write(result)
	w.WriteHeader(http.StatusOK)
	return
}

func RecoveryPasswordHandler(w http.ResponseWriter, r *http.Request) {
	opReply := OperationReply{
		"OP_USER_RECOVERY_PASS_REQUEST",
		true,
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_CANNOT_READ_REQUEST",
			Message: "Cannot read incoming buffer",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}
	tokenPassword := struct {
		Token    string
		Password string
	}{}

	err = json.Unmarshal(body, &tokenPassword)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_CANNOT_READ_REQUEST",
			Message: "Cannot read incoming buffer",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}

	user, err := VerifyToken(tokenPassword.Token, PURP_FORGOT_PASSWORD)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_TOKEN_NOT_FOUND",
			Message: "Cannot find request token",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}
	passwordHash, err := GeneratePasswordHash(user.Email, tokenPassword.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_INTERNAL_SERVER_ERROR",
			Message: "There is something wrong, please try some moment",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}

	user.UpdateUser(User{Password: passwordHash})
	errReply := ErrorReply{}
	result, _ := CreateReply(opReply, errReply, []byte{})
	w.Write(result)
	w.WriteHeader(http.StatusOK)

}
