package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Function to delete a task checking the taskID in URL
func DeleteTask(c *gin.Context, db *sql.DB) {
	id := getTaskID(c)
	if id == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"}) // Respond with an error if the task is not found
		return
	}

	// Prepare the SQL statement for deleting a task(row)
	stmt, err := db.Prepare("DELETE from todoapi where id=$1")
	if err != nil {
		log.Fatal(err) // Error handling
	}
	defer stmt.Close()
	//Executing the prepared statement for the id variable retrieved from get TaskID function
	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err) //Error Handling
	}
	// Response with "message" : "Task deleted"
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
