package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/humamfauzi/go-registration/utils"

	"github.com/humamfauzi/go-registration/exconn"
	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
)

var (
	db       *gorm.DB
	errorMap map[string]string
)

func main() {
	// Initialize DB connection
	db = exconn.ConnectToMySQL()
	errorMap = utils.InitError("./error.json")

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods(http.MethodGet)
	r.HandleFunc("/register", RegisterHandler).Methods(http.MethodPost)
	r.HandleFunc("/login", LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/logout", LogoutHandler).Methods(http.MethodPost)
	r.HandleFunc("/forget_password", ForgotPasswordHandler).Methods(http.MethodPost)
	r.HandleFunc("/recovery_password", RecoveryPasswordHandler).Methods(http.MethodGet)
	r.HandleFunc("/users/update", UpdateUserHandler).Methods(http.MethodPut)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Print("STARTING SERVER")
	log.Fatal(srv.ListenAndServe())
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello world")
	return
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newUser User
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var findUser User
	db.Debug().Where("email = ?", newUser.Email).Find(&findUser)
	if findUser.Id != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorMap["ERR_EMAIL_ALREADY_TAKEN"]))
		return
	}
	newUser.CreateUser()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true}`))
	return
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var loginUser User
	err = json.Unmarshal(body, &loginUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var findUser User
	err = db.Debug().Where("email = ?", loginUser.Email).Find(&findUser).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`ERR_EMAIL_PASS_NOT_MATCH`))
		return
	}
	if findUser.Password != loginUser.Password {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`ERR_EMAIL_PASS_NOT_MATCH`))
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
	return
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
	return
}

func RecoveryPasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
	return
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
	return
}
