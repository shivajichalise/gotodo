package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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
	Id          string `json:"id"`
	Todo        string `json:"todo"`
	IsCompleted bool   `json:"is_completed"`
}

type Response struct {
	Message string `json:"message"`
}

func respond(w http.ResponseWriter, message string, status int) {
	response := Response{Message: message}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func todoExists(todo_id string) (bool, error) {
	count := 0

	//---------------------SQL start-----------------------------------------------//
	stmt, err := db.Prepare(`SELECT COUNT(*) FROM todos WHERE id = ?`)
	if err != nil {
		return false, errors.New("ERROR: could not check for todo with given id")
	}
	defer stmt.Close()

	err = stmt.QueryRow(todo_id).Scan(&count)
	if err != nil {
		return false, errors.New("ERROR: could not query for given todo")
	}
	//---------------------SQL end-----------------------------------------------//

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func AddTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var todo Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Printf("ERROR: could not read request payload: %s\n", err.Error())
		respond(w, "Something went wrong, please try again.", http.StatusInternalServerError)
		return
	}

	//---------------------SQL start-----------------------------------------------//
	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Printf("ERROR: could not generate a UUID: %s\n", err.Error())
		respond(w, "Something went wrong, please try again.", http.StatusInternalServerError)
		return
	}
	todo.Id = uuid.String()

	stmt, err := db.Prepare(`INSERT INTO todos(id, todo) VALUES(?, ?)`)
	if err != nil {
		log.Printf("ERROR: could not prepare query: %s\n", err.Error())
		respond(w, "Something went wrong, please try again.", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, query_err := stmt.Exec(todo.Id, todo.Todo)
	if query_err != nil {
		log.Printf("ERROR: could not add new todo: %s\n", err.Error())
		respond(w, "Could not add a new todo, please try again.", http.StatusInternalServerError)
		return
	}
	//---------------------SQL end-----------------------------------------------//

	json.NewEncoder(w).Encode(todo)
}

func GetTodoHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todos []Todo

	//---------------------SQL start-----------------------------------------------//
	result, err := db.Query("SELECT id, todo, is_completed FROM todos")
	if err != nil {
		log.Printf("ERROR: could not fetch todos: %s\n", err.Error())
		respond(w, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}
	defer result.Close()

	for result.Next() {
		var id string
		var todo string
		var is_completed bool

		scan_err := result.Scan(&id, &todo, &is_completed)
		if scan_err != nil {
			log.Printf("ERROR: could not extract todos data: %s\n", err.Error())
			respond(w, "Something went wrong, please try again", http.StatusInternalServerError)
			return
		}

		todos = append(todos, Todo{id, todo, is_completed})
	}

	result_err := result.Err()
	if result_err != nil {
		log.Printf("ERROR: cannot complete the iteration: %s\n", result_err.Error())
		respond(w, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}
	//---------------------SQL end-----------------------------------------------//

	if len(todos) == 0 {
		json.NewEncoder(w).Encode(make([]string, 0))
		return
	}

	json.NewEncoder(w).Encode(todos)
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	todo_id := vars["todo"]

	exists, err := todoExists(todo_id)
	if err != nil {
		log.Printf("ERROR: %s\n", err)

		respond(w, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}

	if !exists {
		respond(w, "Todo not found", http.StatusNotFound)
		return
	}

	var todo Todo

	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		log.Printf("ERROR: could not read request payload: %s\n", err.Error())
		respond(w, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}

	//---------------------SQL start-----------------------------------------------//
	stmt, err := db.Prepare(`UPDATE todos SET todo = ? WHERE id = ?`)
	if err != nil {
		log.Printf("ERROR: could not prepare query: %s\n", err.Error())
		respond(w, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, query_err := stmt.Exec(todo.Todo, todo_id)
	if query_err != nil {
		log.Printf("ERROR: could not update todo: %s\n", err.Error())
		respond(w, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}
	//---------------------SQL end-----------------------------------------------//

	response := Response{Message: "Todo updated successfully."}
	json.NewEncoder(w).Encode(response)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	todo_id := vars["todo"]

	exists, err := todoExists(todo_id)
	if err != nil {
		log.Printf("ERROR: %s\n", err)

		respond(w, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}

	if !exists {
		respond(w, "Todo not found", http.StatusNotFound)
		return
	}

	stmt, err := db.Prepare(`DELETE FROM todos WHERE id = ?`)
	if err != nil {
		log.Printf("ERROR: could not prepare query: %s\n", err.Error())
		respond(w, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, query_err := stmt.Exec(todo_id)
	if query_err != nil {
		log.Printf("ERROR: could not delete todo: %s\n", err.Error())
		respond(w, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}

	response := Response{Message: "Todo deleted successfully."}
	json.NewEncoder(w).Encode(response)
}

func MarkTodoCompleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	todo_id := vars["todo"]

	exists, err := todoExists(todo_id)
	if err != nil {
		log.Printf("ERROR: %s\n", err)

		respond(w, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}

	if !exists {
		respond(w, "Todo not found", http.StatusNotFound)
		return
	}

	//---------------------SQL start-----------------------------------------------//
	stmt, err := db.Prepare(`UPDATE todos SET is_completed = true WHERE id = ?`)
	if err != nil {
		log.Printf("ERROR: could not prepare query: %s\n", err.Error())
		respond(w, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, query_err := stmt.Exec(todo_id)
	if query_err != nil {
		log.Printf("ERROR: could not mark todo as complete: %s\n", err.Error())
		respond(w, "Something went wrong, please try again", http.StatusInternalServerError)
		return
	}
	//---------------------SQL end-----------------------------------------------//

	response := Response{Message: "Todo marked as completed successfully."}
	json.NewEncoder(w).Encode(response)
}

func main() {
	var db_err error
	db, db_err = sql.Open("sqlite3", "./gotodo.db")
	if db_err != nil {
		log.Fatalf("ERROR: cannot connect to database: %s\n", db_err.Error())
	}
	defer db.Close()

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS todos (id STRING NOT NULL PRIMARY KEY, todo TEXT, is_completed BOOLEAN DEFAULT FALSE);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("ERROR: could not create todos table: %s\n", err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/todos", GetTodoHandler).Methods("GET")
	router.HandleFunc("/api/todos", AddTodoHandler).Methods("POST")
	router.HandleFunc("/api/todos/{todo}", UpdateTodoHandler).Methods("PUT")
	router.HandleFunc("/api/todos/{todo}", DeleteTodoHandler).Methods("DELETE")
	router.HandleFunc("/api/todos/{todo}/complete", MarkTodoCompleteHandler).Methods("PATCH")

	log.Printf("INFO: server is listening on port: %s\n", PORT)

	err = http.ListenAndServe(":"+PORT, router)
	if err != nil {
		log.Fatalf("ERROR: Could not start the server on port: %s\n", PORT)
	}
}
