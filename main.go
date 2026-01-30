package main

import (
	"time"

	"github.com/Andreseag/expense-tracker-api/config"
	"github.com/Andreseag/expense-tracker-api/controllers"
	"github.com/gin-contrib/cors"

	// "github.com/Andreseag/expense-tracker-api/controllers"
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
 
	// // Envolvemos las rutas con el middleware
  // r.GET("/api/tasks", controllers.GetTasks)
	// r.POST("/api/tasks/new", controllers.CreateTask)
	r.GET("/api/expenses", controllers.GetExpenses)
	
	// // ¡Aquí está lo que pediste! El :id es el parámetro
	// r.PUT("/api/tasks/:id", controllers.UpdateTask)
	// r.DELETE("/api/tasks/:id", controllers.DeleteTask) 

	r.Run(":8080")
}