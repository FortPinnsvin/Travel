package models

import (
	"labix.org/v2/mgo"
	"os"
)

type User struct {
	Id        string `bson:"_id,omitempty"`
	FirstName string
	LastName  string
	Email     string
}

func ConnectToDataBase(){
	url := os.Getenv("DB_URL")
	if url == "" {
		url = "localhost"
	}
	database := os.Getenv("DB")
	if database == "" {
		database = "travel"
	}

	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	UserCollection = session.DB(database).C("users")

}