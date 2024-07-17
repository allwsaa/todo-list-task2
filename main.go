package main

import (
	"log"
	"net/http"

	_ "github.com/allwsaa/todo-list-task2/docs"
	"github.com/allwsaa/todo-list-task2/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @Title TODO List API
// @Version 1.0
// @Description This is a simple TODO List API (hl)
// BasePath
func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the Todo List Service"))
	})
	//swagger endpoint
	r.Get("/swagger/*", httpSwagger.Handler())
	//healthcheck endpoint
	r.Get("/health", handlers.HealthCheckHandler)

	r.Post("/api/todo-list/tasks", handlers.CreateTask)
	r.Put("/api/todo-list/tasks/{ID}", handlers.UpdateTask)
	r.Delete("/api/todo-list/tasks/{ID}", handlers.DeleteTask)
	r.Put("/api/todo-list/tasks/{ID}/done", handlers.CompleteTask)
	r.Get("/api/todo-list/tasks", handlers.GetTasks)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}

}
