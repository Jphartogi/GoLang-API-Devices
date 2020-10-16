package global

import "time"

// Message is a token struct standard message
type Message struct {
	Message    string
	DeviceInfo NewDevice
	Token      string
	CreatedAt  time.Time
	ExpiredAt  time.Time
}

// NewDevice struct data type for creating device
type NewDevice struct {
	DeviceID       string    `bson:"deviceId"`
	DeviceName     string    `bson:"deviceName"`
	DeviceCategory string    `bson:"deviceCategory"`
	DeviceLocation string    `bson:"deviceLocation"`
	CreatedAt      time.Time `bson:"createdAt"`
}

//DeviceList is a struct for device main information
type DeviceList struct {
	DeviceID       string `bson:"deviceId"`
	DeviceName     string `bson:"deviceName"`
	DeviceCategory string `bson:"deviceCategory"`
	DeviceLocation string `bson:"deviceLocation"`
}
