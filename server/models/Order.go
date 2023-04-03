package models

import "time"

type Order struct {
	ID                    int       `json:"id"`
	OrderID               string    `json:"order_id"`
	Book                  Book      `json:"book"`
	Recipient             *User     `json:"recipient"`
	City                  string    `json:"city"`
	Address               string    `json:"address"`
	Street                string    `json:"street"`
	Comment               string    `json:"comment"`
	ExpectedReceivingTime time.Time `json:"expected_receiving_time"`
}
