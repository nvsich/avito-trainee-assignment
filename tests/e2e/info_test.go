package e2e

import (
	"fmt"
	"net/http"
	"testing"
)

func TInfo(t *testing.T) {
	username := generateUsername("info-user")
	password := "test-password"
	token, err := getAuthToken(username, password)
	if err != nil {
		t.Fatalf("failed to get auth token: %v", err)
	}

	infoURL := fmt.Sprintf("%s/api/info", e2eURL)
	resp, err := sendRequestWithAuth(http.MethodGet, infoURL, nil, token)
	if err != nil {
		t.Fatalf("failed to get info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 on info, got %v", resp.StatusCode)
	}

	var infoResp InfoResponse
	if err := decodeResponse(resp, &infoResp); err != nil {
		t.Fatalf("failed to decode info response: %v", err)
	}

	if infoResp.Coins != 1000 {
		t.Errorf("expected initial coins to be 1000, got %d", infoResp.Coins)
	}

	if len(infoResp.Inventory) != 0 {
		t.Errorf("expected inventory to be empty, got %d items", len(infoResp.Inventory))
	}

	if len(infoResp.CoinHistory.Received) != 0 || len(infoResp.CoinHistory.Sent) != 0 {
		t.Errorf("expected coinHistory to be empty, got received %d and sent %d items",
			len(infoResp.CoinHistory.Received), len(infoResp.CoinHistory.Sent))
	}
}
