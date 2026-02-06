package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	InitDB()

	// auto create/update table sesuai struct Task
	if err := DB.AutoMigrate(&Task{}); err != nil {
		log.Fatal("Gagal migrate: ", err)
	}

	r := gin.Default()

	r.POST("/tasks", CreateTask)
	r.GET("/tasks", ListTasks)
	r.GET("/tasks/:id", GetTaskByID)
	r.PUT("/tasks/:id", UpdateTask)
	r.DELETE("/tasks/:id", DeleteTask)

	log.Println("Server running di http://localhost:8080")
	r.Run(":8080")
}
