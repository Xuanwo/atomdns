package upstream

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"github.com/Xuanwo/atomdns/pkg/request"
)

// Upstream will actually handle a dns request.
type Upstream interface {
	ServeDNS(r *request.Request)

	Name() string
}

// Config is the upstream config for atomdns
type Config struct {
	Name    string   `hcl:",label"`
	Type    string   `hcl:"type"`
	Addr    string   `hcl:"addr"`
	Options hcl.Body `hcl:",remain"`
}

// New will create a new upstream
func New(cfg *Config) (up Upstream, err error) {
	switch cfg.Type {
	case "udp":
		return NewUDPClient(cfg)
	case "tcp":
		return NewTCPClient(cfg)
	case "dot":
		return NewDoTClient(cfg)
	default:
		return nil, fmt.Errorf("not supported upstream type: %s", cfg.Type)
	}
}
