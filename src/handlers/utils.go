package handlers

import (
	"bankTest/dbwrapper"
	"bankTest/models"
	"fmt"
)

func validateAccount(account models.Account) error {
	if account.Currency != "USD" && account.Currency != "MXN" && account.Currency != "COP" {
		return fmt.Errorf("currency not supported")
	}
	_, err := dbwrapper.GetClient(account.ClientID)

	if err != nil {
		return fmt.Errorf("client doesnt exists")
	}

	return nil
}

func validateTransaction(transaction models.Transaction) error {
	clientAcc, err := dbwrapper.GetAccount(transaction.ClientID)
	if err != nil {
		return fmt.Errorf("client doesnt exists")
	}

	transferAcc, _ := dbwrapper.GetAccount(transaction.TransferID)

	switch transaction.Type {
	case "deposit":
		return deposit(clientAcc, transaction)

	case "withdraw":
		return withdraw(clientAcc, transaction)

	case "transfer":
		return transfer(clientAcc, transferAcc, transaction)

	default:
		return fmt.Errorf("transaction type is invalid")
	}
}

func updateBalance(account *models.Account, amount float64) error {
	account.Balance += amount
	return dbwrapper.UpdateAccountBalance(*account)
}

func deposit(account models.Account, transaction models.Transaction) error {
	if account == (models.Account{}) {
		return fmt.Errorf("account doesnt exists")
	} else if transaction.Amount < 0 {
		return fmt.Errorf("cannot deposit negative amount")
	} else if account.ID <= 0 {
		return fmt.Errorf("account doesnt exists")
	} else if transaction.TransferID != 0 {
		return fmt.Errorf("extra arguments in the request")
	}

	return updateBalance(&account, transaction.Amount)
}

func withdraw(account models.Account, transaction models.Transaction) error {
	if account == (models.Account{}) {
		return fmt.Errorf("account doesnt exists")
	} else if account.Balance-transaction.Amount < 0 {
		return fmt.Errorf("insuficient funds")
	} else if transaction.Amount >= 0 {
		return fmt.Errorf("cannot withdraw positive amount")
	} else if account.ID <= 0 {
		return fmt.Errorf("account doesnt exists")
	} else if transaction.TransferID != 0 {
		return fmt.Errorf("extra arguments in the request")
	}

	return updateBalance(&account, transaction.Amount)
}

func transfer(sender, receiver models.Account, transaction models.Transaction) error {
	if sender == (models.Account{}) {
		return fmt.Errorf("sender account doesnt exists")
	} else if receiver == (models.Account{}) {
		return fmt.Errorf("receiver account doesnt exists")
	} else if sender.Currency != receiver.Currency {
		return fmt.Errorf("accounts manage different currencies")
	} else if transaction.Amount > sender.Balance {
		return fmt.Errorf("insuficiente funds")
	} else if transaction.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	} else if sender.ID == receiver.ID {
		return fmt.Errorf("sender cannot be the same as receiver")
	}

	err := updateBalance(&sender, (transaction.Amount)*-1)

	if err != nil {
		return err
	}

	return updateBalance(&receiver, transaction.Amount)
}
