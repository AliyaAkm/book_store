package book

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type BookResponse struct {
	ID     uint    `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	books, err := h.service.GetAllBooks()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := Response{Status: "fail", Message: "Failed to retrieve books"}
		json.NewEncoder(w).Encode(response)
		return
	}

	var bookResponses []BookResponse
	for _, book := range books {
		bookResponses = append(bookResponses, BookResponse{
			ID:     book.ID,
			Title:  book.Title,
			Author: book.Author,
			Price:  book.Price,
		})
	}
	w.WriteHeader(http.StatusOK)
	response := Response{Status: "success", Message: "Books retrieved successfully"}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := Response{Status: "fail", Message: "invalid book ID"}
		json.NewEncoder(w).Encode(response)
		return
	}

	book, err := h.service.GetBookByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := Response{Status: "fail", Message: "Book not found"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := Response{Status: "success", Message: "books retrieved by id successfully", Data: book}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := Response{Status: "fail", Message: "Некорректное JSON-сообщение"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Создание книги с возвратом ID
	if err := h.service.CreateBook(&book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := Response{Status: "fail", Message: "Failed to create book"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := Response{Status: "success", Message: "Book created successfully", Data: book}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := Response{Status: "fail", Message: "Invalid book ID"}
		json.NewEncoder(w).Encode(response)
		return
	}

	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := Response{Status: "fail", Message: "Invalid Input"}
		json.NewEncoder(w).Encode(response)
		return
	}

	book.ID = uint(id)

	if err := h.service.UpdateBook(book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := Response{Status: "fail", Message: "Failed to update book"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := Response{Status: "success", Message: "Book updated successfully", Data: book}
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := Response{Status: "fail", Message: "Invalid book ID"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := h.service.DeleteBook(uint(id)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := Response{Status: "fail", Message: "Failed to delete book"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := Response{Status: "success", Message: "Book deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
