package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux) {
	r.Post("/add", addTodoHandler)
	r.Delete("/delete/{id}", deleteTodoHandler)
	r.Put("/update/{id}", updateTodoHandler)
	r.Get("/todos", getTodosHandler)
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST request received")
	var newTodo todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//insert into db
	result, err := db.Exec("INSERT INTO todos(title, completed) VALUES (?, ?)", newTodo.Title, newTodo.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	newTodo.Id = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("DELETE FROM todos WHERE id = ?", idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Todo deleted"))
}

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var todos []todo
	for rows.Next() {
		var t todo
		rows.Scan(&t.Id, &t.Title, &t.Completed)
		todos = append(todos, t)
	}
	json.NewEncoder(w).Encode(todos)
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Fetch the current todo item
	var todo todo
	err = db.QueryRow("SELECT completed FROM todos WHERE id = ?", idInt).Scan(&todo.Completed)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Todo not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Toggle the completed status
	newStatus := !todo.Completed

	// Update the todo in the database
	_, err = db.Exec("UPDATE todos SET completed = ? WHERE id = ?", newStatus, idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Todo updated"))
}
