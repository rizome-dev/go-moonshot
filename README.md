# go-moonshot

[![GoDoc](https://pkg.go.dev/badge/github.com/rizome-dev/go-moonshot)](https://pkg.go.dev/github.com/rizome-dev/go-moonshot)
[![Go Report Card](https://goreportcard.com/badge/github.com/rizome-dev/go-moonshot)](https://goreportcard.com/report/github.com/rizome-dev/go-moonshot)
[![CI](https://github.com/rizome-dev/go-moonshot/actions/workflows/ci.yml/badge.svg)](https://github.com/rizome-dev/go-moonshot/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A comprehensive Go SDK for the Moonshot AI API, providing easy access to chat completions, file operations, and more.

built by [rizome labs](https://rizome.dev) | contact: [hi@rizome.dev](mailto:hi@rizome.dev)

## Installation

```bash
go get github.com/rizome-dev/go-moonshot
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/rizome-dev/go-moonshot"
)

func main() {
    // Create a new SDK instance (uses MOONSHOT_API_KEY env var)
    sdk := moonshot.New()
    
    // Create a chat completion
    req := moonshot.ChatCompletionRequest{
        Model: moonshot.ModelMoonshotV18K.String(),
        Messages: []moonshot.Message{
            {
                Role:    "user",
                Content: "Hello, Moonshot!",
            },
        },
    }
    
    resp, err := sdk.Chat.CreateCompletion(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(resp.Choices[0].Message.Content)
}
```

## Configuration

### API Key

The SDK supports multiple ways to provide your API key:

```go
// 1. Environment variable (recommended)
os.Setenv("MOONSHOT_API_KEY", "sk-your-api-key")
sdk := moonshot.New()

// 2. Direct initialization
sdk := moonshot.New("sk-your-api-key")

// 3. With additional options
sdk := moonshot.New("sk-your-api-key", 
    client.WithTimeout(60 * time.Second),
    client.WithBaseURL("https://custom.api.url"),
)
```

## Available Models

```go
// Standard Moonshot models
moonshot.ModelMoonshotV18K   // 8K context window
moonshot.ModelMoonshotV132K  // 32K context window
moonshot.ModelMoonshotV1128K // 128K context window

// Kimi K2 models
moonshot.ModelKimiK2         // Latest Kimi K2
moonshot.ModelKimiK2Base     // Base model for fine-tuning
moonshot.ModelKimiK2Instruct // Instruction-tuned model
```

## Chat Completions

### Basic Usage

```go
req := moonshot.ChatCompletionRequest{
    Model: moonshot.ModelMoonshotV18K.String(),
    Messages: []moonshot.Message{
        {Role: "system", Content: "You are a helpful assistant."},
        {Role: "user", Content: "What is the weather like?"},
    },
    Temperature: moonshot.Float64(0.7),
    MaxTokens:   moonshot.Int(1000),
}

resp, err := sdk.Chat.CreateCompletion(ctx, req)
```

### Streaming Responses

```go
// Callback-based streaming
err := sdk.Chat.CreateCompletionWithCallback(ctx, req, func(chunk *moonshot.ChatCompletionStream) error {
    if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != nil {
        fmt.Print(*chunk.Choices[0].Delta.Content)
    }
    return nil
})

// Reader-based streaming for more control
stream, err := sdk.Chat.CreateCompletionStream(ctx, req)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for {
    chunk, err := stream.Read()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    // Process chunk...
}
```

## File Operations

### Upload Files

```go
// Upload from file path
file, err := sdk.Files.UploadFile(ctx, "/path/to/file.txt", "assistants")

// Upload from reader
reader := strings.NewReader("file content")
file, err := sdk.Files.Upload(ctx, reader, "filename.txt", "assistants")
```

### List and Manage Files

```go
// List all files
files, err := sdk.Files.List(ctx, nil)

// Filter by purpose
files, err := sdk.Files.List(ctx, &moonshot.FileListParams{
    Purpose: "assistants",
})

// Get file details
file, err := sdk.Files.Get(ctx, "file-id")

// Get file content
content, err := sdk.Files.GetContent(ctx, "file-id")

// Delete file
err := sdk.Files.Delete(ctx, "file-id")
```

### Use Files in Chat

```go
req := moonshot.ChatCompletionRequest{
    Model: moonshot.ModelMoonshotV18K.String(),
    Messages: []moonshot.Message{
        {Role: "user", Content: "Summarize the uploaded document"},
    },
    FileIDs: []string{"file-id-1", "file-id-2"},
}
```

## Error Handling

The SDK provides typed errors for better error handling:

```go
resp, err := sdk.Chat.CreateCompletion(ctx, req)
if err != nil {
    // Check if it's an API error
    if apiErr, ok := moonshot.IsAPIError(err); ok {
        fmt.Printf("API Error: %s (Code: %s)\n", apiErr.Message, apiErr.Code)
        
        // Handle specific error codes
        switch apiErr.Code {
        case moonshot.ErrCodeRateLimitExceeded:
            // Handle rate limit
        case moonshot.ErrCodeInvalidRequest:
            // Handle invalid request
        }
    }
}
```

## Advanced Features

### Tool Use / Function Calling

```go
req := moonshot.ChatCompletionRequest{
    Model: moonshot.ModelMoonshotV18K.String(),
    Messages: messages,
    Tools: []moonshot.Tool{
        {
            Type: "function",
            Function: moonshot.Function{
                Name:        "get_weather",
                Description: "Get current weather",
                Parameters: map[string]interface{}{
                    "type": "object",
                    "properties": map[string]interface{}{
                        "location": map[string]string{
                            "type":        "string",
                            "description": "City name",
                        },
                    },
                    "required": []string{"location"},
                },
            },
        },
    },
}
```

### Temperature Note

The Moonshot API automatically adjusts temperature values:
- Actual temperature = Your temperature × 0.6
- The SDK handles this automatically for you

## Examples

See the [examples](examples/) directory for complete working examples:

- [Basic Chat Completion](examples/basic/main.go)
- [Streaming Responses](examples/streaming/main.go)
- [File Operations](examples/files/main.go)

## Testing

Run all tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
