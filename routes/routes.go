package routes

import (
	"api/auth"
	"api/database"
	"api/global"
	"api/helpers"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/joho/godotenv"
)

// ResponseMessage is standard format for welcoming message
type ResponseMessage struct {
	Message string
}

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// Home function for displaying json respnse
func Home(w http.ResponseWriter, r *http.Request) {
	const welcomeMessage = "API is up and running!"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data := ResponseMessage{Message: welcomeMessage}
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		panic(err)
	} else {
		w.Write(jsonResponse)
	}
}

// GetDeviceInfoHandler for getting the list of devices from the user
func GetDeviceInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var user global.User

	err := helpers.DecodeJSONBody(w, r, &user)
	if err != nil {
		var mr *helpers.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	// Authorization part
	if user.Username == "" || user.Password == "" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	result, err := auth.AuthorizeUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if result {
		listDevice, err := database.RetrieveDeviceData(bson.M{"userName": user.Username})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		response, err := json.Marshal(listDevice)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	}
}

// RegisterHandler function to handle registering new devices
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Check if method is not correct
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var x global.NewDevice
	// Check if the request is already as requested
	err := helpers.DecodeJSONBody(w, r, &x)
	if err != nil {
		var mr *helpers.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if x.DeviceName == "" || x.Username == "" || x.Email == "" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	signedToken, expiredTime, err := auth.DeviceTokenGenerator(&x)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	// Add file to database
	errors := database.AddDevicesToDB(&x)

	if errors != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	data := global.Message{
		Message:    "Successfully registered devices",
		DeviceInfo: x,
		Token:      signedToken,
		CreatedAt:  time.Now().Local(),
		ExpiredAt:  expiredTime}

	response, err := json.Marshal(data)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
