package database

import (
	"api/global"
	"context"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

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

//AddDevicesToDB is a function to add new data to database
func AddDevicesToDB(device *global.NewDevice) error {
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
