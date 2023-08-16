package handlers

import (
	"database/sql"
	"go-todo/models"
	"log"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

// Function to extract the task ID from the URL parameter
func getTaskID(c *gin.Context) int {
	id := c.Param("id")
	taskID, err := strconv.Atoi(id) // converts the string value to int and stores in taskID variable
	if err != nil {
		return -1
	}
	return taskID
}

// function to update the task using taskid and changing the other attributes in it
func UpdateTask(c *gin.Context, db *sql.DB) {
	id := getTaskID(c)
	if id == -1 {
		// Respond with an error if the task is not found
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var task models.ToDo
	//binds the request body to task variable
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Respond with an error if there's a problem with JSON binding
		return
	}

	// Prepare the SQL statement for updating description and status on db and returns id,description,status for scanning to task variable later
	stmt, err := db.Prepare("UPDATE todoapi SET description = $1, status = $2 WHERE id = $3 RETURNING id, description, status")
	if err != nil {
		log.Fatal(err) // Error handling
	}
	defer stmt.Close()

	// Execute the prepared statement with values for respective placeholders in the prepared SQL and scan the updated task data to task variable for displaying in JSON response
	err = stmt.QueryRow(task.Description, "to-do", id).Scan(&task.ID, &task.Description, &task.Status)
	if err != nil {
		log.Fatal(err) // Error handling
		return
	}
	// Respond with the updated task
	c.JSON(http.StatusOK, task)
}
