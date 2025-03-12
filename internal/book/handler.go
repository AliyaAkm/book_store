package book

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
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

// все круды с контекстом WithTimeout.если не завершится за время, то автоматически закроется запрос
func (h *Handler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	books, err := h.service.GetAllBooks(ctx)
	if err != nil {
		log.Println("Error retrieving books:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: "fail", Message: "Failed to retrieve books"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: "success", Message: "Books retrieved successfully", Data: books})
}

func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: "fail", Message: "Invalid book ID"})
		return
	}

	book, err := h.service.GetBookByID(ctx, uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{Status: "fail", Message: "Book not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: "success", Message: "Book retrieved successfully", Data: book})
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: "fail", Message: "Invalid JSON format"})
		return
	}

	if err := h.service.CreateBook(ctx, &book); err != nil {
		log.Println("Error creating book:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: "fail", Message: "Failed to create book"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Status: "success", Message: "Book created successfully", Data: book})
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
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

	if err := h.service.UpdateBook(ctx, &book); err != nil {
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
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: "fail", Message: "Invalid book ID"})
		return
	}

	if err := h.service.DeleteBook(ctx, uint(id)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: "fail", Message: "Failed to delete book"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: "success", Message: "Book deleted successfully"})
}

// контекст с Background для фоновых задач
func (h *Handler) BackgroundTask(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Background task executed")
		case <-ctx.Done():
			log.Println("Background task stopped")
			return
		}
	}
}

// контекст с WithCancel. если запрос не завершается за 2 сек, то она принудительно закроется вручную
func (h *Handler) CancelableOperation(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()

	select {
	case <-time.After(2 * time.Second):
		w.WriteHeader(http.StatusOK)
		response := Response{Status: "success", Message: "done!"}
		json.NewEncoder(w).Encode(response)
	case <-ctx.Done():
		w.WriteHeader(http.StatusRequestTimeout)
		json.NewEncoder(w).Encode(Response{Status: "fail", Message: "operation cancelled"})

	}
}

// контекст с WithDeadline. если не завершится за 4 сек, пройдет к блоку ctx.Done
func (h *Handler) LimitedTimeOperation(w http.ResponseWriter, r *http.Request) {
	deadline := time.Now().Add(4 * time.Second)
	ctx, cancel := context.WithDeadline(r.Context(), deadline)
	defer cancel()

	select {
	case <-time.After(3 * time.Second):
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{Status: "success", Message: "Successfully task completed!"})
	case <-ctx.Done():
		w.WriteHeader(http.StatusRequestTimeout)
		json.NewEncoder(w).Encode(Response{Status: "fail", Message: "task deadline exceeded"})
	}
}

// контекст с TODO
func (h *Handler) FutureFunction(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: "info", Message: "Feature under development", Data: ctx})
}
