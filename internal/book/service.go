package book

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllBooks() ([]Book, error) {
	return s.repo.GetAllBooks()
}

func (s *Service) GetBookByID(id uint) (*Book, error) {
	return s.repo.GetBookByID(id)
}

func (s *Service) CreateBook(book *Book) error {
	return s.repo.CreateBook(book)
}
func (s *Service) UpdateBook(book Book) error {
	existingBook, err := s.repo.GetBookByID(book.ID)
	if err != nil {
		return err
	}

	existingBook.Title = book.Title
	existingBook.Author = book.Author
	existingBook.Price = book.Price

	return s.repo.UpdateBook(*existingBook)
}

func (s *Service) DeleteBook(id uint) error {
	return s.repo.DeleteBook(id)
}
