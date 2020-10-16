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
	"net"
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
	if username == "" || password == "" {
		return &proto.UserSuccessRegister{Message: "Please insert username and password correctly"}, errors.New("Please insert username and password correctly")
	}
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

func (apiServer) UpdateDevice(_ context.Context, input *proto.DeviceUpdateRequest) (*proto.DeviceSuccessUpdate, error) {
	id, name, location, category := input.GetDeviceID(), input.GetDeviceName(), input.GetDeviceLocation(), input.GetDeviceCategory()
	if id == "" {
		return &proto.DeviceSuccessUpdate{}, errors.New("Please specify the ID of the device")
	}
	//TODO search in database first the initial value, so when it is empty, get the initial value
	var devices global.NewDevice
	devices.DeviceID = id
	devices.DeviceName = name
	devices.DeviceLocation = location
	devices.DeviceCategory = category
	_, err := auth.UpdateDevice(&devices)
	if err != nil {
		return &proto.DeviceSuccessUpdate{}, err
	}
	return &proto.DeviceSuccessUpdate{DeviceID: id}, nil
}

func (apiServer) RegisterDevice(_ context.Context, input *proto.DeviceRequest) (*proto.DeviceSuccessRegister, error) {
	deviceName, deviceCategory, deviceLocation := input.GetDeviceName(), input.GetDeviceCategory(), input.GetDeviceLocation()
	if deviceName == "" || deviceCategory == "" || deviceLocation == "" {
		return &proto.DeviceSuccessRegister{}, errors.New("Please insert all the required field")
	}
	var devices global.NewDevice
	devices.DeviceName = deviceName
	devices.DeviceLocation = deviceLocation
	devices.DeviceCategory = deviceCategory
	id, token, err := auth.RegisterDevice(&devices)
	if err != nil {
		return &proto.DeviceSuccessRegister{}, errors.New("Error in registering the device")
	}
	return &proto.DeviceSuccessRegister{DeviceID: id, DeviceToken: token}, nil
}

// func (apiServer) UpdateDevice(_ context.Context, input *proto.DeviceRequest) (*proto.SuccessUpdate, error) {

// }

// func (apiServer) DeleteDevice(_ context.Context, input *proto.DeviceDeleteRequest) (*proto.SuccessUpdate, error) {

// }

func main() {
	server := grpc.NewServer()
	proto.RegisterAPIServicesServer(server, apiServer{})
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("Error creating listener: ", err.Error())
	}
	fmt.Println("Start serving grpc in port :5000")

	// for RestAPI request handler
	go requestHandle()
	server.Serve(listener)

}

func requestHandle() {
	const port = 8092
	fmt.Printf("Server is up and running at port %d ", port)
	myRoutes := mux.NewRouter()
	// myRoutes.HandleFunc("/api/v1/", routes.Home).Methods("GET")
	// myRoutes.HandleFunc("/api/v1/register/devices", routes.RegisterHandler).Methods("POST")
	// myRoutes.HandleFunc("/api/v1/get/devices", routes.GetDeviceInfoHandler).Methods("GET")

	myRoutes.HandleFunc("/", routes.Home).Methods("GET")

	err := http.ListenAndServe(":8092", myRoutes)
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
