//go:build mcp_go_client_oauth

// CarbonArc MCP Client — Go Quick Start
// ======================================
//
// Connect to the CarbonArc MCP server using the official
// modelcontextprotocol/go-sdk, authenticate via OAuth 2.0,
// discover available tools, and call the search_entities tool.
//
// Setup:
//
//	go mod tidy
//
// Run:
//
//	go run -tags mcp_go_client_oauth .
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	mcpServerURL = "https://mcp.carbonarc.ai/"
	resourceURL  = "https://mcp.carbonarc.ai"
	callbackPort = 8092
)

var callbackURL = fmt.Sprintf("http://localhost:%d/callback", callbackPort)

func main() {
	ctx := context.Background()

	// Connect with browser-based OAuth (see oauthhelper.go)
	httpClient, err := ConnectWithOAuth(resourceURL, callbackURL, callbackPort)
	if err != nil {
		log.Fatal(err)
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "carbonarc-mcp-client-go",
		Version: "1.0.0",
	}, nil)

	log.Printf("Connecting to %s", mcpServerURL)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{
		Endpoint:   mcpServerURL,
		HTTPClient: httpClient,
	}, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer session.Close()

	log.Printf("Connected (session %s)", session.ID())

	// Discover tools
	toolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	fmt.Printf("\nDiscovered %d tool(s):\n\n", len(toolsResult.Tools))
	for _, tool := range toolsResult.Tools {
		desc := tool.Description
		for i, c := range desc {
			if c == '\n' {
				desc = desc[:i]
				break
			}
		}
		fmt.Printf("  • %s — %s\n", tool.Name, desc)
	}

	// Test: search for "Walmart"
	fmt.Println("\nSearching for Walmart\n")
	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "search_entities",
		Arguments: map[string]any{
			"query": "Walmart",
		},
	})
	if err != nil {
		log.Fatalf("Tool call failed: %v", err)
	}

	for _, content := range result.Content {
		if tc, ok := content.(*mcp.TextContent); ok {
			var pretty json.RawMessage
			if json.Unmarshal([]byte(tc.Text), &pretty) == nil {
				formatted, _ := json.MarshalIndent(pretty, "", "  ")
				fmt.Println(string(formatted))
			} else {
				fmt.Println(tc.Text)
			}
		}
	}

	log.Println("Done")
}
