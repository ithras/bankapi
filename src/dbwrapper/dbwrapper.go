package dbwrapper

import (
	"bankTest/models"
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func CreateClient(cli models.Client) error {
	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO clients (name,created_at) VALUES ($1,$2)")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, cli.Name, time.Now())
	if err != nil {
		tx.Rollback()
	}

	err = tx.Commit()

	return err
}

func GetClient(id int) (models.Client, error) {
	cli := models.Client{ID: id}
	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return cli, err
	}

	row := tx.QueryRow("SELECT name FROM clients WHERE id = $1", cli.ID)

	err = row.Scan(&cli.Name)
	if err != nil {
		tx.Rollback()
	}

	err = tx.Commit()

	return cli, err
}

func CreateAccount(account models.Account) error {
	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO accounts(client_id,currency,balance,created_at,updated_at) VALUES($1,$2,$3,$4,$5)")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, account.ClientID, account.Currency, 0, time.Now(), time.Now())
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

func GetAccount(id int) (models.Account, error) {
	account := models.Account{ID: id}

	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)

	if err != nil {
		return models.Account{}, err
	}

	row := tx.QueryRow("SELECT id,client_id,currency,balance,created_at,updated_at FROM accounts WHERE id = $1", &account.ID)

	err = row.Scan(&account.ID, &account.ClientID, &account.Currency, &account.Balance, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		return models.Account{}, err
	}

	err = tx.Commit()

	return account, err
}

func UpdateAccountBalance(account models.Account) error {
	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	_ = tx.QueryRowContext(ctx, "SELECT * FROM accounts WHERE id = $1 FOR UPDATE", account.ID).Scan(&account.ID)

	stmt, err := tx.PrepareContext(ctx, "UPDATE accounts SET balance = $2, updated_at = $3 WHERE id = $1")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, account.ID, account.Balance, time.Now())
	if err != nil {
		tx.Rollback()
	}

	err = tx.Commit()

	return err
}

func CreateTransaction(transaction models.Transaction) error {
	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO transactions (receiver_id,sender_id,amount,type,created_at) VALUES($1,$2,$3,$4,$5)")

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, transaction.ReceiverID, transaction.SenderID, transaction.Amount, transaction.Type, time.Now())
	if err != nil {
		tx.Rollback()
	}

	err = tx.Commit()

	return err
}

func GetTransactions(accID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	var transaction models.Transaction

	ctx := context.Background()
	tx, err := DB.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	rows, err := tx.Query("SELECT id, sender_id, receiver_id, amount,type, created_at FROM transactions WHERE sender_id = $1 OR receiver_id = $1", accID)

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

	err = tx.Commit()

	return transactions, err
}
