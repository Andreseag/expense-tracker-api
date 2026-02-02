package models

import "gorm.io/gorm"

// Definimos las categorías como constantes para usarlas en todo el sistema
const (
	CatComida     = "comida"
	CatTransporte = "transporte"
	CatOcio       = "ocio"
	CatServicios  = "servicios"
	CatGeneral    = "general"
)

type Expense struct {
	gorm.Model
	Description string  `json:"description" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Category    string  `json:"category" binding:"required" gorm:"type:varchar(50);default:'general'"`
}