package main

import "time"

type Task struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"type:varchar(255);not null"`
	Status    string    `json:"status" gorm:"type:varchar(50);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
