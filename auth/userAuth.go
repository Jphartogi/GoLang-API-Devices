package auth

import (
	"api/database"
	"api/global"
	"errors"
	"fmt"
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

var mySigningKeyUser = []byte(goDotEnvVariable("SECRET_KEY_USER"))

//UserTokenGenerator is to generate token for user
func UserTokenGenerator(user *global.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	expiredTime := time.Now().Add(time.Minute * 10)

	claims["authorized"] = true
	claims["username"] = user.Username
	claims["password"] = user.Password
	claims["exp"] = expiredTime

	signedToken, err := token.SignedString(mySigningKeyUser)

	if err != nil {
		panic(err)
	}

	return signedToken, err
}

func validateUserToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return mySigningKeyUser, nil
	})

	if err != nil {
		return false, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}

	return false, err

}

// HashPassword is a function to hashed the password of a registered device
func HashPassword(user *global.User) ([]byte, string, error) {
	pass := &user.Password
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
		return false, errors.New(err.Error())
	}
	if bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password)) != nil {
		return false, errors.New("Wrong credentials")
	}

	return true, nil
}

//RegisterUser is handler function to add user and password to DB
func RegisterUser(user *global.User) (bool, error) {
	var newUser global.User
	_, stringHashedPassword, err := HashPassword(user)
	if err != nil {
		return false, errors.New("Failed to register user " + err.Error())
	}
	newUser.Username = user.Username
	newUser.Password = stringHashedPassword
	e := database.AddUsertoDB(&newUser)
	if e != nil {
		return false, errors.New("Failed to add user to DB" + err.Error())
	}
	return true, nil
}

//DeleteUser is a handler function to delete user from DB
func DeleteUser(user *global.User) (bool, error) {
	var newUser global.User
	newUser.Username = user.Username
	res, e := database.DeleteUserFromDB(&newUser)
	if e != nil {
		return res, e
	}
	return res, nil
}
