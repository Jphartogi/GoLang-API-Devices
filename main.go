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
		return &proto.UserAuthResponse{}, err
	}
	token, err := auth.UserTokenGenerator(&user)
	return &proto.UserAuthResponse{UserToken: token}, nil

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
		return &proto.UserSuccessRegister{}, err
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
	id, name, location, category, username, token := input.GetDeviceID(), input.GetDeviceName(), input.GetDeviceLocation(), input.GetDeviceCategory(), input.GetUsername(), input.GetDeviceToken()
	if id == "" || username == "" || token == "" {
		return &proto.DeviceSuccessUpdate{}, errors.New("Please specify and fill the required parameter")
	}
	//TODO search in database first the initial value, so when it is empty, get the initial value
	var devices global.DeviceList
	devices.DeviceID = id
	devices.DeviceName = name
	devices.DeviceLocation = location
	devices.DeviceCategory = category
	devices.Username = username
	devices.Token = token
	_, err := auth.UpdateDevice(&devices)
	if err != nil {
		return &proto.DeviceSuccessUpdate{}, errors.New("Error in updating the device " + err.Error())
	}
	return &proto.DeviceSuccessUpdate{DeviceID: id}, nil
}

func (apiServer) RegisterDevice(_ context.Context, input *proto.DeviceRequest) (*proto.DeviceSuccessRegister, error) {
	deviceName, deviceCategory, deviceLocation, username, token := input.GetDeviceName(), input.GetDeviceCategory(), input.GetDeviceLocation(), input.GetUsername(), input.GetUserToken()
	if deviceName == "" || deviceCategory == "" || deviceLocation == "" || username == "" {
		return &proto.DeviceSuccessRegister{}, errors.New("Please insert all the required field")
	}
	if token == "" {
		return &proto.DeviceSuccessRegister{}, errors.New("Please insert your user token")
	}
	var devices global.NewDevice
	devices.DeviceName = deviceName
	devices.DeviceLocation = deviceLocation
	devices.DeviceCategory = deviceCategory
	devices.Username = username
	devices.UserToken = token
	id, token, err := auth.RegisterDevice(&devices)
	if err != nil {
		return &proto.DeviceSuccessRegister{}, errors.New("Error in registering the device " + err.Error())
	}
	return &proto.DeviceSuccessRegister{DeviceID: id, DeviceToken: token}, nil
}

func (apiServer) GetDeviceData(_ context.Context, input *proto.DeviceGetDataRequest) (*proto.DeviceDataResponse, error) {
	username, token := input.GetUsername(), input.GetDeviceToken()
	if username == "" || token == "" {
		return &proto.DeviceDataResponse{}, errors.New("Please insert all the required field")
	}
	var dev global.DeviceDataSearch
	dev.Token = token
	dev.Username = username
	res, err := auth.GetDeviceData(&dev)
	if err != nil {
		return &proto.DeviceDataResponse{}, errors.New("Errors on getting the data" + err.Error())
	}

	var resultList []*proto.DeviceUpdateRequest
	for _, x := range res {
		var p proto.DeviceUpdateRequest

		p.DeviceID = x.DeviceID
		p.DeviceName = x.DeviceName
		p.DeviceCategory = x.DeviceCategory
		p.DeviceLocation = x.DeviceLocation
		p.Username = x.Username

		// dlist = append(dlist, &x)
		resultList = append(resultList, &p)

	}

	finalResult := proto.DeviceDataResponse{}
	finalResult.Data = resultList

	return &finalResult, nil
}

// func (apiServer) DeleteDevice(_ context.Context, input *proto.DeviceDeleteRequest) (*proto.DeviceSuccessDelete, error) {

// }

func main() {
	port := 9001
	address := fmt.Sprintf(":%d", port)
	server := grpc.NewServer()
	proto.RegisterAPIServicesServer(server, apiServer{})
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Error creating listener: ", err.Error())
	}

	// for RestAPI request handler
	go requestHandle()
	// go func() {
	log.Fatal("Serving gRPC: ", server.Serve(listener).Error())
	// }()

	// grpcWebServer := grpcweb.WrapServer(server)
	// httpServer := &http.Server{
	// 	Addr: ":9001",
	// 	Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		if r.ProtoMajor == 2 {
	// 			grpcWebServer.ServeHTTP(w, r)
	// 		} else {
	// 			w.Header().Set("Access-Control-Allow-Origin", "*")
	// 			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// 			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
	// 			w.Header().Set("grpc-status", "")
	// 			w.Header().Set("grpc-message", "")
	// 			if grpcWebServer.IsGrpcWebRequest(r) {
	// 				grpcWebServer.ServeHTTP(w, r)
	// 			}
	// 		}
	// 	}), &http2.Server{}),
	// }
	// log.Fatal("Serving Proxy: ", httpServer.ListenAndServe().Error())

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
