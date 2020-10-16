package auth

import (
	"api/global"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// DeviceTokenGenerator is to generate a token for devices
func DeviceTokenGenerator(device *global.NewDevice) (string, time.Time, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	expiredTime := time.Now().Add(time.Minute * 10)

	claims["authorized"] = true
	claims["device_name"] = device.DeviceName
	claims["username"] = device.Username
	claims["exp"] = expiredTime

	signedToken, err := token.SignedString(mySigningKey)

	if err != nil {
		panic(err)
	}

	return signedToken, expiredTime, nil
}
