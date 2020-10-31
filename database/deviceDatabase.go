package database

import (
	"api/global"
	"context"
	"errors"
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
		return dlist, errors.New("Failed to find devices " + err.Error())
	}

	for cur.Next(context.TODO()) {
		var x global.DeviceList
		err = cur.Decode(&x)
		if err != nil {
			return dlist, err
		}
		dlist = append(dlist, &x)
	}
	return dlist, nil
}

//AddDevicesToDB is a function to add new data to database
func AddDevicesToDB(device *global.NewDevice) (string, error) {
	uuid := UUIDGenerator()
	device.DeviceID = uuid.String()
	device.CreatedAt = time.Now().Local()
	_, err := DB.Collection("devices").InsertOne(context.Background(), *device)
	if err != nil {
		return "", errors.New("Error on inserting new device to DB")
	}
	return uuid.String(), nil
}

//UpdateDeviceOnDB function to update the device on DB
func UpdateDeviceOnDB(device *global.DeviceList) (bool, error) {
	result, status := checkDeviceIDExist(device)
	if !status {
		return false, errors.New("No Device Found")
	}
	if device.DeviceCategory == "" {
		device.DeviceCategory = result.DeviceCategory
	}
	if device.DeviceName == "" {
		device.DeviceName = result.DeviceName
	}
	if device.DeviceLocation == "" {
		device.DeviceLocation = result.DeviceLocation
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel() // releases resources if slowOperation completes before timeout elapses
	updateItems := bson.M{
		"$set": bson.M{
			"deviceName":     device.DeviceName,
			"deviceCategory": device.DeviceCategory,
			"deviceLocation": device.DeviceLocation,
			"userName":       device.Username},
	}
	_, err := DB.Collection("devices").UpdateOne(
		ctx,
		bson.M{"deviceId": device.DeviceID},
		updateItems)

	if err != nil {
		return false, errors.New("Failed to update device " + err.Error())
	}
	return true, nil
}

//StoreDeviceDataToDatabase is a function to store device data to database
func StoreDeviceDataToDatabase(data *global.DeviceData) error {
	data.TimeStamp = time.Now().Local()
	_, err := DB.Collection("deviceData").InsertOne(context.Background(), *data)
	if err != nil {
		return err
	}
	return nil
}

// //GetListOfDeviceByUsername function to access database to search all device registered to user
// func GetListOfDeviceByUsername(data *global.DeviceDataSearch) (global.DataSearchResult, error) {

// }

func checkDeviceIDExist(device *global.DeviceList) (global.DeviceList, bool) {
	var deviceInfo global.DeviceList
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel() // releases resources if slowOperation completes before timeout elapses
	err := DB.Collection("devices").FindOne(ctx, bson.M{"deviceId": device.DeviceID}).Decode(&deviceInfo)
	if err != nil {
		return deviceInfo, false
	}

	return deviceInfo, true
}

// UUIDGenerator is a func who generate new UUID for new devices
func UUIDGenerator() uuid.UUID {
	uuid := uuid.NewV4()
	return uuid
}
