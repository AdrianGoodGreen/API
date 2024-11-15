package main

// Adrian Manuel Escogido Antonio
// Alberto Brenes Fernandez
// Topicos para el despliegue de aplicaciones

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Activity struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	TeacherID        int    `json:"teacher_id"`
	EnrolledStudents []int  `json:"enrolled_students"` // Lista de IDs de estudiantes inscritos
}

var activities = []Activity{
	{ID: 1, Name: "Fútbol", Description: "Entrenamiento de fútbol para principiantes", TeacherID: 101, EnrolledStudents: []int{1, 2}},
	{ID: 2, Name: "Pintura", Description: "Clases de pintura avanzada", TeacherID: 102, EnrolledStudents: []int{3}},
	{ID: 3, Name: "Robótica", Description: "Introducción a la robótica", TeacherID: 103, EnrolledStudents: []int{}},
}

var nextID = 4 // Variable para gestionar el ID único de las actividades

func main() {
	router := gin.Default()

	router.GET("/activities", getActivities)             // Obtener todas las actividades
	router.GET("/activities/:id", getActivityByID)       // Obtener una actividad específica
	router.POST("/activities", createActivity)           // Crear una nueva actividad
	router.PUT("/activities/:id", updateActivity)        // Actualizar una actividad
	router.DELETE("/activities/:id", deleteActivity)     // Eliminar una actividad
	router.POST("/activities/:id/enroll", enrollStudent) // Inscribir a un estudiante en una actividad

	router.Run(":8083")
}

// Obtener todas las actividades
func getActivities(c *gin.Context) {
	c.JSON(http.StatusOK, activities)
}

// Obtener una actividad específica por ID
func getActivityByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	for _, activity := range activities {
		if activity.ID == id {
			c.JSON(http.StatusOK, activity)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Actividad no encontrada"})
}

// Crear una nueva actividad
func createActivity(c *gin.Context) {
	var newActivity Activity
	if err := c.BindJSON(&newActivity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newActivity.ID = nextID
	nextID++
	activities = append(activities, newActivity)
	c.JSON(http.StatusCreated, newActivity)
}

// Actualizar una actividad existente por ID
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

	for i, activity := range activities {
		if activity.ID == id {
			updatedActivity.ID = id
			activities[i] = updatedActivity
			c.JSON(http.StatusOK, updatedActivity)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Actividad no encontrada"})
}

// Eliminar una actividad existente por ID
func deleteActivity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	for i, activity := range activities {
		if activity.ID == id {
			activities = append(activities[:i], activities[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"status": "Actividad eliminada", "activity": activity.Name})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Actividad no encontrada"})
}

// Inscribir a un estudiante en una actividad
func enrollStudent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var student struct {
		StudentID int `json:"student_id"`
	}
	if err := c.BindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, activity := range activities {
		if activity.ID == id {
			// Verificar si el estudiante ya está inscrito
			for _, studentID := range activities[i].EnrolledStudents {
				if studentID == student.StudentID {
					c.JSON(http.StatusBadRequest, gin.H{"error": "El estudiante ya está inscrito en esta actividad"})
					return
				}
			}
			activities[i].EnrolledStudents = append(activities[i].EnrolledStudents, student.StudentID)
			c.JSON(http.StatusOK, gin.H{"status": "Estudiante inscrito en la actividad"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Actividad no encontrada"})
}
