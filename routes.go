package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var todoListItems todoList

func SetupRoutes(r *chi.Mux) {
	r.Get("/", homeHandler)
	r.Post("/add", addTodoHandler)
	r.Delete("/delete/{id}", deleteTodoHandler)
	r.Put("/update/{id}", updateTodoHandler)
	r.Get("/todos", getTodosHandler)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome to my todo app"))
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST request received")
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	var newTodo todo
	newTodo.Title = r.FormValue("title")
	if newTodo.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	fmt.Println("Title:", newTodo.Title)
	newTodo.Done, _ = strconv.ParseBool(r.FormValue("done"))
	fmt.Println("Done:", newTodo.Done)
	newTodo.ID = len(todoListItems) + 1

	fmt.Printf("New todo: %+v\n", newTodo)

	todoListItems = append(todoListItems, newTodo)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)

}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for i, v := range todoListItems {
		if v.ID == idInt {
			todoListItems = append(todoListItems[:i], todoListItems[i+1:]...)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func getTodosHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("%+v", todoListItems)))
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for i, v := range todoListItems {
		if v.ID == idInt {
			todoListItems[i].Done = !todoListItems[i].Done
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
