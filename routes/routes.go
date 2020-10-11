package routes

import (
	"api/helpers"
	"api/token"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
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

// NewDevice struct data type for creating device
type NewDevice struct {
	DeviceID  uuid.UUID
	CreatedAt time.Time
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

// TokenHandler function to handle request in token realm
func TokenHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var x token.Device
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

	if x.DeviceName == "" {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	signedToken, expiredTime, err := token.Generator(&x)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	data := token.TokenMessage{
		Message:    "Successfully generated Token",
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

// UUIDHandler is a func who generate new UUID for new devices
func UUIDHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	response := NewDevice{
		DeviceID:  uuid,
		CreatedAt: time.Now().Local()}

	UUIDResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(UUIDResponse)

}
