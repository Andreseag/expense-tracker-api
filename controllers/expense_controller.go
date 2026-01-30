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

func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	var expense models.Expense
	if err := config.DB.First(&expense, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expense no encontrado"})
	}

	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
	}

	config.DB.Delete(&expense)

} 

func GetSummary(c *gin.Context) {
	var summary float64

	// GORM ejecutará: SELECT SUM(amount) FROM expenses WHERE deleted_at IS NULL
	err := config.DB.Model(&models.Expense{}).Select("SUM(amount)").Row().Scan(&summary)

	if err != nil {
			// Si la tabla está vacía, Scan podría dar error o devolver 0
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo calcular el total"})
			return
	}

	c.JSON(http.StatusOK, gin.H{"total": summary})

}