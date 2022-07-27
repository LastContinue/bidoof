package main

import (
	"bidoof/helpers"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/streadway/amqp"
)

//Not sure if this is correct syntax but idea is to keep "magic" strings down
const NA_TYPE = "NA"

type Registration struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Tour       string `json:"tour"`
	IslandType string `json:"islandType"`
}

func jsonToRegistration(jsonMsg string) (Registration, error) {
	var registrationMessage Registration
	err := json.Unmarshal([]byte(jsonMsg), &registrationMessage)
	return registrationMessage, err
}

func persistRegistration(registration *Registration, db *sql.DB) (*sql.Rows, error) {
	//Ideally we'd have some sort of abstraction or ORM in front of this so end-engineers wouldn't need
	//to mess around with raw SQL
	insert, err := db.Query("INSERT INTO attendees(name, email, tour) VALUES(?,?,?)", registration.Name, registration.Email, registration.Tour)
	defer insert.Close()
	return insert, err
}

func formatEmailBody(registration *Registration) string {
	emailBody := fmt.Sprintf("Thank you for registering to Illumicon, %s!", registration.Name)

	if registration.IslandType != NA_TYPE {
		emailBody += fmt.Sprintf("\n For your tour to %s, please pack appropriately as it is a %s-y place!", registration.Tour, registration.IslandType)
	}

	return emailBody
}

func sendRegistrationEmail(registration *Registration) {
	//Here's where we'd actually send an email
	emailBody := formatEmailBody(registration)

	log.Print(emailBody)
	//since email can be time-intenstive depending on the complexity of the attachements, etc
	//simulate some time to justify the go-routine
	time.Sleep(30 * time.Second)
}

func main() {
	//Connect to RabbitMQ
	mqUrl := os.Getenv("MESSAGE_QUEUE_URL")
	mqQueueName := os.Getenv("MESSAGE_QUEUE_NAME")

	conn, err := amqp.Dial(mqUrl)
	helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	helpers.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	messages, err := ch.Consume(
		mqQueueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	helpers.FailOnError(err, "Failed to register a consumer")

	//Connect to DB
	db, err := sql.Open("mysql", helpers.MakeDbConnectionString())
	helpers.FailOnError(err, "Failed to open DB connection")
	defer db.Close()

	go func() {
		for message := range messages {
			//again, not a great idea to be exposing user data including logging - but good for demo'ing
			log.Printf("Received a message: %s", message.Body)

			registration, err := jsonToRegistration(string(message.Body))
			helpers.FailOnError(err, "Could not unmarshall JSON")

			_, err = persistRegistration(&registration, db)
			helpers.FailOnError(err, "Failed to insert to MySql")

			//sending email can be time intensive so make this its own thread
			go sendRegistrationEmail(&registration)
		}
	}()

	log.Printf(" *** Waiting for messages. To exit press CTRL+C")
	forever := make(chan bool)
	<-forever
}
