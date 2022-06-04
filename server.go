package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func handleSubmit(db *sql.DB, ch *amqp.Channel, q amqp.Queue) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if form, err := c.MultipartForm(); err == nil {
			files := form.File["file"]
			if len(files) == 0 {
				return c.Redirect("/error")
			}
			if err := c.SaveFile(files[0], fmt.Sprintf("./files/%s", files[0].Filename)); err != nil {
				return err
			}
			jobId := uuid.New().String()

			_, err := db.Query("INSERT INTO jobs (id, location, status) VALUES ($1, $2, $3)", jobId, files[0].Filename, 0)

			if err != nil {
				log.Fatal(err)
				return c.Redirect("/error")
			}

			body := fmt.Sprintf("%s$%s", files[0].Filename, jobId)
			err = ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			FailOnError(err, "Failed to publish a message")
			return c.Redirect("/jobs?id=" + jobId)
		}
		return c.Redirect("/error")
	}
}

func handleGetJobStatus(db *sql.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		jobId := c.Params("jobId")
		var status int
		err := db.QueryRow("SELECT status FROM jobs WHERE id = $1", jobId).Scan(&status)
		if err != nil {
			return c.SendStatus(500)
		}
		var link string
		var title string
		var description string
		var similarity int
		var result []fiber.Map
		refs, err := db.Query("SELECT link, title, description, similarity FROM refs WHERE jobId = $1", jobId)
		if err != nil {
			return c.SendStatus(500)
		}
		for refs.Next() {
			err := refs.Scan(&link, &title, &description, &similarity)
			if err != nil {
				return c.SendStatus(500)
			}
			obj := fiber.Map{
				"title":       title,
				"description": description,
				"similarity":  similarity,
				"link":        link,
			}
			result = append(result, obj)

		}
		return c.JSON(fiber.Map{
			"status": status,
			"result": result,
		})
	}
}

func Server(ch *amqp.Channel, q amqp.Queue, db *sql.DB) {
	app := fiber.New()
	app.Static("/", "./static")
	app.Post("/", handleSubmit(db, ch, q))
	app.Get("/status/:jobId", handleGetJobStatus(db))
	app.Use(func(c *fiber.Ctx) error {
		return c.Redirect("/404.html")
	})
	log.Printf("Server started")
	app.Listen(":8080")
}
