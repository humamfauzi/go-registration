package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/humamfauzi/go-registration/utils"

	"github.com/humamfauzi/go-registration/exconn"
	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
	pb "github.com/humamfauzi/go-registration/proto"
)

var (
	db *gorm.DB
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

func InitRPCServer() {
	flag.Parse()
	listen, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterRouteGuideServer(grpcServer, newServer())
	grpcServer.Serve(listen)

}

type routeGuideServer struct {
	pb.UnimplementedRouteGuideServer
	savedFeatures []*pb.Feature // read-only after initialized

	mu         sync.Mutex // protects routeNotes
	routeNotes map[string][]*pb.RouteNote
}

func newServer() *routeGuideServer {
	s := &routeGuideServer{routeNotes: make(map[string][]*pb.RouteNote)}
	s.loadFeatures(*jsonDBFile)
	return s
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello world")
	return
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
			Code: "ERR_CANNOT_READ_REQUEST",
			Message: "Cannot read incoming buffer"
		}
		opReply.Flip()
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
	var findUser User
	db.Debug().Where("email = ?", newUser.Email).Find(&findUser)
	if findUser.Id != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorMap["ERR_EMAIL_ALREADY_TAKEN"]))
		errReply := ErrorReply{
			Code: "ERR_EMAIL_ALREADY_TAKEN",
			Message: "Email already taken please use another email"
		}
		opReply.Flip()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}
	passwordHash, err := GeneratePasswordHash(newUser.Email, newUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code: "ERR_INTERNAL_SERVER_ERROR",
			Message: "There is something wrong, please try some moment"
		}
		opReply.Flip()
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
			Code: "ERR_CANNOT_READ_REQUEST",
			Message: "Cannot read incoming buffer"
		}
		opReply.Flip()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}

	var loginUser User
	err = json.Unmarshal(body, &loginUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code: "ERR_CANNOT_READ_REQUEST",
			Message: "Cannot read incoming buffer"
		}
		opReply.Flip()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}
	var findUser User
	err = db.Debug().Where("email = ?", loginUser.Email).Find(&findUser).Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code: "ERR_EMAIL_PASS_NOT_MATCH",
			Message: "Combination of Email and Password not found"
		}
		opReply.Flip()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}
	combined := loginUser.Email + ":" + loginUser.Password + ":" + PASSWORD_SALT
	ok := ValidatePasswordHash(combined, findUser.Password)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code: "ERR_EMAIL_PASS_NOT_MATCH",
			Message: "Combination of Email and Password not found"
		}
		opReply.Flip()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}
	w.WriteHeader(http.StatusOK)
	errReply := ErrorReply{}
	result, _ := CreateReply(opReply, errReply, []byte{})
	w.Write(result)
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
