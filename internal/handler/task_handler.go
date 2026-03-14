package handler

import (
    "net/http"
    "github.com/gorilla/mux"
    "your_module/models"
)

// CreateTask handles the creation of a new task
func CreateTask(w http.ResponseWriter, r *http.Request) {
    // Implementation for creating a new task
}

// GetTask handles fetching a task by ID
func GetTask(w http.ResponseWriter, r *http.Request) {
    // Implementation for fetching a task
}

// UpdateTask handles updating an existing task
func UpdateTask(w http.ResponseWriter, r *http.Request) {
    // Implementation for updating a task
}

// DeleteTask handles deleting a task by ID
func DeleteTask(w http.ResponseWriter, r *http.Request) {
    // Implementation for deleting a task
}

// RegisterRoutes registers the task routes
func RegisterRoutes(r *mux.Router) {
    r.HandleFunc("/tasks", CreateTask).Methods("POST")
    r.HandleFunc("/tasks/{id}", GetTask).Methods("GET")
    r.HandleFunc("/tasks/{id}", UpdateTask).Methods("PUT")
    r.HandleFunc("/tasks/{id}", DeleteTask).Methods("DELETE")
}