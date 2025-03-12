package book

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllBooks(ctx context.Context) ([]Book, error) {
	return s.repo.GetAllBooks(ctx)
}
func (s *Service) GetBookByID(ctx context.Context, id uint) (*Book, error) {
	return s.repo.GetBookByID(ctx, id)
}

func (s *Service) CreateBook(ctx context.Context, book *Book) error {
	return s.repo.CreateBook(ctx, book)
}

func (s *Service) UpdateBook(ctx context.Context, book *Book) error {
	return s.repo.UpdateBook(ctx, book)
}

func (s *Service) DeleteBook(ctx context.Context, id uint) error {
	return s.repo.DeleteBook(ctx, id)
}
