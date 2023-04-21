package chatgpt

import (
	"encoding/json"
	"fmt"
	"github.com/cloudsrc/api.awaymail.v1.go/config"
	"github.com/cloudsrc/api.awaymail.v1.go/src/shared/restclient"
	"net/http"
	"os"
	"time"
)

// wrapper ...
type wrapper struct {
	cfg    config.Config
	client restclient.Client
}

// Setup ...
func Setup(config config.Config) Wrapper {
	return &wrapper{
		cfg: config,
		client: restclient.New(restclient.Options{
			Address: config.Chatgpt.ChatGPTURL,
			Timeout: 10 * time.Second,
			SkipTLS: config.Chatgpt.ChatGPTSkipTLS,
		}),
	}
}
func (w *wrapper) CreateSummary(content string) (resp *SummaryResponse, err error) {
	headers := http.Header{}
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", w.cfg.Chatgpt.ChatGPTToken))

	req := &SummaryRequest{
		Model: "gpt-3.5-turbo",
		Messages: []SummaryMessage{
			{Role: "assistant", Content: content},
		},
		MaxTokens:   14,
		Temperature: 0,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return
	}

	body, statusCode, httpErr := w.client.Post("/v1/chat/completions", headers, body)

	if os.IsTimeout(httpErr) {
		err = httpErr
		return
	}

	if statusCode != http.StatusOK {
		err = fmt.Errorf("%v", string(body[:]))
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return
	}

	return
}
