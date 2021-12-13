package main

import (
	"bankTest/dbwrapper"
	"bankTest/handlers"
	"database/sql"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"os"
	"time"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	err := connectDB()
	if err != nil {
		e.Logger.Error(err)
	}

	for dbwrapper.DB.Ping() != nil {
		err := connectDB()
		if err != nil {
			e.Logger.Error(err)
		}

		time.Sleep(5 * time.Second)
	}

	defer dbwrapper.DB.Close()

	e.POST("/client", handlers.HandlerCreateUser)
	e.GET("/client/:id", handlers.HandlerGetUser)

	e.POST("/account", handlers.HandlerCreateAccount)
	e.GET("/account/:account_id", handlers.HandlerGetAccount)

	e.POST("/transaction", handlers.HandlerCreateTransaction)
	e.GET("/transactions/:account_id", handlers.HandlerGetTransactions)

	e.Logger.Fatal(e.Start("0.0.0.0:8080"))
}

func connectDB() error {
	var err error
	dbwrapper.DB, err = sql.Open("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		return err
	}
	return nil
}
