package user

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllUsers() ([]User, error) {
	return s.repo.GetAllUsers()
}

func (s *Service) GetUserByID(id uint) (*User, error) {
	return s.repo.GetUserByID(id)
}

func (s *Service) CreateUser(user *User) error {
	return s.repo.CreateUser(user)
}

func (s *Service) UpdateUser(user User) error {
	existingUser, err := s.repo.GetUserByID(user.ID)
	if err != nil {
		return err // Если пользователь не найден
	}

	existingUser.Username = user.Username
	existingUser.Password = user.Password

	return s.repo.UpdateUser(*existingUser)
}

func (s *Service) DeleteUser(id uint) error {
	return s.repo.DeleteUser(id)
}
