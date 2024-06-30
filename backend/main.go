package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	PORT = "3000"
)

type Todo struct {
	Id   string `json:"id"`
	Todo string `json:"todo"`
}

var todos []Todo

func add_todo(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todo := Todo{Id: "1", Todo: "Test"}
	todos = append(todos, todo)
	log.Printf("INFO: add todo is hit\n")
	json.NewEncoder(w).Encode(todos)
}

func get_todos(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("INFO: get todos is hit\n")
	json.NewEncoder(w).Encode(todos)
}

func update_todo(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Update todo")
	log.Printf("INFO: update todo is hit\n")
}

func delete_todo(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Delete todo")
	log.Printf("INFO: delete todo is hit\n")
}

func mark_todo_as_complete(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Complete todo")
	log.Printf("INFO: mark todo complete is hit\n")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/todos", get_todos).Methods("GET")
	router.HandleFunc("/api/todos", add_todo).Methods("POST")
	router.HandleFunc("/api/todos/{todo}", update_todo).Methods("PUT")
	router.HandleFunc("/api/todos/{todo}", delete_todo).Methods("DELETE")
	router.HandleFunc("/api/todos/{todo}", mark_todo_as_complete).Methods("PATCH")

	log.Printf("INFO: server is listening on port: %v\n", PORT)

	err := http.ListenAndServe(":"+PORT, router)
	if err != nil {
		log.Fatalf("ERROR: Could not start the server on port: %v\n", PORT)
	}
}
