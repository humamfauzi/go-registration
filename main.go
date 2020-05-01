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

	"github.com/gorilla/mux"
	// pb "github.com/humamfauzi/go-registration/protobuf"
)

func main() {
	// Initialize DB connection
	db = exconn.ConnectToMySQL()
	errorMap = utils.InitError("./error.json")

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Methods(http.MethodGet)
	r.HandleFunc("/users/register", RegisterHandler).Methods(http.MethodPost)
	r.HandleFunc("/users/login", LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/users/logout", LogoutHandler).Methods(http.MethodPost)
	r.HandleFunc("/users/forget_password", ForgotPasswordHandler).Methods(http.MethodPost)
	r.HandleFunc("/users/recovery_password", RecoveryPasswordHandler).Methods(http.MethodGet)
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

// func InitRPCServer() {
// 	flag.Parse()
// 	listen, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
// 	if err != nil {
// 		log.Fatalf("Failed to Listen: %v", err)
// 	}
// 	var opts []grpc.ServerOption
// 	grpcServer := grpc.NewServer(opts...)

// 	pb.RegisterRouteGuideServer(grpcServer, newServer())
// 	grpcServer.Serve(listen)

// }

// type routeGuideServer struct {
// 	pb.UnimplementedRouteGuideServer
// 	savedFeatures []*pb.Feature // read-only after initialized

// 	mu         sync.Mutex // protects routeNotes
// 	routeNotes map[string][]*pb.RouteNote
// }

// func newServer() *routeGuideServer {
// 	s := &routeGuideServer{routeNotes: make(map[string][]*pb.RouteNote)}
// 	s.loadFeatures(*jsonDBFile)
// 	return s
// }

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello world")
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
	opReply := OperationReply{
		"OP_USER_UPDATE",
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_UNREADBLE_PAYLOAD",
			Message: "Cannot parse incoming payload",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}

	var userUpdate User
	err = json.Unmarshal(body, userUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errReply := ErrorReply{
			Code:    "ERR_UNREADBLE_PAYLOAD",
			Message: "Cannot parse incoming payload",
		}
		opReply.SetFail()
		result, _ := CreateReply(opReply, errReply, []byte{})
		w.Write(result)
		return
	}

	sanitizedUser := User{
		Name:  userUpdate.Name,
		Phone: userUpdate.Phone,
	}
	user.UpdateUser(sanitizedUser)
	errReply := ErrorReply{}
	result, _ := CreateReply(opReply, errReply, []byte{})
	w.Write(result)
	return
}
