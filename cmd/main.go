package main

import (
	"book_store/internal/auth"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"book_store/config"
	"book_store/internal/book"
	"book_store/internal/user"

	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()
	db := config.DB

	router := mux.NewRouter()

	ctx, cancel := context.WithCancel(context.Background())

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

	go bookHandler.BackgroundTask(ctx)

	authRoutes := router.PathPrefix("/api").Subrouter()
	authRoutes.Use(authMiddleware)

	// Book routes
	authRoutes.HandleFunc("/books", bookHandler.GetAllBooks).Methods("GET")
	authRoutes.HandleFunc("/books/{id}", bookHandler.GetBookByID).Methods("GET")
	authRoutes.HandleFunc("/books", bookHandler.CreateBook).Methods("POST")
	authRoutes.HandleFunc("/books/{id}", bookHandler.UpdateBook).Methods("PUT")
	authRoutes.HandleFunc("/books/{id}", bookHandler.DeleteBook).Methods("DELETE")
	authRoutes.HandleFunc("/cancelable", bookHandler.CancelableOperation).Methods("GET")
	authRoutes.HandleFunc("/limited", bookHandler.LimitedTimeOperation).Methods("GET")

	// User routes
	authRoutes.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	authRoutes.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	authRoutes.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	authRoutes.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	authRoutes.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	server := &http.Server{
		Addr:    ":8084",
		Handler: router,
	}

	go func() {
		// Ctrl+C
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		log.Println("Shutting down server...")

		// Отменяем контекст для фоновой задачи
		cancel()

		ctxShutdown, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(ctxShutdown); err != nil {
			log.Fatal("Server shutdown failed:", err)
		}
		log.Println("Server stopped gracefully")
	}()

	log.Println("Server is running on port 8084")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed:", err)
	}
}
