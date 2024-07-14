package main

import (
	"log"
	"net/http"

	"github.com/allwsaa/todo-list-task2/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the Todo List Service"))
	})

	r.Post("/api/todo-list/tasks", handlers.CreateTask)
	r.Put("/api/todo-list/tasks/{ID}", handlers.UpdateTask)
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}

}
