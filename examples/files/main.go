package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rizome-dev/go-moonshot"
)

func main() {
	// Create a new SDK instance
	sdk := moonshot.New()

	ctx := context.Background()

	// Example 1: Upload a file from disk
	fmt.Println("=== File Upload Example ===")
	
	// Create a temporary file for demonstration
	tempFile, err := os.CreateTemp("", "moonshot-example-*.txt")
	if err != nil {
		log.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write some content to the file
	content := `This is a sample document for the Moonshot AI API.
It contains information about various topics:

1. Introduction to AI
2. Machine Learning basics
3. Natural Language Processing

This file will be used to demonstrate file operations with the Moonshot API.`

	if _, err := tempFile.WriteString(content); err != nil {
		log.Fatalf("Error writing to temp file: %v", err)
	}
	tempFile.Close()

	// Upload the file
	file, err := sdk.Files.UploadFile(ctx, tempFile.Name(), "assistants")
	if err != nil {
		log.Fatalf("Error uploading file: %v", err)
	}

	fmt.Printf("File uploaded successfully!\n")
	fmt.Printf("- ID: %s\n", file.ID)
	fmt.Printf("- Filename: %s\n", file.Filename)
	fmt.Printf("- Size: %d bytes\n", file.Bytes)
	fmt.Printf("- Purpose: %s\n", file.Purpose)
	fmt.Printf("- Created: %d\n", file.CreatedAt)

	// Example 2: Upload file from memory
	fmt.Println("\n=== Upload from Memory ===")
	
	jsonContent := `{
  "settings": {
    "model": "moonshot-v1-8k",
    "temperature": 0.7,
    "max_tokens": 1000
  },
  "prompts": [
    "Explain quantum computing",
    "Write a poem about the moon"
  ]
}`

	reader := strings.NewReader(jsonContent)
	file2, err := sdk.Files.Upload(ctx, reader, "config.json", "assistants")
	if err != nil {
		log.Fatalf("Error uploading from memory: %v", err)
	}

	fmt.Printf("JSON file uploaded: %s\n", file2.ID)

	// Example 3: List all files
	fmt.Println("\n=== List Files ===")
	
	fileList, err := sdk.Files.List(ctx, nil)
	if err != nil {
		log.Fatalf("Error listing files: %v", err)
	}

	fmt.Printf("Total files: %d\n", len(fileList.Data))
	for i, f := range fileList.Data {
		fmt.Printf("%d. %s (ID: %s, Purpose: %s)\n", i+1, f.Filename, f.ID, f.Purpose)
		if i >= 4 {
			fmt.Println("   ... (showing first 5 files)")
			break
		}
	}

	// Example 4: Get file details
	fmt.Println("\n=== Get File Details ===")
	
	fileDetails, err := sdk.Files.Get(ctx, file.ID)
	if err != nil {
		log.Fatalf("Error getting file details: %v", err)
	}

	fmt.Printf("Retrieved file: %s\n", fileDetails.Filename)
	if fileDetails.Status != "" {
		fmt.Printf("Status: %s\n", fileDetails.Status)
	}

	// Example 5: Get file content
	fmt.Println("\n=== Get File Content ===")
	
	content2, err := sdk.Files.GetContent(ctx, file.ID)
	if err != nil {
		log.Fatalf("Error getting file content: %v", err)
	}

	fmt.Printf("File content (%d bytes):\n", len(content2))
	if len(content2) > 200 {
		fmt.Printf("%s... (truncated)\n", string(content2[:200]))
	} else {
		fmt.Println(string(content2))
	}

	// Example 6: Use file in chat completion
	fmt.Println("\n=== Using File in Chat ===")
	
	chatReq := moonshot.ChatCompletionRequest{
		Model: moonshot.ModelMoonshotV18K.String(),
		Messages: []moonshot.Message{
			{
				Role:    "user",
				Content: "Please summarize the main topics in the uploaded document.",
			},
		},
		FileIDs: []string{file.ID}, // Reference the uploaded file
	}

	resp, err := sdk.Chat.CreateCompletion(ctx, chatReq)
	if err != nil {
		log.Fatalf("Error in chat completion: %v", err)
	}

	if len(resp.Choices) > 0 {
		fmt.Println("Summary:", resp.Choices[0].Message.Content)
	}

	// Example 7: Delete file
	fmt.Println("\n=== Delete File ===")
	
	err = sdk.Files.Delete(ctx, file.ID)
	if err != nil {
		log.Fatalf("Error deleting file: %v", err)
	}
	fmt.Printf("File %s deleted successfully\n", file.ID)

	// Also delete the second file
	err = sdk.Files.Delete(ctx, file2.ID)
	if err != nil {
		log.Printf("Error deleting second file: %v", err)
	}

	// Example 8: Filter files by purpose
	fmt.Println("\n=== Filter Files by Purpose ===")
	
	assistantFiles, err := sdk.Files.List(ctx, &moonshot.FileListParams{
		Purpose: "assistants",
	})
	if err != nil {
		log.Fatalf("Error listing assistant files: %v", err)
	}

	fmt.Printf("Files with purpose 'assistants': %d\n", len(assistantFiles.Data))
}