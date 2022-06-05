package main

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"

	"database/sql"

	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load()
	FailOnError(err, "Error loading .env file")

	var myEnv map[string]string
	myEnv, err = godotenv.Read()
	FailOnError(err, "Error loading .env file")

	db, err := sql.Open("postgres", myEnv["POSTGRESQL_CONNECTION_URL"])
	FailOnError(err, "Error connecting to database")
	defer db.Close()

	conn, err := amqp.Dial(myEnv["RABBITMQ_CONNECTION_URL"])
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		myEnv["RABBITMQ_CHANNEL"],
		false,
		false,
		false,
		false,
		nil,
	)
	FailOnError(err, "Failed to declare a queue")

	ctx := context.Background()
	customsearchService, err := customsearch.NewService(ctx, option.WithAPIKey(myEnv["GOOGLE_CUSTOM_SERACH_API_KEY"]))
	FailOnError(err, "Failed to create custom search service")

	// launching checker in another go routine
	go Checker(ch, q, db, customsearchService, myEnv["GOOGLE_CUSTOM_SEARCH_ENGINE_ID"])
	Server(ch, q, db)
}
