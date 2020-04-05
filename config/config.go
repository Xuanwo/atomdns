package config

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/Xuanwo/atomdns/match"
	"github.com/Xuanwo/atomdns/upstream"
)

// Config is the main config for atomdns
type Config struct {
	Listen    string             `hcl:"listen"`
	Upstreams []*upstream.Config `hcl:"upstream,block"`
	Matches   []*match.Config    `hcl:"match,block"`
	Rules     map[string]string  `hcl:"rules"`
}

// Parse will parse file content into valid config.
func Parse(src []byte, filename string) (c *Config, err error) {
	var diags hcl.Diagnostics

	file, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("config parse: %w", diags)
	}

	c = &Config{}

	diags = gohcl.DecodeBody(file.Body, nil, c)
	if diags.HasErrors() {
		return nil, fmt.Errorf("config parse: %w", diags)
	}

	return c, nil
}

// Load will load content from file and parse into config.
func Load(filename string) (c *Config, err error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("config load: %w", err)
	}

	return Parse(content, filename)
}
