package handlers

import (
	"database/sql"
	"go-todo/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// function to create a new task
func CreateTasks(c *gin.Context, db *sql.DB) {
	// Bind JSON data to 'task' struct
	var task models.ToDo
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Prepare the SQL statement for inserting description and status, returning the newly generated ID
	stmt, err := db.Prepare("INSERT INTO todoapi (description, status) VALUES ($1, $2) RETURNING id")
	if err != nil {
		log.Fatal(err) //Error handling
	}
	defer stmt.Close()

	var newID int // To store the new ID returned by RETURNING from db.Prepare

	// Execute the prepared statement and scan the new ID
	err = stmt.QueryRow(task.Description, "to-do").Scan(&newID)
	if err != nil {
		log.Fatal(err)
	}
	task.ID = newID
	task.Status = "to-do"
	// Respond with the newly created task
	c.JSON(http.StatusCreated, task)
}
