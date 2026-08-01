package main

type Todo struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Title     string `gorm:"not null" json:"title"`
	Completed bool   `gorm:"default:false" json:"completed"`
}
