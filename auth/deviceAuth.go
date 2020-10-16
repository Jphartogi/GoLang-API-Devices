package auth

import (
	"api/database"
	"api/global"
	"log"

	"github.com/dgrijalva/jwt-go"
)

//RegisterDevice is a function to register a device
func RegisterDevice(device *global.NewDevice) (string, string, error) {
	id, err := database.AddDevicesToDB(device)
	if err != nil {
		log.Fatal(err)
	}

	Token, errs := DeviceTokenGenerator(device)

	if errs != nil {
		log.Fatal(err)
	}
	return id, Token, nil
}

//UpdateDevice function handler to update the device using ID
func UpdateDevice(device *global.NewDevice) (bool, error) {
	_, err := database.UpdateDeviceOnDB(device)
	if err != nil {
		log.Fatal(err)
	}
	return true, nil
}

// DeviceTokenGenerator is to generate a token for devices
func DeviceTokenGenerator(device *global.NewDevice) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	// expiredTime := time.Now().Add(time.Minute * 10)

	claims["authorized"] = true
	claims["device_name"] = device.DeviceName
	// claims["exp"] = expiredTime

	signedToken, err := token.SignedString(mySigningKey)

	if err != nil {
		panic(err)
	}

	return signedToken, nil
}
