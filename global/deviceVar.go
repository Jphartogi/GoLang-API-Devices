package global

import "time"

// DeviceDataSearch is a token struct standard message
type DeviceDataSearch struct {
	Username string
	Token    string
}

// NewDevice struct data type for creating device
type NewDevice struct {
	DeviceID       string `bson:"deviceId"`
	DeviceName     string `bson:"deviceName"`
	DeviceCategory string `bson:"deviceCategory"`
	DeviceLocation string `bson:"deviceLocation"`
	Username       string `bson:"userName"`
	UserToken      string
	CreatedAt      time.Time `bson:"createdAt"`
}

//DeviceList is a struct for device main information
type DeviceList struct {
	DeviceID       string `bson:"deviceId"`
	DeviceName     string `bson:"deviceName"`
	DeviceCategory string `bson:"deviceCategory"`
	DeviceLocation string `bson:"deviceLocation"`
	Username       string `bson:"userName"`
	Token          string
}

//DeviceAuth is a struct for device data and its token
type DeviceAuth struct {
	Data  DeviceData
	Token string
}

//DeviceData is a struct for device data
type DeviceData struct {
	DeviceID    string    `bson:"deviceId"`
	DeviceValue int       `bson:"deviceValue"`
	TimeStamp   time.Time `bson:"timestamp"`
}
