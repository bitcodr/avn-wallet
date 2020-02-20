package model

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID `json:"id" msgpack:"id" validate:"required"`
	Charge    float64    `json:"charge" msgpack:"charge" validate:"required"`
	CreatedAt time.Time `json:"createdAt" msgpack:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" msgpack:"updatedAt"`
	User      *User     `json:"user" msgpack:"user" validate:"required"`
}
