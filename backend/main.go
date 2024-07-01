package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

const (
	PORT = "3000"
)

var db *sql.DB

type Todo struct {
	Id   string `json:"id"`
	Todo string `json:"todo"`
}

type Response struct {
	Message string `json:"message"`
}

func AddTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var todo Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Fatalf("ERROR: could not read request payload: %v", err.Error())
	}

	//---------------------SQL start-----------------------------------------------//
	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("ERROR: could not generate a UUID: %v", err.Error())
	}
	todo.Id = uuid.String()

	stmt, err := db.Prepare(`INSERT INTO todos(id, todo) VALUES(?, ?)`)
	if err != nil {
		log.Fatalf("ERROR: could not prepare query: %v", err.Error())
	}
	defer stmt.Close()

	_, query_err := stmt.Exec(todo.Id, todo.Todo)
	if query_err != nil {
		log.Fatalf("ERROR: could not add new todo: %v", err.Error())
	}
	//---------------------SQL end-----------------------------------------------//

	json.NewEncoder(w).Encode(todo)
}

func GetTodoHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todos []Todo

	//---------------------SQL start-----------------------------------------------//
	sql_stmt := `SELECT * FROM todos;`
	result, err := db.Query(sql_stmt)
	if err != nil {
		log.Fatalf("ERROR: could not fetch todos: %v", err.Error())
	}
	defer result.Close()

	for result.Next() {
		var id string
		var todo string

		scan_err := result.Scan(&id, &todo)
		if scan_err != nil {
			log.Fatalf("ERROR: could not extract todos data: %v", err.Error())
		}

		todos = append(todos, Todo{id, todo})
	}

	result_err := result.Err()
	if result_err != nil {
		log.Fatalf("ERROR: cannot complete the iteration: %v", result_err.Error())
	}
	//---------------------SQL end-----------------------------------------------//

	json.NewEncoder(w).Encode(todos)
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	todo_id := vars["todo"]

	var todo Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Fatalf("ERROR: could not read request payload: %v", err.Error())
	}

	stmt, err := db.Prepare(`UPDATE todos SET todo = ? WHERE id = ?`)
	if err != nil {
		log.Fatalf("ERROR: could not prepare query: %v", err.Error())
	}
	defer stmt.Close()

	_, query_err := stmt.Exec(todo.Todo, todo_id)
	if query_err != nil {
		log.Fatalf("ERROR: could not update todo: %v", err.Error())
	}

	response := Response{Message: "Todo updated successfully."}
	json.NewEncoder(w).Encode(response)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	todo_id := vars["todo"]

	stmt, err := db.Prepare(`DELETE FROM todos WHERE id = ?`)
	if err != nil {
		log.Fatalf("ERROR: could not prepare query: %v", err.Error())
	}
	defer stmt.Close()

	_, query_err := stmt.Exec(todo_id)
	if query_err != nil {
		log.Fatalf("ERROR: could not delete todo: %v", err.Error())
	}

	response := Response{Message: "Todo deleted successfully."}
	json.NewEncoder(w).Encode(response)
}

func MarkTodoCompleteHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Complete todo")
	log.Printf("INFO: mark todo complete is hit\n")
}

func main() {
	var db_err error
	db, db_err = sql.Open("sqlite3", "./gotodo.db")
	if db_err != nil {
		log.Fatalf("ERROR: cannot connect to database: %v", db_err.Error())
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/api/todos", GetTodoHandler).Methods("GET")
	router.HandleFunc("/api/todos", AddTodoHandler).Methods("POST")
	router.HandleFunc("/api/todos/{todo}", UpdateTodoHandler).Methods("PUT")
	router.HandleFunc("/api/todos/{todo}", DeleteTodoHandler).Methods("DELETE")
	router.HandleFunc("/api/todos/{todo}", MarkTodoCompleteHandler).Methods("PATCH")

	log.Printf("INFO: server is listening on port: %v\n", PORT)

	err := http.ListenAndServe(":"+PORT, router)
	if err != nil {
		log.Fatalf("ERROR: Could not start the server on port: %v\n", PORT)
	}
}
