package e2e

import (
	"avito-shop/internal/config"
	"fmt"
	"github.com/google/uuid"
	"net"
	"testing"
	"time"
)

var e2eURL string

func TestE2E(t *testing.T) {
	cfg := config.MustLoad("../../local.e2e.env")
	e2eURL = "http://" + net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port)

	t.Run("auth", TAuth)
	t.Run("buy", TBuy)
	t.Run("send coin", TSendCoin)
	t.Run("Info", TInfo)
}

func generateUsername(prefix string) string {
	return fmt.Sprintf("%s-%s-%d", prefix, uuid.New().String()[:8], time.Now().UnixNano()%10000)
}
