package main

import (
	"go-todo/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var todotasks []model.ToDo
var taskIDcount int

func createTasks(c *gin.Context) {
	var task model.ToDo
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = taskIDcount
	taskIDcount++
	task.Status = "to-do"
	todotasks = append(todotasks, task)

	c.JSON(http.StatusCreated, task)
}

func getTasks(c *gin.Context) {
	c.JSON(http.StatusOK, todotasks)
}
func getTaskID(c *gin.Context) int {
	id := c.Param("id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		return -1
	}
	return taskID
}

func updateTask(c *gin.Context) {
	id := getTaskID(c)
	if id == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var task model.ToDo
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = id
	todotasks[id] = task
	task.Status = "to-do"
	c.JSON(http.StatusOK, task)
}

func deleteTask(c *gin.Context) {
	id := getTaskID(c)
	if id == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	todotasks = append(todotasks[:id], todotasks[id+1:]...)

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

func markTaskAsDone(c *gin.Context) {
	id := getTaskID(c)
	if id == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	todotasks[id].Status = "Done"

	c.JSON(http.StatusOK, todotasks[id])
}

func main() {
	r := gin.Default()
	r.POST("/tasks", createTasks)
	r.GET("/tasks", getTasks)
	r.PUT("/tasks/:id", updateTask)
	r.DELETE("/tasks/:id", deleteTask)
	r.PUT("/tasks/:id/done", markTaskAsDone)
	r.Run()

}
