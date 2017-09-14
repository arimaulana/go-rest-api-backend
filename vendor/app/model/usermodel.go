package model

import (
	"app/shared/database"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// User structurized the user model
type User struct {
	UserName  string    `json:"username,omitempty" bson:"username,omitempty"`
	FirstName string    `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email     string    `json:"email,omitempty" bson:"email,omitempty"`
	Password  string    `json:"password,omitempty" bson:"password,omitempty"`
	StatusID  uint8     `json:"status_id,omitempty" bson:"status_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Deleted   uint8     `json:"deleted,omitempty" bson:"deleted,omitempty"`
}

// GetUser gets user information from email or username
func GetUser(usernameoremail string) (User, error) {
	var err error

	result := User{}

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().Database).C("user")
		if strings.Contains(usernameoremail, "@") {
			err = c.Find(bson.M{"email": usernameoremail}).One(&result)
		} else {
			err = c.Find(bson.M{"username": usernameoremail}).One(&result)
		}
	} else {
		err = ErrUnavailable
	}

	return result, standardizeError(err)
}

// CreateUser registering user
func CreateUser(username, first, last, email, hashpass string) error {
	var err error

	now := time.Now()

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().Database).C("user")
		user := &User{
			UserName:  username,
			FirstName: first,
			LastName:  last,
			Email:     email,
			Password:  hashpass,
			StatusID:  1,
			CreatedAt: now,
			UpdatedAt: now,
			Deleted:   0,
		}
		err = c.Insert(user)
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}

// UpdateUser : updating user
func UpdateUser(username, first, last, email, hashpass string) error {
	var err error

	now := time.Now()

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().Database).C("user")
		var user User
		user, err = GetUser(username)
		if err == nil {
			user.UpdatedAt = now
			user.FirstName = first
			user.LastName = last
			user.Email = email
			user.Password = hashpass
			err = c.Update(bson.M{"username": username}, &user)
		} else {
			err = ErrUnauthorized
		}
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}
