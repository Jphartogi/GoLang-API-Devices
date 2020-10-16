package main

import (
	"api/auth"
	"api/database"
	"api/global"
	"api/proto"
	"api/routes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type apiServer struct{}

func (apiServer) Login(_ context.Context, in *proto.UserLoginRequest) (*proto.UserAuthResponse, error) {
	login, password := in.GetUsername(), in.GetPassword()
	var user global.User
	user.Username = login
	user.Password = password
	_, err := auth.AuthorizeUser(&user)
	if err != nil {
		log.Fatal(err)
	}
	token, err := auth.UserTokenGenerator(&user)
	return &proto.UserAuthResponse{Token: token}, nil

}

func (apiServer) SignUp(_ context.Context, input *proto.UserSignUpRequest) (*proto.UserSuccessRegister, error) {
	username, password := input.GetUsername(), input.GetPassword()
	var user global.User
	user.Username = username
	user.Password = password
	result := database.CheckUsernameAvailability(&user)
	if !result {
		return &proto.UserSuccessRegister{Message: "Username is taken"}, errors.New("Username is taken")
	}
	_, err := auth.RegisterUser(&user)
	if err != nil {
		log.Fatal(err)
	}
	return &proto.UserSuccessRegister{Message: "Successfully added user to DB"}, nil

}

func (apiServer) UsernameTaken(_ context.Context, input *proto.UsernameTakenRequest) (*proto.UsernameisTaken, error) {
	username := input.GetUsername()
	var user global.User
	user.Username = username
	result := database.CheckUsernameAvailability(&user)

	return &proto.UsernameisTaken{Status: result}, nil

}

func (apiServer) DeleteUser(_ context.Context, input *proto.UserDeleteRequest) (*proto.UserSuccessDelete, error) {
	username := input.GetUsername()
	var user global.User
	user.Username = username
	result, err := auth.DeleteUser(&user)
	if err != nil {
		return &proto.UserSuccessDelete{Message: "Failed to delete from DB"}, errors.New(err.Error())
	}
	if !result {
		return &proto.UserSuccessDelete{Message: "Failed to delete from DB"}, errors.New("Failed to delete from DB")
	}
	return &proto.UserSuccessDelete{Message: "Successfully delete user from DB"}, nil

}

func requestHandle() {
	server := grpc.NewServer()
	proto.RegisterAPIServicesServer(server, apiServer{})

	const port = 8092
	fmt.Printf("Server is up and running at port %d ", port)
	myRoutes := mux.NewRouter()
	myRoutes.HandleFunc("/api/v1/", routes.Home).Methods("GET")
	myRoutes.HandleFunc("/api/v1/register/devices", routes.RegisterHandler).Methods("POST")
	myRoutes.HandleFunc("/api/v1/get/devices", routes.GetDeviceInfoHandler).Methods("GET")

	myRoutes.HandleFunc("/", routes.Home).Methods("GET")

	err := http.ListenAndServe(":8092", myRoutes)
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

func main() {

	requestHandle()

}
