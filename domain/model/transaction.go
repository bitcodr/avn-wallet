package model

import "github.com/google/uuid"

type Transaction struct {
	ID          uuid.UUID `json:"id" msgpack:"id" validate:"required"`
	Type        string    `json:"type" msgpack:"type" validate:"required"`
	Status      string    `json:"status" msgpack:"status" validate:"required"`
	Description string    `json:"description" msgpack:"status" validate:"required"`
	CreatedAt   string    `json:"createdAt" msgpack:"createdAt"`
	Cause       string    `json:"cause" msgpack:"cause"`
	Balance     float64   `json:"balance" msgpack:"balance" validate:"required"`
	User        *User     `json:"user" msgpack:"user"`
}
