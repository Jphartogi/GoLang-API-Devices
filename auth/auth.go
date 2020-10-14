package auth

import (
	"api/database"
	"api/global"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

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
func TokenGenerator(device *global.NewDevice) (string, time.Time, error) {

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

// HashPassword is a function to hashed the password of a registered device
func HashPassword(device *global.NewDevice) ([]byte, string, error) {
	pass := &device.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*pass), 10)
	if err != nil {
		panic(err)
	}
	return hashedPassword, string(hashedPassword[:]), nil
}

//AuthorizeUser is function to authorize wheter the user has registered or not
func AuthorizeUser(user *global.User) (bool, error) {
	// username, password := user.Username, user.Password
	result, err := database.RetrieveUserData(bson.M{"userName": user.Username})
	if err != nil {
		log.Fatal(err)
	}
	if bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password)) != nil {
		return false, errors.New("Wrong credentials")
	}

	return true, nil
}
