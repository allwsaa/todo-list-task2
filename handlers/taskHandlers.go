package handlers

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/allwsaa/todo-list-task2/models"
	"github.com/araddon/dateparse"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var tasks = map[string]*models.Task{}
var tasksLock sync.Mutex

// CreateTask godoc
// @Summary Create a new task
// @Description Creates a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task to create"
// @Success 201 {object} map[string]string "id of the created task"
// @Failure 400 {string} string "Bad request"
// @Router /api/todo-list/tasks [post]
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

// UpdateTask godoc
// @Summary Update an existing task
// @Description Updates an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param ID path string true "Task ID"
// @Param task body models.Task true "Updated task data"
// @Success 204 {string} string "Task updated successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Task not found"
// @Router /api/todo-list/tasks/{ID} [put]
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

// CompleteTask godoc
// @Summary Complete a task
// @Description Marks a task as done
// @Tags tasks
// @Param ID path string true "Task ID"
// @Success 204 {string} string "Task completed successfully"
// @Failure 404 {string} string "Task not found"
// @Router /api/todo-list/tasks/{ID}/done [put]
func CompleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")
	tasksLock.Lock()
	defer tasksLock.Unlock()

	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	task.Done = true
	tasks[id] = task

	w.WriteHeader(http.StatusNoContent)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Deletes a task
// @Tags tasks
// @Param ID path string true "Task ID"
// @Success 204 {string} string "Task deleted successfully"
// @Failure 404 {string} string "Task not found"
// @Router /api/todo-list/tasks/{ID} [delete]
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")
	tasksLock.Lock()
	defer tasksLock.Unlock()

	if _, exists := tasks[id]; !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusNoContent)
}

// GetTask godoc
// @Summary Get a task by ID
// @Description Retrieves a task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param ID path string true "Task ID"
// @Success 200 {object} models.Task "Task details"
// @Failure 404 {string} string "Task not found"
// @Router /api/todo-list/tasks/{ID} [get]

func GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")
	tasksLock.Lock()
	defer tasksLock.Unlock()

	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// GetTasks godoc
// @Summary Get tasks
// @Description Retrieves a list of tasks
// @Tags tasks
// @Param status query string false "Task status" Enums(active, done) default(active)
// @Accept json
// @Produce json
// @Success 200 {array} models.Task "List of tasks"
// @Router /api/todo-list/tasks [get]
func GetTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status != "active" && status != "done" {
		status = "active"
	}
	tasksLock.Lock()
	defer tasksLock.Unlock()

	var res []*models.Task
	now := time.Now()

	for _, task := range tasks {
		activeAt, err := time.Parse("2006-01-02", task.ActiveAt)
		if err != nil {
			continue
		}

		if status == "done" && task.Done {
			res = append(res, task)
		} else if status == "active" && !task.Done && (activeAt.Before(now) || activeAt.Equal(now)) {
			res = append(res, task)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ActiveAt < res[j].ActiveAt
	})

	weekday := now.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		for _, task := range res {
			task.Title = "ВЫХОДНОЙ - " + task.Title
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
