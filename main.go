package main

import (
	"database/sql"
	"go-todo/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Defined a slice todotasks to store the tasks
var todotasks []models.ToDo
var db *sql.DB

// function to create a new task
func createTasks(c *gin.Context) {
	// Bind JSON data to 'task' struct
	var task models.ToDo
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Prepare the SQL statement for inserting description and status, returning the newly generated ID
	stmt, err := db.Prepare("INSERT INTO todoapi (description, status) VALUES ($1, $2) RETURNING id")
	if err != nil {
		log.Fatal(err)
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

// function to retrieve all todo tasks from database
func getTasks(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	// Query all tasks from the database and saved in variable rows
	rows, err := db.Query("SELECT id,description,status FROM todoapi")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	// Loop through the result rows and transfer rows to 'todotasks' slice
	for rows.Next() {
		var task models.ToDo
		err := rows.Scan(&task.ID, &task.Description, &task.Status)
		if err != nil {
			log.Fatal(err)
		}
		todotasks = append(todotasks, task)
	}

	// responds with every task in the slice
	c.JSON(http.StatusOK, todotasks)
}

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
func updateTask(c *gin.Context) {
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
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the prepared statement with values for respective placeholders in the prepared SQL and scan the updated task data to task variable for displaying in JSON response
	err = stmt.QueryRow(task.Description, "to-do", id).Scan(&task.ID, &task.Description, &task.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}
	// Respond with the updated task
	c.JSON(http.StatusOK, task)
}

// Function to delete a task checking the taskID in URL
func deleteTask(c *gin.Context) {
	id := getTaskID(c)
	if id == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"}) // Respond with an error if the task is not found
		return
	}

	// Prepare the SQL statement for deleting a task(row)
	stmt, err := db.Prepare("DELETE from todoapi where id=$1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	//Executing the prepared statement for the id variable retrieved from get TaskID function
	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
	// Response with "message" : "Task deleted"
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// Function to mark a task as done for the taskID
func markTaskAsDone(c *gin.Context) {
	id := getTaskID(c)
	if id == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	// Preparing the SQL statement for changing the status of the task to "Done" for id
	stmt, err := db.Prepare("UPDATE todoapi SET status = $1 WHERE id = $2")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Execute the prepared statement and changes the status to "Done" for id in the URL
	_, err = stmt.Exec("Done", id)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Response with "message: Task marked as Done"
	c.JSON(http.StatusOK, gin.H{"message": "Task marked as Done"})
}

func main() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()                       // created gin router
	r.POST("/tasks", createTasks)            // Create a new task using POST
	r.GET("/tasks", getTasks)                // Retrieve all tasks using GET
	r.PUT("/tasks/:id", updateTask)          // Update a task using PUT
	r.DELETE("/tasks/:id", deleteTask)       // Delete a task using DELETE
	r.PUT("/tasks/:id/done", markTaskAsDone) // Mark a task as done using PUT
	r.Run(":3000")                           // Start the server at localhost:3000

}
