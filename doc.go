// Package reve provides an unofficial Go SDK for the Reve Image Generation API.
//
// Reve is a powerful AI image generation platform known for stunning aesthetics,
// accurate text rendering, prompt adherence, and natural-language image edits.
//
// # Installation
//
//	go get github.com/shamspias/reve-go
//
// # Quick Start
//
//	package main
//
//	import (
//		"context"
//		"log"
//		"os"
//
//		reve "github.com/shamspias/reve-go"
//	)
//
//	func main() {
//		client := reve.NewClient(os.Getenv("REVE_API_KEY"))
//
//		result, err := client.Images.Create(context.Background(), &reve.CreateParams{
//			Prompt: "A beautiful mountain landscape at sunset",
//		})
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		if err := result.SaveTo("landscape.png"); err != nil {
//			log.Fatal(err)
//		}
//	}
//
// # Architecture
//
// The SDK uses a service-based architecture:
//
//	client := reve.NewClient(apiKey)
//	client.Images.Create(...)   // Create images
//	client.Images.Edit(...)     // Edit images
//	client.Images.Remix(...)    // Remix images
//
// # Proxy Support
//
//	// HTTP/HTTPS proxy
//	client := reve.NewClient(apiKey,
//		reve.WithHTTPProxy("http://proxy.example.com:8080"),
//	)
//
//	// SOCKS5 proxy
//	client := reve.NewClient(apiKey,
//		reve.WithSOCKS5Proxy("127.0.0.1:1080", "user", "pass"),
//	)
//
// For more examples, see the examples directory.
package reve

// Version is the SDK version.
const Version = "1.0.0"
