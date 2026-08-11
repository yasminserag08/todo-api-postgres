package main

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null" json:"-"` // json:"-" so that hash never appears in API responses
	Role     string `gorm:"default:user"`
}
