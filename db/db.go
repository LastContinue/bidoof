package db

import (
	"bidoof/helpers"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

const RETRY_COUNT = 3

//this syntax will make sense for mocking later - probably 100 better ways to do this
var makeDbConnectionString = func() string {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUrl := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")

	return fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbUrl, dbName)
}

func Connect() (*sql.DB, error) {
	//Connect to DB - same comment as above for Rabbit: would like to abstract this into a function
	db, err := sql.Open("mysql", makeDbConnectionString())
	dbRetryCount := 0
	for err != nil && dbRetryCount < RETRY_COUNT {
		dbRetryCount++
		helpers.RetrySleep("DB")
		db, err = sql.Open("mysql", makeDbConnectionString())
	}

	return db, err
}

func PersistRegistration(registration *helpers.Registration, db *sql.DB) (*sql.Rows, error) {
	//Ideally we'd have some sort of abstraction or ORM in front of this so end-engineers wouldn't need
	//to mess around with raw SQL
	insert, err := db.Query("INSERT INTO attendees(name, email, tour) VALUES(?,?,?)", registration.Name, registration.Email, registration.Tour)
	//I just want to return the error to the caller so I don't need/want to defer here... I think...
	insert.Close()
	return insert, err
}
