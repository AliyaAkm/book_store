package user

import (
	"errors"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllUsers() ([]User, error) {
	var users []User
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *Repository) GetUserByID(id uint) (*User, error) {
	var user User
	result := r.db.First(&user, id)
	return &user, result.Error
}

func (r *Repository) CreateUser(user *User) error {
	return r.db.Create(&user).Error
}

func (r *Repository) UpdateUser(user User) error {
	result := r.db.Model(&User{}).Where("id = ?", user.ID).Updates(user)
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return result.Error
}

func (r *Repository) DeleteUser(id uint) error {
	return r.db.Delete(&User{}, id).Error
}
