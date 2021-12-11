package handlers

import (
	"bankTest/dbwrapper"
	"bankTest/models"
	"fmt"
)

func accountValidator(account models.Account) error {
	if account.Currency != "USD" && account.Currency != "MXN" && account.Currency != "COP" {
		return fmt.Errorf("currency not supported")
	}
	_, err := dbwrapper.GetClient(account.ClientID)

	if err != nil {
		return fmt.Errorf("client doesnt exists")
	}

	return nil
}

func transactionValidator(transaction models.Transaction) error {
	senderAcc, _ := dbwrapper.GetAccount(transaction.SenderID)

	receiverAcc, _ := dbwrapper.GetAccount(transaction.ReceiverID)

	switch transaction.Type {
	case "deposit":
		return transactionDeposit(receiverAcc, transaction)

	case "withdraw":
		return transactionWithdraw(senderAcc, transaction)

	case "transfer":
		return transactionTransfer(senderAcc, receiverAcc, transaction)
	}

	return nil
}

func updateBalance(account *models.Account, amount float64) error {
	account.Balance += amount
	return dbwrapper.UpdateAccountBalance(*account)
}

func transactionDeposit(account models.Account, transaction models.Transaction) error {
	if transaction.Amount < 0 {
		return fmt.Errorf("cannot deposit negative amount")
	} else if account.ID <= 0 {
		return fmt.Errorf("account doesnt exists")
	}

	return updateBalance(&account, transaction.Amount)
}

func transactionWithdraw(account models.Account, transaction models.Transaction) error {
	if transaction.Amount > account.Balance {
		return fmt.Errorf("insuficient funds")
	} else if transaction.Amount >= 0 {
		return fmt.Errorf("cannot withdraw positive amount")
	} else if account.ID <= 0 {
		return fmt.Errorf("account doesnt exists")
	}

	return updateBalance(&account, transaction.Amount)
}

func transactionTransfer(sender, receiver models.Account, transaction models.Transaction) error {
	if sender.Currency != receiver.Currency {
		return fmt.Errorf("insuficiente funds")
	} else if transaction.Amount > sender.Balance {
		return fmt.Errorf("accounts manage different currencies")
	} else if transaction.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	err := updateBalance(&sender, (transaction.Amount)*-1)

	if err != nil {
		return err
	}

	return updateBalance(&receiver, transaction.Amount)
}
