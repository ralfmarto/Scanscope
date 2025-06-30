package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

// AI handles OpenAI calls.
type AI struct {
	apiKey string
	model  string
	client *http.Client
}

// NewAI creates AI client with provided model.
func NewAI(model string) *AI {
	return &AI{
		apiKey: os.Getenv("OPENAI_API_KEY"),
		model:  model,
		client: &http.Client{},
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type choice struct {
	Message chatMessage `json:"message"`
}

type chatResponse struct {
	Choices []choice `json:"choices"`
}

// Validate asks OpenAI if snippet represents risk.
func (a *AI) Validate(prompt, snippet string) (bool, error) {
	if a.apiKey == "" {
		return false, errors.New("missing OPENAI_API_KEY")
	}
	reqBody := chatRequest{
		Model: a.model,
		Messages: []chatMessage{
			{Role: "system", Content: "Você é um sistema de análise de código."},
			{Role: "user", Content: prompt + "\n" + snippet},
		},
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return false, err
	}
	request, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(data))
	if err != nil {
		return false, err
	}
	request.Header.Set("Authorization", "Bearer "+a.apiKey)
	request.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(request)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var r chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return false, err
	}
	if len(r.Choices) == 0 {
		return false, errors.New("no response")
	}
	content := r.Choices[0].Message.Content
	return content != "false" && content != "no", nil
}
