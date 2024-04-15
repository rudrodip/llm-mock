package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type APIServer struct {
	listenAddr string
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{listenAddr}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", makeHTTPHandleFunc(s.handlePing))
	mux.HandleFunc("POST /chat/completions", makeHTTPHandleFunc(s.handleCompletions))

	log.Println("API server listening on", s.listenAddr)
	http.ListenAndServe(s.listenAddr, mux)
}

func (s *APIServer) handlePing(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, map[string]string{"message": "pong"})
}

func (s *APIServer) handleCompletions(w http.ResponseWriter, r *http.Request) error {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	res := Response{
		Id:     uuid.New().String(),
		Object: "chat.completion",
		Model:  "gpt-3.5-turbo",
		Usage: Usage{
			Prompt_tokens:     len(req.Messages),
			Completion_tokens: 100,
			Total_tokens:      len(req.Messages) + 100,
		},
		Choices: []Choice{
			{
				Message: Message{
					Role:    "assistant",
					Content: req.Messages[0].Content + "Hello, how can I help you today?",
				},
				Logprobs:      nil,
				Finish_reason: "stop",
				Index:         0,
			},
		},
	}

	return WriteJSON(w, http.StatusOK, res)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
