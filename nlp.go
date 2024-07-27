package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/vertexai/genai"
	"github.com/joho/godotenv"
)

func generateToDoFromText(w io.Writer, todoFromUser string) (todo, error) {
	if err := godotenv.Load(); err != nil {
		return todo{}, fmt.Errorf("error loading .env file: %w", err)
	}
	projectID := os.Getenv("PROJECT_ID")
	location := os.Getenv("LOCATION")
	modelName := os.Getenv("MODEL_NAME")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return todo{}, fmt.Errorf("error creating client: %w", err)
	}
	gemini := client.GenerativeModel(modelName)
	text := fmt.Sprintf("Turn this user todo text %s into a struct with keys title, completed as boolean of false, only give output and no explanation or any template literals like ```json or ```", todoFromUser)
	prompt := genai.Text(text)

	resp, err := gemini.GenerateContent(ctx, prompt)
	if err != nil {
		return todo{}, fmt.Errorf("error generating content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return todo{}, fmt.Errorf("no content generated")
	}

	jsonStr, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return todo{}, fmt.Errorf("unexpected response format")
	}

	fmt.Printf("The response from the llm is: %+v\n", jsonStr)

	var newTodo todo
	err = json.Unmarshal([]byte(jsonStr), &newTodo)
	if err != nil {
		return todo{}, fmt.Errorf("error unmarshalling json: %w", err)
	}

	rb, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return todo{}, fmt.Errorf("json.MarshalIndent: %w", err)
	}
	fmt.Fprintln(w, string(rb))
	return newTodo, nil
}
