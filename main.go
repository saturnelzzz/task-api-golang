package main

import (
	"log"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "task-api/docs" 
)

// @title Task Management API
// @version 1.0
// @description REST API untuk manajemen task (CRUD + pagination + filter status)
// @host localhost:8080
// @BasePath /
func main() {
	InitDB()
	if err := DB.AutoMigrate(&Task{}); err != nil {
		log.Fatal("Gagal migrate: ", err)
	}

	r := gin.Default()

	r.POST("/tasks", CreateTask)
	r.GET("/tasks", ListTasks)
	r.GET("/tasks/:id", GetTaskByID)
	r.PUT("/tasks/:id", UpdateTask)
	r.DELETE("/tasks/:id", DeleteTask)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server running di http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Gagal menjalankan server: ", err)
	}
}
