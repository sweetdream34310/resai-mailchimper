package chatgpt

type SummaryRequest struct {
	Model       string           `json:"model"`
	Messages    []SummaryMessage `json:"messages"`
	MaxTokens   int              `json:"max_tokens"`
	Temperature int              `json:"temperature"`
}

type SummaryMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SummaryResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}
