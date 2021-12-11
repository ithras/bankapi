package main

import (
	"bankTest/dbwrapper"
	"bankTest/handlers"
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "banktest", "bank")
	dbwrapper.DB, _ = sql.Open("postgres", postgresInfo)
	defer dbwrapper.DB.Close()

	e := echo.New()

	e.POST("/client", handlers.HandlerCreateUser)
	e.GET("/client/:id", handlers.HandlerGetUser)

	e.POST("/account", handlers.HandlerCreateAccount)
	e.GET("/account/:account_id", handlers.HandlerGetAccount)

	e.POST("/transaction", handlers.HandlerCreateTransaction)
	e.GET("/transactions/:account_id", handlers.HandlerGetTransactions)

	e.Logger.Fatal(e.Start(":8080"))
}

/*
func getUsers(c echo.Context) error {
	return c.String()
}
*/
