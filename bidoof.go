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
	//I just want to return the error to the caller so I don't need/want to defer here... I think...
	insert.Close()
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

	retrySleep := func(message string) {
		log.Printf("Failed to connect to %s... will sleep on it and try again", message)
		time.Sleep(2 * time.Second)
	}

	conn, err := amqp.Dial(mqUrl)
	rabbitRetryCount := 0
	//When trying to run this with other containers, I noticed this would try to start
	//before the others were fully functional, so a few retries might do the trick here

	//Would be better to abstract this into a function
	//Needs to be tested as well, but my Go skills can't quite figure out how
	//to look for fatal's

	//Probably 100 better ways of doing this...
	for err != nil && rabbitRetryCount < 5 {
		rabbitRetryCount++
		retrySleep("RabbitMQ")
		conn, err = amqp.Dial(mqUrl)
	}
	helpers.FailOnError(err, "Failed to connect to RabbitMQ after several attempts")

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

	//Connect to DB - same comment as above for Rabbit: would like to abstract this into a function
	db, err := sql.Open("mysql", helpers.MakeDbConnectionString())
	dbRetryCount := 0
	for err != nil && dbRetryCount < 5 {
		dbRetryCount++
		retrySleep("DB")
		db, err = sql.Open("mysql", helpers.MakeDbConnectionString())
	}
	helpers.FailOnError(err, "Failed to open DB connection after several attempts")
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
