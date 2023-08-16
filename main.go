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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	// Set the task ID and update the task details in the slice
	task.ID = id
	todotasks[id] = task
	task.Status = "to-do"

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
	// Remove the task from the todotasks slice by appending tasks that are before and after the index id and excludes the task with id index
	todotasks = append(todotasks[:id], todotasks[id+1:]...)
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
	// Change the task status of the task with index 'id' to "Done"
	todotasks[id].Status = "Done"
	// Respond with the updated task with the new task status as "done"
	c.JSON(http.StatusOK, todotasks[id])
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
	r.Run()                                  // Start the server

}
