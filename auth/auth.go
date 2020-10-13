package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

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
	DeviceID   string    `bson:"deviceId"`
	DeviceName string    `bson:"deviceName"`
	Username   string    `bson:"userName"`
	Email      string    `bson:"email"`
	Password   string    `bson:"password"`
	Latitude   float32   `bson:"lat"`
	Longitude  float32   `bson:"long"`
	CreatedAt  time.Time `bson:"createdAt"`
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

// TokenGenerator is to generate a Token to be consumed for API usage
func TokenGenerator(device *NewDevice) (string, time.Time, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	expiredTime := time.Now().Add(time.Minute * 10)

	claims["authorized"] = true
	claims["device_name"] = device.DeviceName
	claims["username"] = device.Username
	claims["password"] = device.Password
	claims["exp"] = expiredTime

	signedToken, err := token.SignedString(mySigningKey)

	if err != nil {
		panic(err)
	}

	return signedToken, expiredTime, nil
}

// UUIDGenerator is a func who generate new UUID for new devices
func UUIDGenerator() uuid.UUID {
	uuid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return uuid

}

// HashPassword is a function to hashed the password of a registered device
func HashPassword(device *NewDevice) ([]byte, string, error) {
	pass := &device.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*pass), 10)
	if err != nil {
		panic(err)
	}
	return hashedPassword, string(hashedPassword[:]), nil
}
