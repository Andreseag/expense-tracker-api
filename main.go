package main

import (
	"time"

	"github.com/Andreseag/expense-tracker-api/config"
	"github.com/Andreseag/expense-tracker-api/controllers"
	"github.com/gin-contrib/cors"

	"github.com/Andreseag/expense-tracker-api/models"
	"github.com/gin-gonic/gin"
)


func main() {
	config.ConectarDB()
	config.DB.AutoMigrate(&models.Expense{})

	r := gin.Default()

	// Configuración de CORS profesional
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		MaxAge:           12 * time.Hour,
	}))
 
	// Requests
	r.GET("/api/expenses", controllers.GetExpenses)
	r.POST("/api/expenses/new", controllers.CreateExpense)
	r.PUT("/api/expenses/:id", controllers.UpdateExpense)
	r.GET("/api/expenses/summary", controllers.GetSummary) 
	r.DELETE("/api/expenses/:id", controllers.DeleteExpense)
	r.GET("/api/expenses/export", controllers.ExportExpensesCSV)


	r.Run(":8080")
}