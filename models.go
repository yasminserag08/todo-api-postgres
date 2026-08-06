package main

import "time"

type Todo struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Completed   bool       `gorm:"default:false;not null" json:"completed"`
	Category    string     `gorm:"type:varchar(50);not null;default:'General'" json:"category"`
	Priority    string     `gorm:"type:varchar(10);not null;default:'Medium'" json:"priority"`
	CompletedAt *time.Time `json:"completedAt"`
	DueDate     *time.Time `json:"dueDate"`
}

// used pointers to *time.Time so value can be nil and to make comparisons between dates easier than string.
