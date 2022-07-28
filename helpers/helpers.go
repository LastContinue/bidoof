package helpers

import (
	"log"
	"time"
)

const RETRY_TIME = 10

type Registration struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Tour       string `json:"tour"`
	IslandType string `json:"islandType"`
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func RetrySleep(message string) {
	log.Printf("Failed to connect to %s... will sleep on it and try again", message)
	time.Sleep(RETRY_TIME * time.Second)
}
