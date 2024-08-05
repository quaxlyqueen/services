package services

type Communication struct {
	Communication string `json:"communication"`
	Hash          string `json:"hash"`
}

type Response struct {
	Model         string `json:"model"`
	CreatedAt     string `json:"created_at"`
	Response      string `json:"response"`
	Done          bool   `json:"done"`
	DoneReason    string `json:"done_reason"`
	Context       []int  `json:"context"`
	TotalDuration int    `json:"total_duration"`
	LoadDuration  int    `json:"load_duration"`
	PromptEC      int    `json:"prompt_eval_count"`
	PromptED      int    `json:"prompt_eval_duration"`
	EvalCount     int    `json:"eval_count"`
	EvalDuration  int    `json:"eval_duration"`
}

// Chats are updated once the content is decrypted. Since this is never leaving
// the server once decrypted, and possibly will not be saved (depending on user
// settings), it will remain decrypted until it's time to transform into a
// Communication JSON object, which is the HTTP response that is re-encrypted.

// TODO: If the user has enabled conversation history, then save the encrypted
// chats to the server's drive. Additionally, only accept incoming additional
// messages from the user, rather than having the client re-send Chats already
// stored on the server.
type Chat struct {
	Role    bool   `json:"role"`
	Content string `json:"content"`
}

type PromptWHistory struct {
	Model    string `json:"model"`
	Messages []Chat `json:"messages"`
	Stream   bool   `json:"stream"`
}

type ResponseWHistory struct {
	Model         string `json:"model"`
	CreatedAt     string `json:"created_at"`
	Message       Chat   `json:"message"`
	Done          bool   `json:"done"`
	DoneReason    string `json:"done_reason"`
	Context       []int  `json:"context"`
	TotalDuration int    `json:"total_duration"`
	LoadDuration  int    `json:"load_duration"`
	PromptEC      int    `json:"prompt_eval_count"`
	PromptED      int    `json:"prompt_eval_duration"`
	EvalCount     int    `json:"eval_count"`
	EvalDuration  int    `json:"eval_duration"`
}
