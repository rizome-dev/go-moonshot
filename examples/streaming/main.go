package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/rizome-dev/go-moonshot"
)

func main() {
	// Create a new SDK instance
	sdk := moonshot.New()

	// Create a chat completion request
	req := moonshot.ChatCompletionRequest{
		Model: moonshot.ModelMoonshotV18K.String(),
		Messages: []moonshot.Message{
			{
				Role:    "system",
				Content: "You are a helpful storyteller.",
			},
			{
				Role:    "user",
				Content: "Tell me a short story about a robot learning to paint.",
			},
		},
		Temperature: moonshot.Float64(0.8),
		MaxTokens:   moonshot.Int(1000),
	}

	fmt.Println("=== Streaming Response ===")
	fmt.Print("Assistant: ")

	// Option 1: Use callback-based streaming
	err := sdk.Chat.CreateCompletionWithCallback(context.Background(), req, func(chunk *moonshot.ChatCompletionStream) error {
		// Print each chunk as it arrives
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != nil {
			fmt.Print(*chunk.Choices[0].Delta.Content)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Streaming error: %v", err)
	}
	fmt.Println()

	// Option 2: Use reader-based streaming for more control
	fmt.Println("=== Reader-based Streaming ===")
	
	req.Messages = []moonshot.Message{
		{
			Role:    "user",
			Content: "Now tell me a haiku about that robot.",
		},
	}

	stream, err := sdk.Chat.CreateCompletionStream(context.Background(), req)
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}
	defer stream.Close()

	fmt.Print("Assistant: ")
	var fullResponse string
	
	for {
		chunk, err := stream.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error reading stream: %v", err)
		}

		// Process each chunk
		if len(chunk.Choices) > 0 {
			choice := chunk.Choices[0]
			
			// Check if this is the initial chunk with role
			if choice.Delta.Role != nil {
				// Initial chunk, role is set
				continue
			}
			
			// Print content
			if choice.Delta.Content != nil {
				content := *choice.Delta.Content
				fmt.Print(content)
				fullResponse += content
			}
			
			// Check if we're done
			if choice.FinishReason != nil && *choice.FinishReason == "stop" {
				fmt.Println("\n\nStream finished!")
			}
		}

		// Some chunks may include usage information at the end
		if chunk.Usage != nil {
			fmt.Printf("\nTokens used: %d (prompt: %d, completion: %d)\n",
				chunk.Usage.TotalTokens,
				chunk.Usage.PromptTokens,
				chunk.Usage.CompletionTokens)
		}
	}

	fmt.Printf("\nFull response collected: %d characters\n", len(fullResponse))

	// Example with context tracking
	fmt.Println("\n=== Streaming with Context ===")
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req.Messages = []moonshot.Message{
		{
			Role:    "system",
			Content: "You are a helpful assistant that provides step-by-step explanations.",
		},
		{
			Role:    "user",
			Content: "Explain how to make a simple HTTP server in Go, step by step.",
		},
	}

	charCount := 0
	err = sdk.Chat.CreateCompletionWithCallback(ctx, req, func(chunk *moonshot.ChatCompletionStream) error {
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != nil {
			content := *chunk.Choices[0].Delta.Content
			fmt.Print(content)
			charCount += len(content)
			
			// Example: Cancel streaming after 500 characters
			if charCount > 500 {
				fmt.Println("\n\n[Streaming cancelled after 500 characters]")
				cancel()
				return io.EOF
			}
		}
		return nil
	})

	if err != nil && err != io.EOF {
		log.Printf("Streaming error: %v", err)
	}
}