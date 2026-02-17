//go:build mcp_go_client_oauth

// Package main provides a reusable OAuth helper for MCP clients.
//
// BrowserOAuthFlow implements the browser-based authorization code flow
// (with PKCE) expected by auth.NewHTTPTransport. It discovers the
// authorization server metadata, opens the user's browser, spins up a
// local callback server, and exchanges the code for a token — so the
// main MCP client example can stay focused on the protocol.
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/oauthex"
	"golang.org/x/oauth2"
)

// BrowserOAuthFlow is an auth.OAuthHandler that performs the full
// browser-based authorization code flow.  Pass it to auth.NewHTTPTransport.
func BrowserOAuthFlow(resourceURL, callbackURL string, callbackPort int) func(*http.Request, *http.Response) (oauth2.TokenSource, error) {
	return func(req *http.Request, res *http.Response) (oauth2.TokenSource, error) {
		ctx := req.Context()

		// 1. Discover authorization server metadata
		asm, err := oauthex.GetAuthServerMeta(ctx, resourceURL, nil)
		if err != nil {
			return nil, fmt.Errorf("auth server metadata: %w", err)
		}
		log.Printf("Auth server: %s", asm.Issuer)

		// 2. Generate PKCE code verifier + challenge (S256)
		verifierBytes := make([]byte, 32)
		if _, err := rand.Read(verifierBytes); err != nil {
			return nil, err
		}
		codeVerifier := base64.RawURLEncoding.EncodeToString(verifierBytes)
		h := sha256.Sum256([]byte(codeVerifier))
		codeChallenge := base64.RawURLEncoding.EncodeToString(h[:])

		// 3. Build authorization URL
		authURL, err := url.Parse(asm.AuthorizationEndpoint)
		if err != nil {
			return nil, err
		}
		q := authURL.Query()
		q.Set("response_type", "code")
		q.Set("client_id", "carbonarc-mcp-fixed-client")
		q.Set("redirect_uri", callbackURL)
		q.Set("code_challenge", codeChallenge)
		q.Set("code_challenge_method", "S256")
		q.Set("resource", resourceURL+"/")
		authURL.RawQuery = q.Encode()

		// 4. Start local callback server and open browser
		codeCh := make(chan string, 1)
		mux := http.NewServeMux()
		mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code != "" {
				fmt.Fprint(w, "<h2>Authorized — you can close this tab.</h2>")
				codeCh <- code
			} else {
				http.Error(w, "missing code", http.StatusBadRequest)
			}
		})
		srv := &http.Server{Addr: fmt.Sprintf(":%d", callbackPort), Handler: mux}
		go srv.ListenAndServe()

		fmt.Printf("\nOpening browser for login: %s\n\n", authURL)
		openBrowser(authURL.String())

		// 5. Wait for the auth code
		code := <-codeCh
		srv.Shutdown(ctx)

		// 6. Exchange the code for a token
		oauthConf := &oauth2.Config{
			ClientID:    "carbonarc-mcp-fixed-client",
			Endpoint:    oauth2.Endpoint{TokenURL: asm.TokenEndpoint},
			RedirectURL: callbackURL,
		}
		token, err := oauthConf.Exchange(ctx, code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
		if err != nil {
			return nil, fmt.Errorf("token exchange: %w", err)
		}
		log.Printf("Authenticated successfully")

		return oauthConf.TokenSource(ctx, token), nil
	}
}

func openBrowser(rawURL string) {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "linux":
		cmd = "xdg-open"
	default:
		cmd = "start"
	}
	exec.Command(cmd, rawURL).Start()
}

// ConnectWithOAuth creates an http.Client wrapped with the OAuth handler
// and returns a ready-to-use *http.Client for MCP transports.
func ConnectWithOAuth(resourceURL, callbackURL string, callbackPort int) (*http.Client, error) {
	httpTransport, err := auth.NewHTTPTransport(
		BrowserOAuthFlow(resourceURL, callbackURL, callbackPort), nil,
	)
	if err != nil {
		return nil, err
	}
	return &http.Client{Transport: httpTransport}, nil
}
