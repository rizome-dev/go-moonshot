package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rizome-dev/go-moonshot"
)

func main() {
	// Create a new SDK instance
	// Uses MOONSHOT_API_KEY environment variable by default
	sdk := moonshot.New()

	// Or provide API key directly
	// sdk := moonshot.New("sk-your-api-key")

	// Create a chat completion request
	req := moonshot.ChatCompletionRequest{
		Model: moonshot.ModelKimiK2Instruct.String(),
		Messages: []moonshot.Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role:    "user",
				Content: "Hello! Can you explain what the Moonshot API is?",
			},
		},
		Temperature: moonshot.Float64(0.7),
		MaxTokens:   moonshot.Int(500),
	}

	// Make the API call
	resp, err := sdk.Chat.CreateCompletion(context.Background(), req)
	if err != nil {
		// Check if it's an API error
		if apiErr, ok := moonshot.IsAPIError(err); ok {
			log.Fatalf("API error: %s (code: %s)", apiErr.Message, apiErr.Code)
		}
		log.Fatalf("Error: %v", err)
	}

	// Print the response
	if len(resp.Choices) > 0 {
		fmt.Println("Assistant:", resp.Choices[0].Message.Content)
		fmt.Printf("\nTokens used: %d (prompt: %d, completion: %d)\n",
			resp.Usage.TotalTokens,
			resp.Usage.PromptTokens,
			resp.Usage.CompletionTokens)
	}

	// Example with a longer context model
	fmt.Println("\n--- Using a longer context model ---")

	req.Model = moonshot.ModelMoonshotV132K.String()
	req.Messages = append(req.Messages, moonshot.Message{
		Role:    "assistant",
		Content: resp.Choices[0].Message.Content.(string),
	})
	req.Messages = append(req.Messages, moonshot.Message{
		Role:    "user",
		Content: "Can you give me a code example of how to use it?",
	})

	resp2, err := sdk.Chat.CreateCompletion(context.Background(), req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if len(resp2.Choices) > 0 {
		fmt.Println("Assistant:", resp2.Choices[0].Message.Content)
	}

	// Example using the latest Kimi K2 model
	if os.Getenv("USE_KIMI_K2") == "true" {
		fmt.Println("\n--- Using Kimi K2 model ---")

		req.Model = moonshot.ModelKimiK2.String()
		req.Messages = []moonshot.Message{
			{
				Role:    "user",
				Content: "What makes Kimi K2 special compared to other models?",
			},
		}

		resp3, err := sdk.Chat.CreateCompletion(context.Background(), req)
		if err != nil {
			log.Printf("Kimi K2 error (this model may require special access): %v", err)
		} else if len(resp3.Choices) > 0 {
			fmt.Println("Kimi K2:", resp3.Choices[0].Message.Content)
		}
	}
}
