package database

import (
	"api/auth"
	"context"
	"log"
	"time"

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
func AddToDB(device *auth.NewDevice) error {
	uuid := auth.UUIDGenerator()
	device.DeviceID = uuid.String()
	device.CreatedAt = time.Now().Local()
	_, err := DB.Collection("devices").InsertOne(context.Background(), *device)
	if err != nil {
		log.Fatalln("Error on inserting new devices", err)
	}
	return nil
}
