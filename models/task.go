package models

type Task struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	ActiveAt string `json:"activeAt"`
	Done     bool   `json:"done"`
}
