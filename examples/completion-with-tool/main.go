package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bbang94/go-zhipuai"
	"github.com/bbang94/go-zhipuai/jsonschema"
)

func main() {
	ctx := context.Background()
	client := zhipuai.NewClient(os.Getenv("zhipuai_API_KEY"))

	// describe the function & its inputs
	params := jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"location": {
				Type:        jsonschema.String,
				Description: "The city and state, e.g. San Francisco, CA",
			},
			"unit": {
				Type: jsonschema.String,
				Enum: []string{"celsius", "fahrenheit"},
			},
		},
		Required: []string{"location"},
	}
	f := zhipuai.FunctionDefinition{
		Name:        "get_current_weather",
		Description: "Get the current weather in a given location",
		Parameters:  params,
	}
	t := zhipuai.Tool{
		Type:     zhipuai.ToolTypeFunction,
		Function: &f,
	}

	// simulate user asking a question that requires the function
	dialogue := []zhipuai.ChatCompletionMessage{
		{Role: zhipuai.ChatMessageRoleUser, Content: "What is the weather in Boston today?"},
	}
	fmt.Printf("Asking zhipuai '%v' and providing it a '%v()' function...\n",
		dialogue[0].Content, f.Name)
	resp, err := client.CreateChatCompletion(ctx,
		zhipuai.ChatCompletionRequest{
			Model:    zhipuai.GPT4TurboPreview,
			Messages: dialogue,
			Tools:    []zhipuai.Tool{t},
		},
	)
	if err != nil || len(resp.Choices) != 1 {
		fmt.Printf("Completion error: err:%v len(choices):%v\n", err,
			len(resp.Choices))
		return
	}
	msg := resp.Choices[0].Message
	if len(msg.ToolCalls) != 1 {
		fmt.Printf("Completion error: len(toolcalls): %v\n", len(msg.ToolCalls))
		return
	}

	// simulate calling the function & responding to zhipuai
	dialogue = append(dialogue, msg)
	fmt.Printf("zhipuai called us back wanting to invoke our function '%v' with params '%v'\n",
		msg.ToolCalls[0].Function.Name, msg.ToolCalls[0].Function.Arguments)
	dialogue = append(dialogue, zhipuai.ChatCompletionMessage{
		Role:       zhipuai.ChatMessageRoleTool,
		Content:    "Sunny and 80 degrees.",
		Name:       msg.ToolCalls[0].Function.Name,
		ToolCallID: msg.ToolCalls[0].ID,
	})
	fmt.Printf("Sending zhipuai our '%v()' function's response and requesting the reply to the original question...\n",
		f.Name)
	resp, err = client.CreateChatCompletion(ctx,
		zhipuai.ChatCompletionRequest{
			Model:    zhipuai.GPT4TurboPreview,
			Messages: dialogue,
			Tools:    []zhipuai.Tool{t},
		},
	)
	if err != nil || len(resp.Choices) != 1 {
		fmt.Printf("2nd completion error: err:%v len(choices):%v\n", err,
			len(resp.Choices))
		return
	}

	// display zhipuai's response to the original question utilizing our function
	msg = resp.Choices[0].Message
	fmt.Printf("zhipuai answered the original request with: %v\n",
		msg.Content)
}
