package match

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"

	"github.com/Xuanwo/atomdns/pkg/request"
)

// Match will check input request and tell you whether it's match or not.
type Match interface {
	Name() string

	IsMatch(r *request.Request) bool
}

// Config is the match config for atomdns
type Config struct {
	Name    string   `hcl:",label"`
	Type    string   `hcl:"type"`
	Options hcl.Body `hcl:",remain"`
}

// New will create a match.
func New(cfg *Config) (m Match, err error) {
	switch cfg.Type {
	case "in_domain_list":
		return NewInDomainList(cfg)
	default:
		return nil, fmt.Errorf("not supported match type: %s", cfg.Type)
	}
}
