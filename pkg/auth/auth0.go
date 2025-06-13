package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/fatih/color"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/browser"
	"github.com/zalando/go-keyring"
)

const (
	keyringService = "aphelion-cli"
	tokenKey       = "auth0_token"
)

// Auth0Client handles Auth0 authentication
type Auth0Client struct {
	Domain      string
	ClientID    string
	Audience    string
	RedirectURI string
}

// TokenResponse represents Auth0 token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// TokenInfo contains parsed token information
type TokenInfo struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	UserInfo     UserInfo
}

// UserInfo contains user information from token
type UserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

// NewAuth0Client creates a new Auth0 client
func NewAuth0Client(domain, clientID, audience, redirectURI string) *Auth0Client {
	return &Auth0Client{
		Domain:      domain,
		ClientID:    clientID,
		Audience:    audience,
		RedirectURI: redirectURI,
	}
}

// Login performs Auth0 PKCE flow authentication
func (a *Auth0Client) Login(ctx context.Context) (*TokenInfo, error) {
	// Generate PKCE parameters
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code verifier: %w", err)
	}
	
	codeChallenge := generateCodeChallenge(codeVerifier)
	state := generateState()
	
	// Build authorization URL
	authURL := a.buildAuthURL(codeChallenge, state)
	
	color.Blue("Opening browser for authentication...")
	color.Cyan("If the browser doesn't open automatically, visit:")
	fmt.Printf("%s\n\n", authURL)
	
	// Open browser
	if err := browser.OpenURL(authURL); err != nil {
		color.Red("Failed to open browser: %v", err)
	}
	
	// Start local server to receive callback
	code, receivedState, err := a.startCallbackServer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to receive callback: %w", err)
	}
	
	// Verify state
	if receivedState != state {
		return nil, fmt.Errorf("state mismatch, possible CSRF attack")
	}
	
	// Exchange code for token
	tokenResp, err := a.exchangeCodeForToken(code, codeVerifier)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	
	// Parse token
	tokenInfo, err := a.parseToken(tokenResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	
	// Store token securely
	if err := a.storeToken(tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to store token: %w", err)
	}
	
	return tokenInfo, nil
}

// GetStoredToken retrieves stored token
func (a *Auth0Client) GetStoredToken() (*TokenInfo, error) {
	tokenData, err := keyring.Get(keyringService, tokenKey)
	if err != nil {
		return nil, fmt.Errorf("no stored token found: %w", err)
	}
	
	var tokenInfo TokenInfo
	if err := json.Unmarshal([]byte(tokenData), &tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to parse stored token: %w", err)
	}
	
	// Check if token is expired
	if time.Now().After(tokenInfo.ExpiresAt) {
		return nil, fmt.Errorf("stored token is expired")
	}
	
	return &tokenInfo, nil
}

// RefreshToken refreshes an expired token
func (a *Auth0Client) RefreshToken(refreshToken string) (*TokenInfo, error) {
	if refreshToken == "" {
		return nil, fmt.Errorf("no refresh token available")
	}
	
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {a.ClientID},
		"refresh_token": {refreshToken},
	}
	
	resp, err := http.PostForm(fmt.Sprintf("https://%s/oauth/token", a.Domain), data)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token refresh failed: %s", string(body))
	}
	
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	
	tokenInfo, err := a.parseToken(&tokenResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse refreshed token: %w", err)
	}
	
	// Store updated token
	if err := a.storeToken(tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to store refreshed token: %w", err)
	}
	
	return tokenInfo, nil
}

// Logout removes stored token
func (a *Auth0Client) Logout() error {
	err := keyring.Delete(keyringService, tokenKey)
	if err != nil {
		return fmt.Errorf("failed to delete stored token: %w", err)
	}
	return nil
}

// ValidateToken validates a token and returns user info
func (a *Auth0Client) ValidateToken(token string) (*UserInfo, error) {
	// Parse JWT without verification (we'll verify signature separately)
	parser := jwt.Parser{}
	claims := jwt.MapClaims{}
	
	_, _, err := parser.ParseUnverified(token, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	
	// Extract user info
	userInfo := UserInfo{}
	if sub, ok := claims["sub"].(string); ok {
		userInfo.Sub = sub
	}
	if email, ok := claims["email"].(string); ok {
		userInfo.Email = email
	}
	if emailVerified, ok := claims["email_verified"].(bool); ok {
		userInfo.EmailVerified = emailVerified
	}
	if name, ok := claims["name"].(string); ok {
		userInfo.Name = name
	}
	if picture, ok := claims["picture"].(string); ok {
		userInfo.Picture = picture
	}
	
	return &userInfo, nil
}

// Helper methods

func (a *Auth0Client) buildAuthURL(codeChallenge, state string) string {
	params := url.Values{
		"response_type":         {"code"},
		"client_id":             {a.ClientID},
		"redirect_uri":          {a.RedirectURI},
		"scope":                 {"openid profile email"},
		"code_challenge":        {codeChallenge},
		"code_challenge_method": {"S256"},
		"state":                 {state},
	}
	
	if a.Audience != "" {
		params.Set("audience", a.Audience)
	}
	
	return fmt.Sprintf("https://%s/authorize?%s", a.Domain, params.Encode())
}

func (a *Auth0Client) startCallbackServer(ctx context.Context) (string, string, error) {
	codeChan := make(chan string, 1)
	stateChan := make(chan string, 1)
	errChan := make(chan error, 1)
	
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		errorParam := r.URL.Query().Get("error")
		
		if errorParam != "" {
			errChan <- fmt.Errorf("auth error: %s", errorParam)
			return
		}
		
		if code == "" {
			errChan <- fmt.Errorf("no authorization code received")
			return
		}
		
		// Redirect to success page
		http.Redirect(w, r, "https://aphelion.exmplr.ai/auth/success", http.StatusFound)
		
		codeChan <- code
		stateChan <- state
	})
	
	server := &http.Server{
		Addr:    ":8765",
		Handler: mux,
	}
	
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("failed to start callback server: %w", err)
		}
	}()
	
	defer server.Shutdown(context.Background())
	
	select {
	case code := <-codeChan:
		state := <-stateChan
		return code, state, nil
	case err := <-errChan:
		return "", "", err
	case <-ctx.Done():
		return "", "", fmt.Errorf("authentication cancelled")
	case <-time.After(5 * time.Minute):
		return "", "", fmt.Errorf("authentication timeout")
	}
}

func (a *Auth0Client) exchangeCodeForToken(code, codeVerifier string) (*TokenResponse, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {a.ClientID},
		"code":          {code},
		"redirect_uri":  {a.RedirectURI},
		"code_verifier": {codeVerifier},
	}
	
	resp, err := http.PostForm(fmt.Sprintf("https://%s/oauth/token", a.Domain), data)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}
	
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	
	return &tokenResp, nil
}

func (a *Auth0Client) parseToken(tokenResp *TokenResponse) (*TokenInfo, error) {
	// Parse access token for user info
	userInfo, err := a.ValidateToken(tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate access token: %w", err)
	}
	
	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	
	return &TokenInfo{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    expiresAt,
		UserInfo:     *userInfo,
	}, nil
}

func (a *Auth0Client) storeToken(tokenInfo *TokenInfo) error {
	tokenData, err := json.Marshal(tokenInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}
	
	return keyring.Set(keyringService, tokenKey, string(tokenData))
}

// PKCE helper functions

func generateCodeVerifier() (string, error) {
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func generateState() string {
	data := make([]byte, 16)
	rand.Read(data)
	return base64.RawURLEncoding.EncodeToString(data)
}