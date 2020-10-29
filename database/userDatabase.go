package database

import (
	"api/global"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// dburi       = "mongodb://localhost:27017"
	// uri for docker integration
	dburi       = "mongodb://database:27017"
	dbname      = "api"
	testDBname  = "testAPI"
	performance = 100
)

// DB holds Database Connection
var DB mongo.Database

//TestDB holds test Database Connection
var TestDB mongo.Database

// ConnecttoDB is a function to connect to DB
func ConnecttoDB() {
	ctx, cancel := NewDBContext(10 * time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal("Error connect to DB: ", err.Error())
	}
	DB = *client.Database(dbname)
}

//ConnectToTestDB is a test connection function to testDB
func ConnectToTestDB() {
	ctx, cancel := NewDBContext(10 * time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal("Error connect to DB: ", err.Error())
	}
	TestDB = *client.Database(testDBname)
}

// NewDBContext returns a new Context according to app performance
func NewDBContext(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d*performance/100)
}

//RetrieveUserData function to get the user data from Database
func RetrieveUserData(filter bson.M) (global.User, error) {
	var user global.User
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel() // releases resources if slowOperation completes before timeout elapses
	err := DB.Collection("user").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return user, errors.New("Wrong login or password")
	}
	return user, nil
}

//AddUsertoDB is a function to add new user to database
func AddUsertoDB(user *global.User) error {
	_, err := DB.Collection("user").InsertOne(context.Background(), *user)
	if err != nil {
		log.Fatalln("Error on inserting new user", err)
	}
	return nil
}

//CheckUsernameAvailability is a function to check whether the username is still available
func CheckUsernameAvailability(user *global.User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel() // releases resources if slowOperation completes before timeout elapses
	err := DB.Collection("user").FindOne(ctx, bson.M{"userName": user.Username}).Decode(&user)
	if err != nil {
		// if no same username is return, availability is true
		return true
	}
	return false
}

func checkUsernameExist(user *global.User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel() // releases resources if slowOperation completes before timeout elapses
	err := DB.Collection("user").FindOne(ctx, bson.M{"userName": user.Username}).Decode(&user)
	if err != nil {
		// if no same username is return, username is not yet exist, returning false
		return false
	}
	return true
}

// DeleteUserFromDB is a function to delete a single user
func DeleteUserFromDB(user *global.User) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel() // releases resources if slowOperation completes before timeout elapses

	// First check if user exist
	result := checkUsernameExist(user)
	if !result {
		return false, errors.New("Username is not yet exists")
	}
	_, err := DB.Collection("user").DeleteOne(ctx, bson.M{"userName": user.Username})
	if err != nil {
		return false, err
	}

	return true, nil
}
