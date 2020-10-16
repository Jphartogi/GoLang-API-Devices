package main

import (
	"api/database"
	"api/global"
	"api/proto"
	"context"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func Test_apiServer_Login(t *testing.T) {
	database.ConnecttoDB()
	pw, _ := bcrypt.GenerateFromPassword([]byte("example"), bcrypt.DefaultCost)
	database.DB.Collection("user").InsertOne(context.Background(), global.User{Username: "test", Password: string(pw)})
	server := apiServer{}
	_, err := server.Login(context.Background(), &proto.UserLoginRequest{Username: "test", Password: "example"})
	if err != nil {
		t.Error("1: An error was returned: ", err.Error())
	}
}

func Test_apiServer_SignUp(t *testing.T) {
	var err error
	database.ConnecttoDB()
	database.DB.Collection("user").InsertOne(context.Background(), global.User{Username: "Ronaldo"})
	server := apiServer{}
	_, err = server.SignUp(context.Background(), &proto.UserSignUpRequest{Username: "Justin", Password: "phartogi"})
	if err != nil {
		t.Error("1: An error was returned: ", err.Error())
	}
	_, err = server.SignUp(context.Background(), &proto.UserSignUpRequest{Username: "Ronaldo", Password: "phartogi"})
	if err == nil {
		t.Error("2: An error was returned: ", err.Error())
	}
}

func Test_apiServer_UsernameTaken(t *testing.T) {
	database.ConnecttoDB()
	database.DB.Collection("user").InsertOne(context.Background(), global.User{Username: "Messi"})
	server := apiServer{}
	res, err := server.UsernameTaken(context.Background(), &proto.UsernameTakenRequest{Username: "Joshua P"})
	if err != nil {
		t.Error("1: An error was returned: ", err.Error())
	}
	if !res.Status {
		t.Error("1: Wrong result")
	}
	res, err = server.UsernameTaken(context.Background(), &proto.UsernameTakenRequest{Username: "Messi"})
	if err != nil {
		t.Error("2: An error was returned: ", err.Error())
	}
	if res.Status {
		t.Error("2: Wrong result")
	}
}

func Test_apiServer_DeleteUser(t *testing.T) {
	database.ConnecttoDB()
	database.DB.Collection("user").InsertOne(context.Background(), global.User{Username: "Joshua"})
	server := apiServer{}
	_, err := server.DeleteUser(context.Background(), &proto.UserDeleteRequest{Username: "Phartogi"})
	if err == nil {
		t.Error("1: An error was returned: ", err.Error())
	}
	_, errs := server.DeleteUser(context.Background(), &proto.UserDeleteRequest{Username: "Joshua"})
	if errs != nil {
		t.Error("2: An error was returned: ", err.Error())
	}

}
