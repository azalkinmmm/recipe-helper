package main

import (
	"log"
	"net/http"
	"os"

	"github.com/azalkinmmm/recipe-helper/internal/handler"
	"github.com/azalkinmmm/recipe-helper/internal/recipeprovider"
)

const PORT = ":8085"

func main() {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Fatal("api key not found")
	}

	recipeProvider := recipeprovider.NewGroqRecipeProvider(apiKey)
	handler := handler.GetDichesHandler(recipeProvider)

	log.Println("Starting server on port", PORT)
	err := http.ListenAndServe(PORT, handler)
	if err != nil {
		log.Fatal(err)
	}
}
