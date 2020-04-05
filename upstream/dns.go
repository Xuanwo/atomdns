package upstream

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/miekg/dns"

	"github.com/Xuanwo/atomdns/pkg/request"
)

type client struct {
	config *Config

	// DoT related config
	TLSServerName string `hcl:"tls_server_name,optional"`

	c *dns.Client
}

// Name implements upstream.Name
func (c *client) Name() string {
	return c.config.Name
}

// ServeDNS implements upstream.ServeDNS
func (c *client) ServeDNS(r *request.Request) {
	m, _, err := c.c.Exchange(r.R, c.config.Addr)
	if err != nil {
		log.Printf("serve dns: %v", err)

		m = new(dns.Msg)
		m.SetRcode(r.R, dns.RcodeServerFailure)
	}
	err = r.W.WriteMsg(m)
	if err != nil {
		log.Printf("write msg: %v", err)
	}
}

// NewTCPClient create a new tcp client.
func NewTCPClient(cfg *Config) (u Upstream, err error) {
	c := &client{config: cfg}
	c.c = &dns.Client{
		Net: "tcp",
	}
	return c, nil
}

// NewDoTClient create a new dot client
func NewDoTClient(cfg *Config) (u Upstream, err error) {
	c := &client{config: cfg}

	var diags hcl.Diagnostics
	diags = gohcl.DecodeBody(cfg.Options, nil, c)
	if diags.HasErrors() {
		return nil, fmt.Errorf("new domain list: %w", diags)
	}

	c.c = &dns.Client{
		Net: "tcp-tls",
		TLSConfig: &tls.Config{
			ServerName: "",
		},
	}
	return c, nil
}

// NewUDPClient create a new udp client.
func NewUDPClient(cfg *Config) (u Upstream, err error) {
	c := &client{config: cfg}
	c.c = &dns.Client{}
	return c, nil
}
