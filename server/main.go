package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Transaction struct {
	Payer     string `json:"payer"`
	Points    int    `json:"points"`
	Timestamp string `json:"timestamp"`
}

func main() {
	fmt.Print("Attempting to run server")

	// declare array of points objects
	transactions := []Transaction{}

	// initialize app
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// health check
	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Post("/api/add-transaction", func(c *fiber.Ctx) error {
		// instantiate transaction object
		transaction := &Transaction{}

		// check transaction for errors
		if err := c.BodyParser(transaction); err != nil {
			c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
			return err
		}

		// destructure payload into struct
		payload := struct {
			Payer  string `json:"payer"`
			Points int    `json:"points"`
		}{}

		// parse request body into payload or return error
		if err := c.BodyParser(&payload); err != nil {
			c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
			return err
		}

		// set transaction data points from payload
		transaction.Payer = payload.Payer
		transaction.Points = payload.Points
		transaction.Timestamp = time.Now().String()

		// append the new transaction to the transactions array
		transactions = append(transactions, *transaction)

		// return list of transactions
		return c.JSON(transactions)
	})

	// spend points

	// get balance
	app.Get("/api/get-balance", func(c *fiber.Ctx) error {
		return c.JSON(transactions)
	})

	// listen on port 8081
	log.Fatal(app.Listen(":8081"))
}
