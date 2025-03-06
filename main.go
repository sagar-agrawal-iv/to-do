package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Task struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

var tasks = []Task{}
var mu sync.Mutex

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func addTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task Task
	json.NewDecoder(r.Body).Decode(&task)

	mu.Lock()
	tasks = append(tasks, task)
	mu.Unlock()

	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	mu.Lock()
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}
	mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func toggleTaskDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	mu.Lock()
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Done = !tasks[i].Done
			json.NewEncoder(w).Encode(tasks[i])
			mu.Unlock()
			return
		}
	}
	mu.Unlock()

	http.Error(w, "Task not found", http.StatusNotFound)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", addTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}/done", toggleTaskDone).Methods("PUT")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	http.ListenAndServe(":8080", handler)
}
