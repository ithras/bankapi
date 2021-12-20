package handlers

import (
	"bankTest/dbwrapper"
	"bankTest/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func HandlerCreateUser(c echo.Context) error {
	cli := models.Client{}
	dec := json.NewDecoder(c.Request().Body)
	err := dec.Decode(&cli)

	if err != nil {
		return c.String(http.StatusInternalServerError, "json format is invalid")
	}

	err = dbwrapper.CreateClient(cli)

	if err != nil {
		return c.String(http.StatusInternalServerError, "invalid database connection")
	}

	return c.String(http.StatusOK, "user created")
}

func HandlerGetUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	cli, err := dbwrapper.GetClient(id)

	if err != nil {
		return c.String(http.StatusBadRequest, "client not found")
	}

	return c.JSON(http.StatusOK, cli)
}

func HandlerCreateAccount(c echo.Context) error {
	account := models.Account{}
	dec := json.NewDecoder(c.Request().Body)
	err := dec.Decode(&account)

	if err != nil {
		return c.String(http.StatusBadRequest, "json format is invalid")
	}

	err = accountValidator(account)

	if err != nil {
		return c.String(http.StatusBadRequest, "client doesnt exists")
	}

	err = dbwrapper.CreateAccount(account)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Invalid database connection")
	}

	return c.String(http.StatusOK, "account created")
}

func HandlerGetAccount(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("account_id"))
	account, err := dbwrapper.GetAccount(id)

	if err != nil {
		return c.String(http.StatusNotFound, "account not found")
	}

	return c.JSON(http.StatusOK, account)
}

func HandlerCreateTransaction(c echo.Context) error {

	transaction := models.Transaction{}
	dec := json.NewDecoder(c.Request().Body)
	err := dec.Decode(&transaction)

	if err != nil {
		return c.String(http.StatusInternalServerError, "json format is invalid")
	}

	err = transactionValidator(transaction)

	if err != nil {
		return c.String(http.StatusBadRequest, "transaction is invalid: "+err.Error())
	}

	err = dbwrapper.CreateTransaction(transaction)

	if err != nil {
		return c.String(http.StatusBadRequest, "invalid database connection")
	}

	return c.String(http.StatusOK, "Transaction Complete")
}

func HandlerGetTransactions(c echo.Context) error {
	account_id, _ := strconv.Atoi(c.Param("account_id"))

	transactions, err := dbwrapper.GetTransactions(account_id)

	if err != nil {
		return c.String(http.StatusNotFound, "account not found")
	}

	return c.JSON(http.StatusOK, transactions)
}
