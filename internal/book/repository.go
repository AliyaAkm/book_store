package book

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

func (r *Repository) GetAllBooks() ([]Book, error) {
	var books []Book
	result := r.db.Find(&books)
	return books, result.Error
}

func (r *Repository) GetBookByID(id uint) (*Book, error) {
	var book Book
	result := r.db.First(&book, id)
	return &book, result.Error
}

func (r *Repository) CreateBook(book *Book) error {
	return r.db.Create(book).Error
}

func (r *Repository) UpdateBook(book Book) error {
	// Обновляем только те поля, которые были переданы
	result := r.db.Model(&Book{}).Where("id = ?", book.ID).Updates(book)
	if result.RowsAffected == 0 {
		return errors.New("book not found")
	}
	return result.Error
}

func (r *Repository) DeleteBook(id uint) error {
	return r.db.Delete(&Book{}, id).Error
}
