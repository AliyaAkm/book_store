package user

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username" gorm:"uniqueIndex"`
	Password string `json:"password"`
}
