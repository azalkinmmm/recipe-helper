package recipeprovider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GroqRecipeProvider struct {
	APIKey string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Message Message `json:"message"`
}

type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

func NewGroqRecipeProvider(APIKey string) *GroqRecipeProvider {
	return &GroqRecipeProvider{APIKey: APIKey}
}

const GroqURL = "https://api.groq.com/openai/v1/chat/completions"
const AImodel = "meta-llama/llama-4-scout-17b-16e-instruct"
const systemMessage = "Ты API. Ты возвращаешь ТОЛЬКО валидный JSON без пояснений."
const userMessage = `
У меня есть эти ингредиенты: %s. Какие блюда можно приготовить из этих ингредиентов?

Верни ТОЛЬКО JSON следующего формата:
{
  "dishes": [
    {
      "name": "название блюда",
      "recipe": "рецепт блюда"
    }
  ]
}

Правила:
- Каждый объект ОБЯЗАТЕЛЬНО содержит "name" и "recipe"
- "recipe" — рецепт блюда в формате: "1. первый шаг\n2. второй шаг и т.д."
- В рецепте нельзя использовать ингредиенты, которых не было в запросе
- В рецепте не обязательно использоавть все перечисленные ингредиенты, важно чтобы они сочетались между собой
- Если подходящих блюд нет, верни {"dishes": []}
- Никакого текста вне JSON
`

func (p *GroqRecipeProvider) GetDishes(ingredients string) ([]Dish, error) {
	reqBody := ChatRequest{
		Model: AImodel,
		Messages: []Message{
			{
				Role:    "system",
				Content: systemMessage,
			},
			{
				Role:    "user",
				Content: fmt.Sprintf(userMessage, ingredients),
			},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, GroqURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+p.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("groq API error %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, err
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}

	answer := chatResp.Choices[0].Message.Content

	var dishes DishResponse
	err = json.Unmarshal([]byte(answer), &dishes)
	if err != nil {
		return nil, err
	}

	return dishes.Dishes, nil
}
