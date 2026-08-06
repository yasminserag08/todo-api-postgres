package main

import "time"

type Todo struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"not null" json:"title"`
	Completed   bool       `gorm:"default:false" json:"completed"`
	Category    string     `gorm:"default:'General'" json:"category"`
	Piority     string     `gorm:"default:'Medium'" json:"priority"`
	CompletedAt *time.Time `json:"completedAt"`
	DueDate     *time.Time `json:"dueDate"`
}

// used pointers to *time.Time so value can be nil and to make comparisons between dates easier than string.
