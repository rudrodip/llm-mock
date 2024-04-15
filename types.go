package main

type Request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
	Streaming   *bool     `json:"streaming"`
}

type Response struct {
	Id      string
	Object  string
	Created uint64
	Model   string
	Usage   Usage
	Choices []Choice
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	Prompt_tokens     int `json:"prompt_tokens"`
	Completion_tokens int `json:"completion_tokens"`
	Total_tokens      int `json:"total_tokens"`
}

type Choice struct {
	Message       Message     `json:"message"`
	Logprobs      interface{} `json:"logprobs"`
	Finish_reason string      `json:"finish_reason"`
	Index         uint        `json:"index"`
}
