package main

import (
	"go-todo/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Defined a slice todotasks to store the tasks
var todotasks []model.ToDo

// Variable to keep track of the task IDs. like a counter which increments when a new task is created.
var taskIDcount int

// function to create a new task
func createTasks(c *gin.Context) {
	var task model.ToDo
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Assign a unique ID from ID counter to the task and increments the ID count and also assigns "to-do" status to the task

	task.ID = taskIDcount
	taskIDcount++
	task.Status = "to-do"

	// appends the new task to the todotasks slice
	todotasks = append(todotasks, task)

	// Respond with the newly created task
	c.JSON(http.StatusCreated, task)
}

// function to retrieve all todo tasks from the slice
func getTasks(c *gin.Context) {
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

	var task model.ToDo
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
	r := gin.Default()                       // created gin router
	r.POST("/tasks", createTasks)            // Create a new task using POST
	r.GET("/tasks", getTasks)                // Retrieve all tasks using GET
	r.PUT("/tasks/:id", updateTask)          // Update a task using PUT
	r.DELETE("/tasks/:id", deleteTask)       // Delete a task using DELETE
	r.PUT("/tasks/:id/done", markTaskAsDone) // Mark a task as done using PUT
	r.Run()                                  // Start the server

}
