package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Transaction struct {
	Payer     string    `json:"payer"`
	Points    int       `json:"points"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	fmt.Print("Attempting to run server")

	// declare array of transaction objects
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
		transaction.Timestamp = time.Now()

		// append the new transaction to the transactions array
		transactions = append(transactions, *transaction)

		// return confirmation message
		// return c.SendString("OPERATION SUCCESSFUL")
		return c.JSON(transactions)
	})

	// spend points
	app.Post("/api/spend-points", func(c *fiber.Ctx) error {
		// instantiate transaction object
		transaction := &Transaction{}

		// check transaction for errors
		if err := c.BodyParser(transaction); err != nil {
			c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
			return err
		}

		// declare result struct
		type Result struct {
			Payer  string `json:"payer"`
			Points int    `json:"points"`
		}
		// declare array of result objects
		results := []Result{}

		// destructure payload into struct
		payload := struct {
			Points int `json:"points"`
		}{}

		// parse request body into payload or return error
		if err := c.BodyParser(&payload); err != nil {
			c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
			return err
		}

		// check that user entered more than 0 points to spend
		if payload.Points < 0 {
			return c.SendString("You must enter a number of points that is greater than 0")
		}

		// check that user has enough points
		totalPoints := 0
		for _, val := range transactions {
			totalPoints += val.Points
		}
		if payload.Points > totalPoints {
			return c.SendString("You do not have enough points")
		}

		// sort transactions from oldest to newest
		sort.Slice(transactions, func(i, j int) bool {
			return transactions[i].Timestamp.Before(transactions[j].Timestamp)
		})

		// loop through the sorted transactions
		for i := 0; i < len(transactions); i++ {
			// if all payload points are spent, break out of the loop
			if payload.Points == 0 {
				break
			}

			// if a spend record is detected, skip to the next iteration
			points := transactions[i].Points
			if points < 0 {
				continue
			}

			// save previous balance
			previousBalance := transactions[i].Points

			// subtract points from the transaction
			transactions[i].Points = transactions[i].Points - payload.Points
			// subtract spent points from payload
			payload.Points -= previousBalance - transactions[i].Points

			// if transaction balance has gone negative, put those points back into the payload,
			// set that transaction's balance to 0, and continue
			if transactions[i].Points < 0 {
				payload.Points += transactions[i].Points * -1
				transactions[i].Points = 0
			}

			// instantiate result object
			result := &Result{}

			// set result data points
			result.Payer = transactions[i].Payer
			result.Points = (previousBalance - transactions[i].Points) * -1

			results = append(results, *result)
		}

		return c.JSON(results)
	})

	// get balances by payer
	app.Get("/api/get-balance", func(c *fiber.Ctx) error {
		// declare map for balances by payer
		var balances = make(map[string]int)

		for key, val := range transactions {
			balances[transactions[key].Payer] += val.Points
		}

		return c.JSON(balances)
	})

	// listen on port 8081
	log.Fatal(app.Listen(":8081"))
}
