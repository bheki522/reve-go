# reve-go

[![Go Reference](https://pkg.go.dev/badge/github.com/shamspias/reve-go.svg)](https://pkg.go.dev/github.com/shamspias/reve-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/shamspias/reve-go)](https://goreportcard.com/report/github.com/shamspias/reve-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

An unofficial Go SDK for the [Reve Image Generation API](https://reve.com) - a powerful AI platform known for stunning aesthetics, accurate text rendering, and natural-language image edits.

## Features

- 🎨 **Create** - Generate images from text descriptions
- ✏️ **Edit** - Modify images with text instructions
- 🔄 **Remix** - Combine multiple images with prompts
- 🌐 **Proxy Support** - HTTP, HTTPS, SOCKS5 proxies
- 🔁 **Auto Retry** - Exponential backoff with jitter
- 📦 **Batch Processing** - Concurrent operations
- 💰 **Cost Estimation** - Estimate before you spend

## Installation

```bash
go get github.com/shamspias/reve-go
```

**Requires Go 1.25+**

## Quick Start

```go
package main

import (
    "context"
    "log"
    "os"

    reve "github.com/shamspias/reve-go"
)

func main() {
    client := reve.NewClient(os.Getenv("REVE_API_KEY"))

    result, err := client.Images.Create(context.Background(), &reve.CreateParams{
        Prompt: "A beautiful mountain landscape at sunset",
    })
    if err != nil {
        log.Fatal(err)
    }

    result.SaveTo("landscape.png")
    log.Printf("Saved! Credits used: %d", result.CreditsUsed)
}
```

## Project Structure

```
reve-go/
├── client.go           # Main client
├── options.go          # Configuration options
├── reve.go             # Re-exports for convenience
├── doc.go              # Package documentation
│
├── types/              # Shared types
│   ├── aspect.go       # AspectRatio type
│   ├── version.go      # ModelVersion type
│   ├── format.go       # OutputFormat type
│   ├── postprocess.go  # Postprocess operations
│   ├── image.go        # Image helper type
│   └── result.go       # Result types
│
├── image/              # Image service
│   ├── service.go      # Service definition
│   ├── create.go       # Create operation
│   ├── edit.go         # Edit operation
│   ├── remix.go        # Remix operation
│   ├── batch.go        # Batch operations
│   └── cost.go         # Cost estimation
│
├── internal/           # Internal packages
│   ├── transport/      # HTTP transport
│   │   ├── client.go
│   │   ├── retry.go
│   │   └── errors.go
│   └── validator/      # Validation
│       └── validator.go
│
└── examples/           # Example applications
    ├── basic/
    ├── create/
    ├── edit/
    ├── remix/
    ├── batch/
    ├── proxy/
    ├── error-handling/
    └── complete/
```

## Documentation

### Client Configuration

```go
// Basic
client := reve.NewClient(apiKey)

// With options
client := reve.NewClient(apiKey,
    reve.WithTimeout(60*time.Second),
    reve.WithRetry(5, time.Second, 30*time.Second),
    reve.WithDebug(true),
)
```

### Proxy Support

```go
// HTTP Proxy
client := reve.NewClient(apiKey,
    reve.WithHTTPProxy("http://proxy:8080"),
)

// SOCKS5 Proxy
client := reve.NewClient(apiKey,
    reve.WithSOCKS5Proxy("127.0.0.1:1080", "user", "pass"),
)

// Environment variables (HTTP_PROXY, HTTPS_PROXY)
client := reve.NewClient(apiKey,
    reve.WithProxyFromEnvironment(),
)
```

### Create Images

```go
result, err := client.Images.Create(ctx, &reve.CreateParams{
    Prompt:          "A cyberpunk cityscape",
    AspectRatio:     reve.Ratio16x9,
    TestTimeScaling: 2,
    Postprocess:     []reve.Postprocess{reve.Upscale(2)},
})
```

### Edit Images

```go
img, _ := reve.NewImageFromFile("photo.jpg")

result, err := client.Images.Edit(ctx, &reve.EditParams{
    Instruction:    "Convert to watercolor painting",
    ReferenceImage: img.Base64(),
    Version:        reve.VersionLatestFast, // 5 credits
})
```

### Remix Images

```go
style, _ := reve.NewImageFromFile("style.png")
content, _ := reve.NewImageFromFile("content.png")

result, err := client.Images.Remix(ctx, &reve.RemixParams{
    Prompt: fmt.Sprintf("Apply %s to %s", reve.Ref(0), reve.Ref(1)),
    ReferenceImages: []string{style.Base64(), content.Base64()},
})
```

### Batch Operations

```go
requests := []*reve.CreateParams{
    {Prompt: "A red apple"},
    {Prompt: "A green pear"},
}

results := client.Images.BatchCreate(ctx, requests, &reve.BatchConfig{
    Concurrency: 3,
})

fmt.Printf("Success: %d/%d\n", reve.SuccessCount(results), len(results))
```

### Error Handling

```go
result, err := client.Images.Create(ctx, params)
if err != nil {
    var apiErr *transport.APIError
    if errors.As(err, &apiErr) {
        if apiErr.IsRateLimit() {
            // Wait and retry
        }
        if apiErr.IsInsufficientFunds() {
            // Need more credits
        }
    }
}
```

### Cost Estimation

```go
cost := reve.EstimateCreate(1, nil)
fmt.Println(cost) // "18 credits (~$0.0240)"

cost = reve.EstimateEdit(true, 1, nil) // Fast mode
fmt.Println(cost) // "5 credits (~$0.0067)"
```

## Examples

Run examples with:

```bash
REVE_API_KEY=your-key go run examples/basic/main.go
REVE_API_KEY=your-key go run examples/create/main.go
REVE_API_KEY=your-key go run examples/edit/main.go
REVE_API_KEY=your-key go run examples/remix/main.go
REVE_API_KEY=your-key go run examples/batch/main.go
REVE_API_KEY=your-key go run examples/proxy/main.go
REVE_API_KEY=your-key go run examples/error-handling/main.go
REVE_API_KEY=your-key go run examples/complete/main.go
```

## Pricing

| Endpoint | Credits | ~USD |
|----------|---------|------|
| Create | 18 | $0.024 |
| Edit | 30 | $0.040 |
| Edit Fast | 5 | $0.007 |
| Remix | 30 | $0.040 |
| Remix Fast | 5 | $0.007 |

## Contributing

1. Fork the repo
2. Create feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add feature'`)
4. Push (`git push origin feature/amazing`)
5. Open Pull Request

## License

MIT License - see [LICENSE](LICENSE)

## Disclaimer

This is an unofficial SDK, not affiliated with Reve.
