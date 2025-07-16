# Moonshot AI SDK Examples

This directory contains examples demonstrating how to use the Moonshot AI Go SDK.

## Examples

### Basic Chat Completion
[basic/main.go](basic/main.go) - Simple chat completion example showing:
- Creating a client
- Making chat completion requests
- Using different models (moonshot-v1-8k, moonshot-v1-32k)
- Handling responses and errors

### Streaming Responses
[streaming/main.go](streaming/main.go) - Streaming chat completions showing:
- Callback-based streaming
- Reader-based streaming for more control
- Context cancellation
- Processing streaming chunks

### File Operations
[files/main.go](files/main.go) - File management examples showing:
- Uploading files from disk
- Uploading files from memory
- Listing files
- Getting file details and content
- Using files in chat completions
- Deleting files

## Running the Examples

1. Set your API key:
```bash
export MOONSHOT_API_KEY="sk-your-api-key"
```

2. Run an example:
```bash
# Basic example
go run examples/basic/main.go

# Streaming example
go run examples/streaming/main.go

# File operations example
go run examples/files/main.go
```

## API Key

You can provide your API key in two ways:

1. Environment variable (recommended):
```bash
export MOONSHOT_API_KEY="sk-your-api-key"
```

2. Directly in code:
```go
sdk := moonshot.New("sk-your-api-key")
```

## Models

The SDK supports various Moonshot models:

- `moonshot-v1-8k` - 8K context window
- `moonshot-v1-32k` - 32K context window
- `moonshot-v1-128k` - 128K context window
- `kimi-k2` - Latest Kimi K2 model (may require special access)

## Error Handling

The SDK provides typed errors for better error handling:

```go
resp, err := sdk.Chat.CreateCompletion(ctx, req)
if err != nil {
    // Check if it's an API error
    if apiErr, ok := moonshot.IsAPIError(err); ok {
        log.Printf("API error: %s (code: %s)", apiErr.Message, apiErr.Code)
    } else {
        log.Printf("Other error: %v", err)
    }
}
```

## Temperature Adjustment

Note: The Moonshot API automatically adjusts temperature values. The SDK handles this for you:
- Your temperature × 0.6 = actual temperature used by the API
- For example, setting temperature to 1.0 will use 0.6 in the API