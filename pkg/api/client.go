package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/exmplrai/aphelion-cli/internal/config"
	"github.com/exmplrai/aphelion-cli/internal/logger"
	"github.com/exmplrai/aphelion-cli/pkg/auth"
)

// Client represents the Aphelion API client
type Client struct {
	httpClient  *http.Client
	baseURL     string
	auth0Client *auth.Auth0Client
}

// NewClient creates a new API client
func NewClient() (*Client, error) {
	profile := config.GetCurrentProfile()
	if profile.Endpoint == "" {
		return nil, fmt.Errorf("no endpoint configured")
	}
	
	client := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: profile.Endpoint,
	}
	
	// Initialize Auth0 client - discover config if not set
	if err := client.ensureAuth0Config(); err != nil {
		return nil, fmt.Errorf("failed to configure Auth0: %w", err)
	}
	
	return client, nil
}

// ensureAuth0Config ensures Auth0 configuration is available
func (c *Client) ensureAuth0Config() error {
	profile := config.GetCurrentProfile()
	
	// If Auth0 is already configured, use existing config
	if profile.Auth.Domain != "" && profile.Auth.ClientID != "" {
		c.auth0Client = auth.NewAuth0Client(
			profile.Auth.Domain,
			profile.Auth.ClientID,
			profile.Auth.Audience,
			profile.Auth.RedirectURI,
		)
		return nil
	}
	
	// Auto-discover Auth0 configuration from the API
	logger.Debug("Auto-discovering Auth0 configuration...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Make request without authentication to get Auth0 info
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/auth/info", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("User-Agent", "aphelion-cli")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to discover Auth0 config: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get Auth0 info: HTTP %d", resp.StatusCode)
	}
	
	var authInfo struct {
		Domain       string `json:"auth0_domain"`
		ClientID     string `json:"client_id"`
		Audience     string `json:"auth0_audience"`
		AuthorizeURL string `json:"auth_url"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&authInfo); err != nil {
		return fmt.Errorf("failed to decode Auth0 info: %w", err)
	}
	
	// Update profile with discovered configuration
	profile.Auth.Domain = authInfo.Domain
	// Use provided client ID if not in API response
	if authInfo.ClientID != "" {
		profile.Auth.ClientID = authInfo.ClientID
	} else {
		profile.Auth.ClientID = "UbXxpQBSr9AsqpS2ln2jzmsmamromaFC"
	}
	profile.Auth.Audience = authInfo.Audience
	
	// Set default redirect URI if not configured
	if profile.Auth.RedirectURI == "" {
		profile.Auth.RedirectURI = "http://localhost:8765/callback"
	}
	
	// Save the updated configuration
	if err := config.SetProfile(config.Get().CurrentProfile, profile); err != nil {
		logger.Warnf("Failed to save Auth0 configuration: %v", err)
		// Continue anyway - we can use the config for this session
	}
	
	// Initialize Auth0 client with discovered config
	c.auth0Client = auth.NewAuth0Client(
		profile.Auth.Domain,
		profile.Auth.ClientID,
		profile.Auth.Audience,
		profile.Auth.RedirectURI,
	)
	
	logger.Debugf("Auto-discovered Auth0 config: domain=%s, client_id=%s", 
		authInfo.Domain, authInfo.ClientID)
	
	return nil
}

// Request makes an authenticated HTTP request
func (c *Client) Request(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	url := c.baseURL + path
	
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}
	
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "aphelion-cli")
	
	// Add authentication header
	token, err := c.getValidToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get valid token: %w", err)
	}
	
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	
	logger.Debugf("Making %s request to %s", method, url)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	
	// Log response status
	logger.Debugf("Response status: %d", resp.StatusCode)
	
	return resp, nil
}

// Get makes a GET request
func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.Request(ctx, "GET", path, nil)
}

// Post makes a POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.Request(ctx, "POST", path, body)
}

// Put makes a PUT request
func (c *Client) Put(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.Request(ctx, "PUT", path, body)
}

// Delete makes a DELETE request
func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.Request(ctx, "DELETE", path, nil)
}

// getValidToken gets a valid authentication token
func (c *Client) getValidToken() (string, error) {
	if c.auth0Client == nil {
		return "", fmt.Errorf("Auth0 client not configured")
	}
	
	// Try to get stored token
	tokenInfo, err := c.auth0Client.GetStoredToken()
	if err != nil {
		// No stored token or invalid
		logger.Debug("No valid stored token found")
		return "", fmt.Errorf("not authenticated, please run 'aphelion auth login'")
	}
	
	// Check if token is close to expiring (refresh if < 5 minutes left)
	if time.Until(tokenInfo.ExpiresAt) < 5*time.Minute {
		logger.Debug("Token is close to expiring, attempting refresh")
		
		if tokenInfo.RefreshToken != "" {
			refreshed, err := c.auth0Client.RefreshToken(tokenInfo.RefreshToken)
			if err != nil {
				logger.Warnf("Failed to refresh token: %v", err)
				return "", fmt.Errorf("token expired and refresh failed, please run 'aphelion auth login'")
			}
			tokenInfo = refreshed
			logger.Debug("Token refreshed successfully")
		} else {
			return "", fmt.Errorf("token expired and no refresh token available, please run 'aphelion auth login'")
		}
	}
	
	return tokenInfo.AccessToken, nil
}

// IsAuthenticated checks if the client is authenticated
func (c *Client) IsAuthenticated() bool {
	if c.auth0Client == nil {
		return false
	}
	
	_, err := c.getValidToken()
	return err == nil
}

// GetUserInfo returns current user information
func (c *Client) GetUserInfo() (*auth.UserInfo, error) {
	if c.auth0Client == nil {
		return nil, fmt.Errorf("Auth0 client not configured")
	}
	
	tokenInfo, err := c.auth0Client.GetStoredToken()
	if err != nil {
		return nil, fmt.Errorf("not authenticated: %w", err)
	}
	
	return &tokenInfo.UserInfo, nil
}

// Login performs authentication
func (c *Client) Login(ctx context.Context) error {
	if c.auth0Client == nil {
		return fmt.Errorf("Auth0 client not configured")
	}
	
	_, err := c.auth0Client.Login(ctx)
	return err
}

// Logout clears authentication
func (c *Client) Logout() error {
	if c.auth0Client == nil {
		return fmt.Errorf("Auth0 client not configured")
	}
	
	return c.auth0Client.Logout()
}

// Health checks API health
func (c *Client) Health(ctx context.Context) error {
	resp, err := c.Get(ctx, "/health")
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}
	
	return nil
}

// HandleErrorResponse handles error responses from the API
func HandleErrorResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP %d (failed to read error response)", resp.StatusCode)
	}
	
	var errorResp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	
	if errorResp.Message != "" {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, errorResp.Message)
	}
	
	if errorResp.Error != "" {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, errorResp.Error)
	}
	
	return fmt.Errorf("HTTP %d", resp.StatusCode)
}