package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TAuth(t *testing.T) {
	newUser := generateUsername("new user")
	testPassword := "test-password"
	wrongPassword := "wrong-password"
	tests := []struct {
		name           string
		requestBody    AuthRequest
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "new user",
			requestBody: AuthRequest{
				Username: newUser,
				Password: testPassword,
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name: "wrong password",
			requestBody: AuthRequest{
				Username: newUser,
				Password: wrongPassword,
			},
			expectedStatus: http.StatusUnauthorized,
			expectToken:    false,
		},
		{
			name: "existing user",
			requestBody: AuthRequest{
				Username: newUser,
				Password: testPassword,
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := sendAuthRequest(tt.requestBody)
			if err != nil {
				t.Fatalf("failed to send request: %v: ", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, resp.StatusCode)
			}

			var authResponse AuthResponse
			var errorResponse ErrorResponse
			if err = decodeResponse(resp, &authResponse); err != nil {
				if err = decodeResponse(resp, &errorResponse); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
			}

			if tt.expectToken && authResponse.Token == "" {
				t.Errorf("expected token in response, got none")
			}
		})
	}
}

func getAuthToken(username, password string) (string, error) {
	resp, err := sendAuthRequest(AuthRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("sendAuthRequest failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected auth status: %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := decodeResponse(resp, &authResp); err != nil {
		return "", fmt.Errorf("failed to decode auth response: %w", err)
	}
	if authResp.Token == "" {
		return "", fmt.Errorf("empty token received")
	}
	return authResp.Token, nil
}

func sendAuthRequest(body AuthRequest) (*http.Response, error) {
	requestBodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, e2eURL+"/api/auth", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(req)
}

func sendRequestWithAuth(method, url string, body io.Reader, token string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	client := &http.Client{}
	return client.Do(req)
}
