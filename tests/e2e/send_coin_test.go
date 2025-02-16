package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TSendCoin(t *testing.T) {
	senderUsername := generateUsername("sender")
	receiverUsername := generateUsername("receiver")
	password := "test-password"

	senderToken, err := getAuthToken(senderUsername, password)
	if err != nil {
		t.Fatalf("failed to get sender token: %v", err)
	}
	receiverToken, err := getAuthToken(receiverUsername, password)
	if err != nil {
		t.Fatalf("failed to get receiver token: %v", err)
	}

	sendCoinURL := fmt.Sprintf("%s/api/sendCoin", e2eURL)
	sendReq := SendCoinRequest{
		ToUser: receiverUsername,
		Amount: 100,
	}
	reqBodyBytes, err := json.Marshal(sendReq)
	if err != nil {
		t.Fatalf("failed to marshal send coin request: %v", err)
	}
	resp, err := sendRequestWithAuth(http.MethodPost, sendCoinURL, bytes.NewBuffer(reqBodyBytes), senderToken)
	if err != nil {
		t.Fatalf("failed to send coin transfer request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 on coin transfer, got %v", resp.StatusCode)
	}

	infoURL := fmt.Sprintf("%s/api/info", e2eURL)

	senderInfo := waitForInfoUpdate(t, infoURL, senderToken, func(info InfoResponse) bool {
		return info.Coins == 900
	})
	if senderInfo.Coins != 900 {
		t.Errorf("expected sender coins to be 900, got %d", senderInfo.Coins)
	}

	foundSent := false
	for _, record := range senderInfo.CoinHistory.Sent {
		if record.User == receiverUsername {
			foundSent = true
			break
		}
	}
	fmt.Println(senderInfo.CoinHistory)
	fmt.Println(receiverUsername)
	if !foundSent {
		t.Errorf("sender coinHistory.sent does not include a record for transfer to %s", receiverUsername)
	}

	receiverInfo := waitForInfoUpdate(t, infoURL, receiverToken, func(info InfoResponse) bool {
		return info.Coins == 1100
	})
	if receiverInfo.Coins != 1100 {
		t.Errorf("expected receiver coins to be 1100, got %d", receiverInfo.Coins)
	}

	foundReceived := false
	for _, record := range receiverInfo.CoinHistory.Received {
		if record.User == senderUsername {
			foundReceived = true
			break
		}
	}
	if !foundReceived {
		t.Errorf("receiver coinHistory.received does not include a record for transfer from %s", senderUsername)
	}
}

func waitForInfoUpdate(t *testing.T, infoURL, token string, check func(info InfoResponse) bool) InfoResponse {
	var infoResp InfoResponse
	timeout := time.After(3 * time.Second)
	tick := time.Tick(200 * time.Millisecond)
	for {
		select {
		case <-timeout:
			t.Fatalf("timeout waiting for info update")
		case <-tick:
			resp, err := sendRequestWithAuth(http.MethodGet, infoURL, nil, token)
			if err != nil {
				continue
			}
			if resp.StatusCode != http.StatusOK {
				resp.Body.Close()
				continue
			}
			err = decodeResponse(resp, &infoResp)
			resp.Body.Close()
			if err != nil {
				continue
			}
			if check(infoResp) {
				return infoResp
			}
		}
	}
}
