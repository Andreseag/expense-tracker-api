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

// func CreateTask(c *gin.Context) {
// 	var t models.Task
	
// 	// ShouldBindJSON es el equivalente a Decode
// 	if err := c.ShouldBindJSON(&t); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
// 		return
// 	}

// 	// VALIDACIÓN
// 	if t.Status != models.StatusTodo && 
// 	   t.Status != models.StatusInProgress && 
// 	   t.Status != models.StatusDone {
// 		if t.Status == "" {
// 			t.Status = models.StatusTodo
// 		} else {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Estado no válido"})
// 			return
// 		}
// 	}
	
// 	config.DB.Create(&t)
// 	c.JSON(http.StatusCreated, t)
// }

// func UpdateTask(c *gin.Context) {
// 	id := c.Param("id") 
	
// 	var task models.Task
// 	if err := config.DB.First(&task, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Tarea no encontrada"})
// 		return
// 	}

// 	if err := c.ShouldBindJSON(&task); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
// 		return
// 	}

// 	config.DB.Save(&task)
// 	c.JSON(http.StatusOK, task)
// }

// func DeleteTask(c *gin.Context) {
// 	id := c.Param("id")

// 	var task models.Task
// 	if err := config.DB.First(&task, id).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Tarea no encontrada"})
// 	}

// 	if err := c.ShouldBindJSON(&task); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
// 	}

// 	config.DB.Delete(&task)

// }