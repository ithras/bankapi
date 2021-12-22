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
	senderAcc, _ := dbwrapper.GetAccount(transaction.SenderID)

	receiverAcc, _ := dbwrapper.GetAccount(transaction.ReceiverID)

	switch transaction.Type {
	case "deposit":
		return deposit(receiverAcc, transaction)

	case "withdraw":
		return withdraw(senderAcc, transaction)

	case "transfer":
		return transfer(senderAcc, receiverAcc, transaction)

	default:
		return fmt.Errorf("transaction type is invalid")
	}
}

func updateBalance(account *models.Account, amount float64) error {
	account.Balance += amount
	return dbwrapper.UpdateAccountBalance(*account)
}

func deposit(account models.Account, transaction models.Transaction) error {
	if transaction.Amount < 0 {
		return fmt.Errorf("cannot deposit negative amount")
	} else if account.ID <= 0 {
		return fmt.Errorf("account doesnt exists")
	}

	return updateBalance(&account, transaction.Amount)
}

func withdraw(account models.Account, transaction models.Transaction) error {
	if transaction.Amount > account.Balance {
		return fmt.Errorf("insuficient funds")
	} else if transaction.Amount >= 0 {
		return fmt.Errorf("cannot withdraw positive amount")
	} else if account.ID <= 0 {
		return fmt.Errorf("account doesnt exists")
	}

	return updateBalance(&account, transaction.Amount)
}

func transfer(sender, receiver models.Account, transaction models.Transaction) error {
	if sender.Currency != receiver.Currency {
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
