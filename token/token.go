package token

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

// TokenMessage is a token struct standard message
type TokenMessage struct {
	Message    string
	DeviceInfo Device
	Token      string
	CreatedAt  time.Time
	ExpiredAt  time.Time
}

// Device struct data type for authorizing device for token
type Device struct {
	DeviceID   uuid.UUID
	DeviceName string
}

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load()

	if err != nil {
		panic(err)
		// log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

var mySigningKey = []byte(goDotEnvVariable("SECRET_KEY"))

// var mySigningKey = []byte("mysecretkey")

// Generator is to generate a Token to be consumed for API usage
func Generator(device *Device) (string, time.Time, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	expiredTime := time.Now().Add(time.Minute * 10)

	claims["authorized"] = true
	claims["id"] = device.DeviceID
	claims["user"] = device.DeviceName
	claims["exp"] = expiredTime

	signedToken, err := token.SignedString(mySigningKey)

	if err != nil {
		panic(err)
	}

	return signedToken, expiredTime, nil
}
