package dbwrapper

import (
	"bankTest/models"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func CreateClient(cli models.Client) error {
	stmt, err := DB.Prepare("INSERT INTO clients (name,created_at) VALUES ($1,$2)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(cli.Name, time.Now())

	return err
}

func GetClient(id int) (models.Client, error) {
	cli := models.Client{ID: id}
	err := DB.QueryRow("SELECT name FROM clients WHERE id = $1", cli.ID).Scan(&cli.Name)

	return cli, err
}

func CreateAccount(account models.Account) error {
	stmt, err := DB.Prepare("INSERT INTO accounts(client_id,currency,balance,created_at,updated_at) VALUES($1,$2,$3,$4,$5)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(account.ClientID, account.Currency, 0, time.Now(), time.Now())

	return err
}

func GetAccount(id int) (models.Account, error) {
	account := models.Account{ID: id}
	//err := DB.QueryRow("SELECT * FROM accounts WHERE id = $1", account.ID).Scan(&account)
	err := DB.QueryRow("SELECT id,client_id,currency,balance,created_at,updated_at FROM accounts WHERE id = $1", account.ID).
		Scan(&account.ID, &account.ClientID, &account.Currency, &account.Balance, &account.CreatedAt, &account.UpdatedAt)

	return account, err
}

func UpdateAccountBalance(account models.Account) error {
	stmt, err := DB.Prepare("UPDATE accounts SET balance = $2, updated_at = $3 WHERE id = $1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(account.ID, account.Balance, time.Now())

	return err
}

func CreateTransaction(transaction models.Transaction) error {
	stmt, err := DB.Prepare("INSERT INTO transactions (receiver_id,sender_id,amount,type,created_at) VALUES($1,$2,$3,$4,$5)")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(transaction.ReceiverID, transaction.SenderID, transaction.Amount, transaction.Type, time.Now())

	return err
}

func GetTransactions(accID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	var transaction models.Transaction

	rows, err := DB.Query("SELECT id, sender_id, receiver_id, amount,type, created_at FROM transactions WHERE sender_id = $1 OR receiver_id = $1", accID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&transaction.ID, &transaction.SenderID, &transaction.ReceiverID, &transaction.Amount, &transaction.Type, &transaction.CreatedAt)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return transactions, err
}
