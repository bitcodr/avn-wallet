package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID      `json:"id" msgpack:"id" validate:"required"`
	Cellphone    uint64         `json:"cellphone" msgpack:"cellphone" validate:"required"`
	FirstName    string         `json:"firstName" msgpack:"firstName" validate:"required,gte=3,lte=70"`
	LastName     string         `json:"lastName" msgpack:"lastName" validate:"required,gte=3,lte=70"`
	Email        string         `json:"email" msgpack:"email" validate:"required,email"`
	Status       string         `json:"status" msgpack:"status" validate:"required,gte=3,lte=25"`
	CreatedAt    string         `json:"createdAt" msgpack:"createdAt"`
	UpdatedAt    string         `json:"updatedAt" msgpack:"updatedAt"`
	Wallet       *Wallet        `json:"wallet" msgpack:"wallet"`
	Transactions []*Transaction `json:"transactions" msgpack:"transactions"`
}
