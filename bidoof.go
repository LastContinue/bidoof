package main

import (
	database "bidoof/db"
	"bidoof/helpers"
	h "bidoof/helpers"
	"bidoof/messagequeue"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

//Not sure if this is correct syntax but idea is to keep "magic" strings down
const NA_TYPE = "NA"

func jsonToRegistration(jsonMsg string) (h.Registration, error) {
	var registrationMessage h.Registration
	err := json.Unmarshal([]byte(jsonMsg), &registrationMessage)
	return registrationMessage, err
}

func formatEmailBody(registration *h.Registration) string {
	emailBody := fmt.Sprintf("Thank you for registering to Illumicon, %s!", registration.Name)

	if registration.IslandType != NA_TYPE {
		emailBody += fmt.Sprintf("\n For your tour to %s, please pack appropriately as it is a %s-y place!", registration.Tour, registration.IslandType)
	}

	return emailBody
}

func sendRegistrationEmail(registration *h.Registration) {
	//Here's where we'd actually send an email
	emailBody := formatEmailBody(registration)

	log.Print(emailBody)
	//since email can be time-intenstive depending on the complexity of the attachements, etc
	//simulate some time to justify the go-routine
	time.Sleep(30 * time.Second)
}

func main() {
	conn, err := messagequeue.Connect()
	h.FailOnError(err, "Failed to connect to RabbitMQ after several attempts")
	defer conn.Close()

	channel, err := messagequeue.Channel(conn)
	helpers.FailOnError(err, "Failed to open a channel")
	defer channel.Close()

	messages, err := messagequeue.Consumer(channel)
	h.FailOnError(err, "Failed to create RabbitMQ consumer")

	db, err := database.Connect()
	h.FailOnError(err, "Failed to open DB connection after several attempts")
	defer db.Close()

	go func() {
		for message := range messages {
			//again, not a great idea to be exposing user data including logging - but good for demo'ing
			log.Printf("Received a message: %s", message.Body)

			registration, err := jsonToRegistration(string(message.Body))
			h.FailOnError(err, "Could not unmarshall JSON")

			_, err = database.PersistRegistration(&registration, db)
			h.FailOnError(err, "Failed to insert to MySql")

			//sending email can be time intensive so make this its own thread
			go sendRegistrationEmail(&registration)
		}
	}()

	log.Printf(" *** Waiting for messages. To exit press CTRL+C")
	forever := make(chan bool)
	<-forever
}
