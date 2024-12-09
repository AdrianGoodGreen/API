package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Define la estructura de la actividad
type Activity struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	TeacherID        int    `json:"teacher_id"`
	EnrolledStudents []int  `json:"enrolled_students"`
}

func main() {
	// Conectar a la base de datos MySQL
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(mysql-service:3306)/activitydb")
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()
	// Configurar el pool de conexiones
	db.SetMaxOpenConns(10) // Ajusta según lo necesario
	db.SetMaxIdleConns(5)  // Ajusta según lo necesario
	db.SetConnMaxLifetime(time.Minute * 5)

	// Configurar el router de Gin
	router := gin.Default()

	// Definir las rutas
	router.GET("/activities", getActivities)
	router.GET("/activities/:id", getActivityByID)
	router.POST("/activities", createActivity)
	router.PUT("/activities/:id", updateActivity)
	router.DELETE("/activities/:id", deleteActivity)

	// Iniciar el servidor
	router.Run(":8083")
}

// Obtener todas las actividades
func getActivities(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, description, teacher_id FROM activities")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var activities []Activity
	for rows.Next() {
		var activity Activity
		if err := rows.Scan(&activity.ID, &activity.Name, &activity.Description, &activity.TeacherID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, activities)
}

// Obtener una actividad específica
func getActivityByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var activity Activity
	err = db.QueryRow("SELECT id, name, description, teacher_id FROM activities WHERE id = ?", id).Scan(&activity.ID, &activity.Name, &activity.Description, &activity.TeacherID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Actividad no encontrada"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

// Crear una nueva actividad
func createActivity(c *gin.Context) {
	var newActivity Activity
	if err := c.BindJSON(&newActivity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insertar la nueva actividad en la base de datos
	stmt, err := db.Prepare("INSERT INTO activities(name, description, teacher_id) VALUES(?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(newActivity.Name, newActivity.Description, newActivity.TeacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Obtener el ID de la actividad insertada
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newActivity.ID = int(id)
	c.JSON(http.StatusCreated, newActivity)
}

// Actualizar una actividad existente
func updateActivity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var updatedActivity Activity
	if err := c.BindJSON(&updatedActivity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Actualizar la actividad en la base de datos
	_, err = db.Exec("UPDATE activities SET name = ?, description = ?, teacher_id = ? WHERE id = ?", updatedActivity.Name, updatedActivity.Description, updatedActivity.TeacherID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedActivity.ID = id
	c.JSON(http.StatusOK, updatedActivity)
}

// Eliminar una actividad
func deleteActivity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Eliminar la actividad de la base de datos
	_, err = db.Exec("DELETE FROM activities WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Actividad eliminada"})
}
