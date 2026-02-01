package controllers

import (
	"net/http"

	"github.com/Andreseag/expense-tracker-api/config"
	"github.com/Andreseag/expense-tracker-api/models"
	"github.com/gin-gonic/gin"
)

func GetExpenses(c *gin.Context) {
	var expenses []models.Expense
	config.DB.Order("created_at desc").Find(&expenses)
	
	// Gin se encarga de los headers y de convertir a JSON
	c.JSON(http.StatusOK, expenses)
}

func CreateExpense(c *gin.Context) {
	var t models.Expense
	
	// ShouldBindJSON es el equivalente a Decode
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// VALIDACIÓN
	if t.Amount <= 0 || t.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Descripción obligatoria y monto debe ser mayor a 0"})
		return
	}
	
	
	config.DB.Create(&t)
	c.JSON(http.StatusCreated, t)
}

func UpdateExpense(c *gin.Context) {
	id := c.Param("id") 
	
	var expense models.Expense
	if err := config.DB.First(&expense, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tarea no encontrada"})
		return
	}	

	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	config.DB.Save(&expense)
	c.JSON(http.StatusOK, expense)
}

func DeleteExpense(c *gin.Context) {
	id := c.Param("id")

	var expense models.Expense
	if err := config.DB.First(&expense, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expense no encontrado"})
	}


	config.DB.Delete(&expense)
  c.JSON(http.StatusOK, expense)
} 

func GetSummary(c *gin.Context) {
	var summary float64
	month := c.Query("month")
	months := map[string]string{
		"1": "enero", "2": "febrero", "3": "marzo", "4": "abril",
		"5": "mayo", "6": "junio", "7": "julio", "8": "agosto",
		"9": "septiembre", "10": "octubre", "11": "noviembre", "12": "diciembre",
	}
	// GORM ejecutará: SELECT SUM(amount) FROM expenses WHERE deleted_at IS NULL
	query := config.DB.Model(&models.Expense{}).Select("COALESCE(SUM(amount), 0)")

	if month != "" {
		// PostgreSQL: EXTRACT(MONTH FROM created_at)
		query = query.Where("EXTRACT(MONTH FROM created_at) = ?", month)
	}

	err := query.Row().Scan(&summary)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al calcular el resumen"})
		return
	}

	monthName := "Todos los meses"
	if val, ok := months[month]; ok {
		monthName = val
	}

	c.JSON(http.StatusOK, gin.H{
		"month": monthName,
		"total": summary,
	})

}