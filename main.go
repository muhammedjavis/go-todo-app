package main

import (
	"database/sql"
	"go-todo/handlers"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	// Establishing connection to Postgres Database
	var db *sql.DB
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default() // created gin router

	// Create a new task using POST request

	r.POST("/tasks", func(c *gin.Context) {
		handlers.CreateTasks(c, db)
	})

	// Retrieve all tasks using GET request

	r.GET("/tasks", func(c *gin.Context) {
		handlers.GetTasks(c, db)
	})

	// Update a task using PUT request

	r.PUT("/tasks/:id", func(c *gin.Context) {
		handlers.UpdateTask(c, db)
	})

	// Delete a task using DELETE request

	r.DELETE("/tasks/:id", func(c *gin.Context) {
		handlers.DeleteTask(c, db)
	})

	// Mark a task as done using PUT request

	r.PUT("/tasks/:id/done", func(c *gin.Context) {
		handlers.MarkTaskAsDone(c, db)
	})

	r.Run(":3000") // Start the server and listen at localhost:3000

}
