package config

import (
	"log"

	"github.com/Andreseag/expense-tracker-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB es una variable global que usaremos en los controladores
var DB *gorm.DB

func ConectarDB() {
	// Estos datos deben coincidir con tu docker-compose.yml
	dsn := "host=localhost user=dev_user password=dev_password dbname=expense-tracker port=5433 sslmode=disable"
	
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("❌ Error al conectar a la DB:", err)
	}

	// 👇 ESTA LÍNEA ES LA MAGIA
	// GORM lee el struct 'Expense' y crea la tabla 'exprenses' en la DB
	err = DB.AutoMigrate(&models.Expense{})
	if err != nil {
		log.Println("❌ Error en la migración:", err)
	}

	log.Println("✅ Base de datos conectada y tablas migradas")
}