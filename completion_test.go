package zhipuai_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bbang94/go-zhipuai"
	"github.com/bbang94/go-zhipuai/internal/test/checks"
)

func TestCompletionsWrongModel(t *testing.T) {
	config := zhipuai.DefaultConfig("whatever")
	config.BaseURL = "http://localhost/v1"
	client := zhipuai.NewClientWithConfig(config)

	_, err := client.CreateCompletion(
		context.Background(),
		zhipuai.CompletionRequest{
			MaxTokens: 5,
			Model:     zhipuai.GPT3Dot5Turbo,
		},
	)
	if !errors.Is(err, zhipuai.ErrCompletionUnsupportedModel) {
		t.Fatalf("CreateCompletion should return ErrCompletionUnsupportedModel, but returned: %v", err)
	}
}

func TestCompletionWithStream(t *testing.T) {
	config := zhipuai.DefaultConfig("whatever")
	client := zhipuai.NewClientWithConfig(config)

	ctx := context.Background()
	req := zhipuai.CompletionRequest{Stream: true}
	_, err := client.CreateCompletion(ctx, req)
	if !errors.Is(err, zhipuai.ErrCompletionStreamNotSupported) {
		t.Fatalf("CreateCompletion didn't return ErrCompletionStreamNotSupported")
	}
}

// TestCompletions Tests the completions endpoint of the API using the mocked server.
func TestCompletions(t *testing.T) {
	client, server, teardown := setupzhipuaiTestServer()
	defer teardown()
	server.RegisterHandler("/v1/completions", handleCompletionEndpoint)
	req := zhipuai.CompletionRequest{
		MaxTokens: 5,
		Model:     "ada",
		Prompt:    "Lorem ipsum",
	}
	_, err := client.CreateCompletion(context.Background(), req)
	checks.NoError(t, err, "CreateCompletion error")
}

// handleCompletionEndpoint Handles the completion endpoint by the test server.
func handleCompletionEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error
	var resBytes []byte

	// completions only accepts POST requests
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	var completionReq zhipuai.CompletionRequest
	if completionReq, err = getCompletionBody(r); err != nil {
		http.Error(w, "could not read request", http.StatusInternalServerError)
		return
	}
	res := zhipuai.CompletionResponse{
		ID:      strconv.Itoa(int(time.Now().Unix())),
		Object:  "test-object",
		Created: time.Now().Unix(),
		// would be nice to validate Model during testing, but
		// this may not be possible with how much upkeep
		// would be required / wouldn't make much sense
		Model: completionReq.Model,
	}
	// create completions
	n := completionReq.N
	if n == 0 {
		n = 1
	}
	for i := 0; i < n; i++ {
		// generate a random string of length completionReq.Length
		completionStr := strings.Repeat("a", completionReq.MaxTokens)
		if completionReq.Echo {
			completionStr = completionReq.Prompt.(string) + completionStr
		}
		res.Choices = append(res.Choices, zhipuai.CompletionChoice{
			Text:  completionStr,
			Index: i,
		})
	}
	inputTokens := numTokens(completionReq.Prompt.(string)) * n
	completionTokens := completionReq.MaxTokens * n
	res.Usage = zhipuai.Usage{
		PromptTokens:     inputTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      inputTokens + completionTokens,
	}
	resBytes, _ = json.Marshal(res)
	fmt.Fprintln(w, string(resBytes))
}

// getCompletionBody Returns the body of the request to create a completion.
func getCompletionBody(r *http.Request) (zhipuai.CompletionRequest, error) {
	completion := zhipuai.CompletionRequest{}
	// read the request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return zhipuai.CompletionRequest{}, err
	}
	err = json.Unmarshal(reqBody, &completion)
	if err != nil {
		return zhipuai.CompletionRequest{}, err
	}
	return completion, nil
}
