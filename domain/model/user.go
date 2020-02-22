package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID      `json:"id" msgpack:"id"`
	Cellphone    uint64         `json:"cellphone" msgpack:"cellphone" validate:"required"`
	FirstName    string         `json:"firstName" msgpack:"firstName"`
	LastName     string         `json:"lastName" msgpack:"lastName"`
	Email        string         `json:"email" msgpack:"email"`
	Status       string         `json:"status" msgpack:"status"`
	CreatedAt    string         `json:"createdAt" msgpack:"createdAt"`
	UpdatedAt    string         `json:"updatedAt" msgpack:"updatedAt"`
	Wallet       *Wallet        `json:"wallet" msgpack:"wallet"`
	Transactions []*Transaction `json:"transactions" msgpack:"transactions"`
}
