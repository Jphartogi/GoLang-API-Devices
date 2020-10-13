package routes

import (
	"api/auth"
	"api/database"
	"api/helpers"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// User is struct type for anyone who use the API
type User struct {
	Name        string    `json:"name"`
	IsMasked    bool      `json:"is_masked"`
	Image       string    `json:"image"`
	Temperature float32   `json:"temperature"`
	CreatedAt   time.Time `json:"created_at"`
}

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

// RegisterHandler function to handle registering new devices
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Check if method is not correct
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var x auth.NewDevice
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

	if x.DeviceName == "" || x.Username == "" || x.Email == "" || x.Password == "" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	signedToken, expiredTime, err := auth.TokenGenerator(&x)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	//Password hashing
	_, pass, err := auth.HashPassword(&x)

	//Update password
	x.Password = pass

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	// Add file to database
	errors := database.AddToDB(&x)

	if errors != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	data := auth.Message{
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
