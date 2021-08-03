package models

import "github.com/google/uuid"

type Cat struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Breed string    `json:"breed"`
	Color string    `json:"color"`
	Age   float32   `json:"age"`
	Price float32   `json:"price"`
}
