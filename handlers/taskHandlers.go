package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/allwsaa/todo-list-task2/models"
	"github.com/araddon/dateparse"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var tasks = map[string]*models.Task{}
var tasksLock sync.Mutex

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title    string `json:"title"`
		ActiveAt string `json:"activeAt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Title) > 200 {
		http.Error(w, "Given title is too long", http.StatusBadRequest)
		return
	}

	_, err := dateparse.ParseAny(req.ActiveAt)
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}
	tasksLock.Lock()
	defer tasksLock.Unlock()

	for _, task := range tasks {
		if task.Title == req.Title && task.ActiveAt == req.ActiveAt {
			http.Error(w, "This task already exists", http.StatusConflict)
			return
		}
	}

	taskNew := &models.Task{
		ID:       uuid.New().String(),
		Title:    req.Title,
		ActiveAt: req.ActiveAt,
		Done:     false,
	}

	tasks[taskNew.ID] = taskNew

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": taskNew.ID})

}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")

	var req struct {
		Title    string `json:"title"`
		ActiveAt string `json:"activeAt"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Title) > 200 {
		http.Error(w, "Given title is too long", http.StatusBadRequest)
		return
	}
	_, err := dateparse.ParseAny(req.ActiveAt)
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}
	tasksLock.Lock()
	defer tasksLock.Unlock()

	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	for _, existingTask := range tasks {
		if existingTask.ID != id && existingTask.Title == req.Title && existingTask.ActiveAt == req.ActiveAt {
			http.Error(w, "Task with this title and date already exists", http.StatusConflict)
			return
		}
	}

	task.Title = req.Title
	task.ActiveAt = req.ActiveAt

	tasks[id] = task

	w.WriteHeader(http.StatusNoContent)

}
