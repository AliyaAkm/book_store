package book

type Book struct {
	ID     uint    `json:"id"` // Указываем ID для JSON
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}
