package database

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

// Mongo session
var Mongo *mgo.Session
var databases MongoDB

// MongoDB configuration
type MongoDB struct {
	URL      string
	Database string
}

// Connect to the database
func Connect(d MongoDB) {
	var err error

	databases = d

	// Connect to MongoDB
	log.Printf("Connecting to %s", d.URL)
	if Mongo, err = mgo.DialWithTimeout(d.URL, 5*time.Second); err != nil {
		log.Println("MongoDB Driver Error", err)
		return
	}
	log.Printf("connected to mongo")

	// Prevents these errors: read tcp 127.0.0.1:27017: i/o timeout
	Mongo.SetSocketTimeout(1 * time.Second)

	// Check if is alive
	if err = Mongo.Ping(); err != nil {
		log.Println("Database Error", err)
	}
}

// CheckConnection returns true if MongoDB is available
func CheckConnection() bool {
	if Mongo == nil {
		Connect(databases)
	}

	if Mongo != nil {
		return true
	}

	return false
}

// ReadConfig returns the database information
func ReadConfig() MongoDB {
	return databases
}
