package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
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
	mux.HandleFunc("POST /chat/completions/streaming", makeHTTPHandleFunc(s.handleStreamingCompletions))

	server := &http.Server{
		Addr:    s.listenAddr,
		Handler: mux,
	}

	log.Println("API server listening on", s.listenAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
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
		Id:     randomIdGenerator(),
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

func (s *APIServer) handleStreamingCompletions(w http.ResponseWriter, r *http.Request) error {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	initialResponse := Response{
		Id:      randomIdGenerator(),
		Object:  "chat.completion",
		Model:   "gpt-3.5-turbo",
		Usage:   Usage{Prompt_tokens: len(req.Messages), Completion_tokens: 0, Total_tokens: len(req.Messages)},
		Choices: []Choice{},
	}
	if err := writeSSE(w, "data: ", initialResponse); err != nil {
		return err
	}
	flusher.Flush()

	for i := 0; i < 10; i++ {
		choice := Choice{
			Message:       Message{Role: "assistant", Content: fmt.Sprintf("Streamed response part %d", i)},
			Logprobs:      nil,
			Finish_reason: "stop",
			Index:         0,
		}
		if err := writeSSE(w, "data: ", choice); err != nil {
			return err
		}
		flusher.Flush()
		time.Sleep(1 * time.Second)
	}

	if err := writeSSE(w, "event: end", nil); err != nil {
		return err
	}
	flusher.Flush()

	return nil
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
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

func randomIdGenerator() string {
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	return fmt.Sprintf("chat.completion-%x", rand.Int63())
}

func writeSSE(w http.ResponseWriter, prefix string, v any) error {
	_, err := w.Write([]byte(prefix))
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(v)
}
