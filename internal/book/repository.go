package book

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAllBooks(ctx context.Context) ([]Book, error) {
	var books []Book
	result := r.db.WithContext(ctx).Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return books, nil
}

func (r *Repository) GetBookByID(ctx context.Context, id uint) (*Book, error) {
	var book Book
	result := r.db.WithContext(ctx).First(&book, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("book not found")
	}
	return &book, result.Error
}

func (r *Repository) CreateBook(ctx context.Context, book *Book) error {
	return r.db.WithContext(ctx).Create(book).Error
}

func (r *Repository) UpdateBook(ctx context.Context, book *Book) error {
	result := r.db.WithContext(ctx).Model(&Book{}).Where("id = ?", book.ID).Updates(book)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("book not found")
	}
	return nil
}

func (r *Repository) DeleteBook(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&Book{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("book not found")
	}
	return nil
}
