package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}

type UpdateTaskRequest struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}

// Helper validasi field tidak kosong
func validateNonEmpty(title, status string) (bool, string) {
	title = strings.TrimSpace(title)
	status = strings.TrimSpace(status)

	if title == "" {
		return false, "title tidak boleh kosong"
	}
	if status == "" {
		return false, "status tidak boleh kosong"
	}
	return true, ""
}

// POST /tasks
func CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "body JSON tidak valid"})
		return
	}

	ok, msg := validateNonEmpty(req.Title, req.Status)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	task := Task{Title: req.Title, Status: req.Status}
	if err := DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GET /tasks?status=todo&page=1&limit=10
func ListTasks(c *gin.Context) {
	status := strings.TrimSpace(c.Query("status"))

	// default pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // biar gak kebablasan
	}
	offset := (page - 1) * limit

	var tasks []Task
	query := DB.Model(&Task{})

	// filter by status (opsional)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// total untuk metadata pagination
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghitung data"})
		return
	}

	// ambil data
	if err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tasks,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// GET /tasks/:id
func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	var task Task

	if err := DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// PUT /tasks/:id
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task Task

	if err := DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task tidak ditemukan"})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "body JSON tidak valid"})
		return
	}

	ok, msg := validateNonEmpty(req.Title, req.Status)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	task.Title = req.Title
	task.Status = req.Status

	if err := DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DELETE /tasks/:id
func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	// cek ada/tidak
	var task Task
	if err := DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task tidak ditemukan"})
		return
	}

	if err := DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal hapus task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task berhasil dihapus"})
}
