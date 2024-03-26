package main

import (
	"context"
	"fmt"
	"github.com/bbang94/go-zhipuai"
)

func main() {
	client := zhipuai.NewClient("")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		zhipuai.ChatCompletionRequest{
			Model:            "glm-4",
			Messages:         []zhipuai.ChatCompletionMessage{{Role: "user", Content: "你好，你谁谁"}},
			Temperature:      0.01,
			TopP:             0,
			Stream:           false,
			Stop:             nil,
			PresencePenalty:  0,
			ResponseFormat:   nil,
			Seed:             nil,
			FrequencyPenalty: 0,
			LogitBias:        nil,
			LogProbs:         false,
			TopLogProbs:      0,
			User:             "",
			ToolChoice:       nil,
		},
	)
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
