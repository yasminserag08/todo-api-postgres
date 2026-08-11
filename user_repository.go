package main

import "gorm.io/gorm"

type UserRepositoryInterface interface {
	CreateUser(user User) (User, error)
	GetUserByUsername(username string) (User, error)
	GetUserByID(id uint) (User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user User) (User, error) {
	result := r.db.Create(&user)
	return user, result.Error
}

func (r *UserRepository) GetUserByUsername(username string) (User, error) {
	user := User{}
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return user, ErrNotFound
	}
	return user, result.Error
}

func (r *UserRepository) GetUserByID(id uint) (User, error) {
	user := User{}
	result := r.db.First(&user, id)
	if result.Error == gorm.ErrRecordNotFound {
		return user, ErrNotFound
	}
	return user, result.Error
}
