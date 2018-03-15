package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

func CreateSession() *mgo.Session {
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{"localhost:27017"},
		Username: "",
		Password: "",
		Timeout:  60 * time.Second,
		Database: "waterfall",
	})

	if err != nil {
		log.Fatalf("[CreateSession] - error create session %s\n", err)
	}
	return session
}
