package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type todo struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type todoList []todo

func main() {

	initDB()
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	SetupRoutes(r)

	fmt.Println("Server is running on port 3000")
	http.ListenAndServe(":3000", r)
}
