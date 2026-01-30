package models

import "gorm.io/gorm"

// Definimos las constantes (el "diccionario" de estados)


type Expense struct {
	gorm.Model
	Description string `json:"description" binding:"required"`
	Amount      float32 `json:"amount" binding:"required"`
}