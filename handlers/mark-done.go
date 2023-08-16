package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Function to mark a task as done for the taskID
func MarkTaskAsDone(c *gin.Context, db *sql.DB) {
	id := getTaskID(c)
	if id == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	// Preparing the SQL statement for changing the status of the task to "Done" for id
	stmt, err := db.Prepare("UPDATE todoapi SET status = $1 WHERE id = $2")
	if err != nil {
		log.Fatal(err) // Error handling
	}
	defer stmt.Close()

	// Execute the prepared statement and changes the status to "Done" for id in the URL
	_, err = stmt.Exec("Done", id)
	if err != nil {
		log.Fatal(err) // Error handling
		return
	}
	// Response with "message: Task marked as Done"
	c.JSON(http.StatusOK, gin.H{"message": "Task marked as Done"})
}
