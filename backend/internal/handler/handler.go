package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/azalkinmmm/recipe-helper/internal/recipeprovider"
)

type RecipeProvider interface {
	GetDishes(ingredients string) (dishes []recipeprovider.Dish, err error)
}

type Query struct {
	Message string `json:"message"`
}

func GetDichesHandler(provider RecipeProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"message":"only POST http method is supported"}`, http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		var query Query
		err := json.NewDecoder(r.Body).Decode(&query)
		if err != nil {
			http.Error(w, `{"message":"failed to unmarshall request body"}`, http.StatusInternalServerError)
			return
		}

		if query.Message == "" {
			http.Error(w, `{"message":"bad request body"}`, http.StatusBadRequest)
			return
		}

		dishes, err := provider.GetDishes(query.Message)
		if err != nil {
			msg := fmt.Sprintf(`{"message":"%s"}`, err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(&dishes); err != nil {
			http.Error(w, `{"message":"error while encoding response"}`, http.StatusInternalServerError)
			return
		}
	})
}
