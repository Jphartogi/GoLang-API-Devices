package auth

import (
	"api/database"
	"api/global"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
)

//RegisterDevice is a function to register a device
func RegisterDevice(device *global.NewDevice) (string, string, error) {
	auth, errs := validateUserToken(device.UserToken)
	if errs != nil {
		return "", "", errors.New("Token is invalid")
	}
	if !auth {
		return "", "", errors.New("Authentication Failed")
	}

	id, err := database.AddDevicesToDB(device)
	if err != nil {
		return "", "", errors.New("Failed to register device to DB" + err.Error())
	}

	Token, errs := DeviceTokenGenerator(device)

	if errs != nil {
		return "", "", errors.New("Failed to generate Token" + err.Error())
	}
	return id, Token, nil
}

//UpdateDevice function handler to update the device using its device ID
func UpdateDevice(device *global.DeviceList) (bool, error) {
	auth, errs := validateDeviceToken(device.Token)
	if errs != nil {
		return false, errors.New("Token is invalid")
	}
	if !auth {
		return false, errors.New("Authentication Failed")
	}

	_, err := database.UpdateDeviceOnDB(device)
	if err != nil {
		return false, err
	}
	return true, nil
}

//GetDeviceData is a function to retrieve device data by username
func GetDeviceData(device *global.DeviceDataSearch) ([]*global.DeviceList, error) {
	auth, err := validateDeviceToken(device.Token)
	if err != nil {
		return []*global.DeviceList{}, errors.New("Token is invalid")
	}
	if !auth {
		return []*global.DeviceList{}, errors.New("Authentication Failed")
	}

	result, errs := database.RetrieveDeviceData(bson.M{"userName": device.Username})
	if errs != nil {
		return []*global.DeviceList{}, errs

	}
	return result, nil
}

//StoreDeviceData is a middleware to store device data and auth the token first
func StoreDeviceData(data *global.DeviceAuth) (bool, error) {
	auth, err := validateDeviceToken(data.Token)
	if err != nil {
		return false, errors.New("Token is invalid")
	}
	if !auth {
		return false, errors.New("Authentication Failed")
	}

	errs := database.StoreDeviceDataToDatabase(&data.Data)
	if errs != nil {
		return false, errors.New("Failed to store data to DB " + err.Error())
	}

	return true, nil
}

/************************** Helper Function *********************/

var mySigningKeyDevice = []byte(goDotEnvVariable("SECRET_KEY_DEVICE"))

func validateDeviceToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return mySigningKeyDevice, nil
	})

	if err != nil {
		return false, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}

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

	signedToken, err := token.SignedString(mySigningKeyDevice)

	if err != nil {
		panic(err)
	}

	return signedToken, nil
}
