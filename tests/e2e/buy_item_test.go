package e2e

import (
	"fmt"
	"net/http"
	"testing"
)

func TBuy(t *testing.T) {
	const password = "test-password"
	username := generateUsername("buy-user")

	token, err := getAuthToken(username, password)
	if err != nil {
		t.Fatalf("failed to get auth token: %v", err)
	}

	buyURL := fmt.Sprintf("%s/api/buy/t-shirt", e2eURL)
	resp, err := sendRequestWithAuth(http.MethodGet, buyURL, nil, token)
	if err != nil {
		t.Fatalf("failed to send buy request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 on buy, got %v", resp.StatusCode)
	}

	infoURL := fmt.Sprintf("%s/api/info", e2eURL)
	infoRespHTTP, err := sendRequestWithAuth(http.MethodGet, infoURL, nil, token)
	if err != nil {
		t.Fatalf("failed to get info: %v", err)
	}
	defer infoRespHTTP.Body.Close()
	if infoRespHTTP.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 on info, got %v", infoRespHTTP.StatusCode)
	}

	var infoResp InfoResponse
	if err := decodeResponse(infoRespHTTP, &infoResp); err != nil {
		t.Fatalf("failed to decode info response: %v", err)
	}

	expectedCoins := 1000 - 80
	if infoResp.Coins != expectedCoins {
		t.Errorf("expected coins %d, got %d", expectedCoins, infoResp.Coins)
	}

	found := false
	for _, item := range infoResp.Inventory {
		if item.Type == "t-shirt" && item.Quantity == 1 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected to find one t-shirt in inventory, but it was not present")
	}
}
