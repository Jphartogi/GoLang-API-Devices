package global

//NilUser is the nil value for user
var NilUser User

//User is a struct for user credentials
type User struct {
	Username string `bson:"userName"`
	Password string `bson:"password"`
}
