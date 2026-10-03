package transport

import (
	"strings"
	"testing"

	"github.com/anonvector/slipgate/internal/config"
)

func buildTestCaddyfile(decoy string) string {
	return buildCaddyfile(&config.TunnelConfig{
		Tag:       "naive",
		Transport: config.TransportNaive,
		Domain:    "n.example.org",
		Enabled:   true,
		Naive: &config.NaiveConfig{
			Email:    "admin@example.org",
			DecoyURL: decoy,
			Port:     443,
			User:     "user1",
			Password: "secretpass",
		},
	})
}

// HTTP/3 (QUIC) must stay disabled: Caddy opens a UDP listener for it and treats
// a bind failure as fatal, which crash-loops the service whenever another daemon
// owns UDP/<port> (e.g. AmneziaWG's proxy on 443/udp). SlipGate's firewall only
// opens TCP/<port>, so HTTP/3 is unreachable from outside anyway.
func TestBuildCaddyfileDisablesHTTP3(t *testing.T) {
	got := buildTestCaddyfile("https://www.wikipedia.org")

	if !strings.Contains(got, "protocols h1 h2") {
		t.Errorf("Caddyfile does not restrict protocols to h1 h2; Caddy will try to bind UDP/443:\n%s", got)
	}
}
