package auth

import (
	"api/database"
	"api/global"
	"errors"
	"fmt"
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

var mySigningKey = []byte(goDotEnvVariable("SECRET_KEY"))

//StoreDeviceData is a middleware to store device data and auth the token first
func StoreDeviceData(data *global.DeviceAuth) (bool, error) {
	auth, err := validateToken(data.Token)
	if err != nil {
		log.Fatal(err)
	}
	if !auth {
		return false, errors.New("Authentication Failed")
	}

	errs := database.StoreDataToDatabase(&data.Data)
	if errs != nil {
		log.Fatal(errs)
	}

	return true, nil
}

func validateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return mySigningKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["device_name"])
		return true, nil
	}

	fmt.Println(err)
	return false, err

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
