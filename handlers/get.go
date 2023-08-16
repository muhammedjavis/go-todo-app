package handlers

import (
	"database/sql"
	"go-todo/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// function to retrieve all todo tasks from database
func GetTasks(c *gin.Context, db *sql.DB) {
	c.Header("Content-Type", "application/json")
	// Query all tasks from the database and saved in variable rows
	rows, err := db.Query("SELECT id,description,status FROM todoapi")
	if err != nil {
		log.Fatal(err) // Error handling
	}
	defer rows.Close()

	// Defined a slice todotasks to store the tasks
	var todotasks []models.ToDo

	// Loop through the result rows and transfer rows to 'todotasks' slice
	for rows.Next() {
		var task models.ToDo
		err := rows.Scan(&task.ID, &task.Description, &task.Status)
		if err != nil {
			log.Fatal(err) // Error handling
		}

		todotasks = append(todotasks, task)
	}

	// responds with every task in the slice
	c.JSON(http.StatusOK, todotasks)
}
