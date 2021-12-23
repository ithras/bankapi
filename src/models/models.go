package models

import (
	"time"
)

type Client struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Account struct {
	ClientID  int       `json:"client_id"`
	ID        int       `json:"account_id"`
	Currency  string    `json:"currency"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Transaction struct {
	ClientID   int       `json:"client_id"`
	TransferID int       `json:"transfer_id"`
	ID         int       `json:"transaction_id"`
	Type       string    `json:"type"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}
