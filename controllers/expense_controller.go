package controllers

import (
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/Andreseag/expense-tracker-api/config"
	"github.com/Andreseag/expense-tracker-api/models"
	"github.com/gin-gonic/gin"
)

// Función auxiliar para validar
func isValidCategory(cat string) bool {
	switch cat {
	case models.CatComida, models.CatTransporte, models.CatOcio, models.CatServicios, models.CatGeneral:
		return true
	}
	return false
}

func GetExpenses(c *gin.Context) {
	var expenses []models.Expense
	config.DB.Order("created_at desc").Find(&expenses)
	
	// Filters
	category := c.Query("category")

	query := config.DB.Order("created_at desc")

	if category != "" {
		// 1. Validar si la categoría existe en nuestra lista permitida
		if !isValidCategory(category) {
			// Opción A: Devolver error 400 (Bad Request)
			c.JSON(http.StatusBadRequest, gin.H{"error": "La categoría solicitada no existe"})
			return
		}
		// 2. Si es válida, filtramos
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los gastos"})
		return
	}
	
	// Gin se encarga de los headers y de convertir a JSON
	c.JSON(http.StatusOK, expenses)
}

func CreateExpense(c *gin.Context) {
	var exp models.Expense

	if err := c.ShouldBindJSON(&exp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos o incompletos"})
		return
	}

	// VALIDACIÓN DE CATEGORÍA
	switch exp.Category {
	case models.CatComida, models.CatTransporte, models.CatOcio, models.CatServicios, models.CatGeneral:
		// Es válida, no hacemos nada
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Categoría no permitida"})
		return
	}

	if err := config.DB.Create(&exp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo guardar el gasto"})
		return
	}

	c.JSON(http.StatusCreated, exp)
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

func ExportExpensesCSV(c *gin.Context) {
	var expenses []models.Expense
	// Traemos todos los gastos
	config.DB.Find(&expenses)

	// 1. Configuramos los Headers para la descarga
	c.Header("Content-Disposition", "attachment; filename=gastos.csv")
	c.Header("Content-Type", "text/csv")
	c.Header("Transfer-Encoding", "chunked")

	// 2. Creamos el escritor de CSV que apunta directamente al body de la respuesta
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// 3. Escribimos la cabecera (los títulos de las columnas)
	header := []string{"ID", "Descripcion", "Monto", "Fecha"}
	if err := writer.Write(header); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el archivo"})
		return
	}

	// 4. Recorremos los gastos y los escribimos en el CSV
	for _, e := range expenses {
		row := []string{
			fmt.Sprintf("%d", e.ID),
			e.Description,
			fmt.Sprintf("%.2f", e.Amount),
			e.CreatedAt.Format("2006-01-02 15:04:05"), // Formato de fecha legible
		}
		writer.Write(row)
	}
}

