package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type todo struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {

	initDB()
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	SetupRoutes(r)
	// Call generateContentFromText
	todo, err := generateToDoFromText(os.Stdout, "Buy milk")
	fmt.Printf("the todo from llm is:  %+v\n", todo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating content: %v\n", err)
	}

	fmt.Println("Server is running on port 3000")
	http.ListenAndServe(":3000", r)
}
