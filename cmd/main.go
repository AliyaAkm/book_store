package main

import (
	"log"
	"net/http"

	"book_store/config"
	"book_store/internal/auth"
	"book_store/internal/book"
	"book_store/internal/user"

	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()
	db := config.DB

	router := mux.NewRouter()

	// мидлвар для аутентификации
	authMiddleware := auth.NewAuthMiddleware(db)

	// User Handlers
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// Book Handlers
	bookRepo := book.NewRepository(db)
	bookService := book.NewService(bookRepo)
	bookHandler := book.NewHandler(bookService)

	authRoutes := router.PathPrefix("/api").Subrouter()
	authRoutes.Use(authMiddleware)

	// Book routes
	authRoutes.HandleFunc("/books", bookHandler.GetAllBooks).Methods("GET")
	authRoutes.HandleFunc("/books/{id}", bookHandler.GetBookByID).Methods("GET")
	authRoutes.HandleFunc("/books", bookHandler.CreateBook).Methods("POST")
	authRoutes.HandleFunc("/books/{id}", bookHandler.UpdateBook).Methods("PUT")
	authRoutes.HandleFunc("/books/{id}", bookHandler.DeleteBook).Methods("DELETE")

	// User routes
	authRoutes.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	authRoutes.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	authRoutes.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	authRoutes.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	authRoutes.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Запуск сервера
	log.Println("Server is running on port 8084")
	log.Fatal(http.ListenAndServe(":8084", router))
}
