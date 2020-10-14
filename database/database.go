package database

import (
	"api/global"
	"context"
	"errors"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi       = "mongodb://localhost:27017"
	dbname      = "api"
	performance = 100
)

// DB holds Database Connection
var DB mongo.Database

func connectToDB() {
	ctx, cancel := NewDBContext(10 * time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal("Error connect to DB: ", err.Error())
	}
	DB = *client.Database(dbname)
}

// NewDBContext returns a new Context according to app performance
func NewDBContext(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d*performance/100)
}

//AddToDB is a function to add new data to database
func AddToDB(device *global.NewDevice) error {
	uuid := UUIDGenerator()
	device.DeviceID = uuid.String()
	device.CreatedAt = time.Now().Local()
	_, err := DB.Collection("devices").InsertOne(context.Background(), *device)
	if err != nil {
		log.Fatalln("Error on inserting new devices", err)
	}
	return nil
}

// UUIDGenerator is a func who generate new UUID for new devices
func UUIDGenerator() uuid.UUID {
	uuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return uuid

}

//RetrieveDeviceData function to get the data from Database and implement filter
func RetrieveDeviceData(filter bson.M) ([]*global.DeviceList, error) {
	var dlist []*global.DeviceList
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel() // releases resources if slowOperation completes before timeout elapses
	cur, err := DB.Collection("devices").Find(ctx, filter)
	if err != nil {
		log.Fatal("Error on finding documents", err)
	}

	for cur.Next(context.TODO()) {
		var x global.DeviceList
		err = cur.Decode(&x)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		dlist = append(dlist, &x)
	}
	return dlist, nil
}

//RetrieveUserData function to get the user data from Database
func RetrieveUserData(filter bson.M) (global.User, error) {
	var user global.User
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel() // releases resources if slowOperation completes before timeout elapses
	err := DB.Collection("devices").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return user, errors.New("Wrong login or password")
	}
	return user, nil
}
