package upstream

import (
	"log"

	"github.com/miekg/dns"

	"github.com/Xuanwo/atomdns/pkg/request"
)

type client struct {
	config *Config

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

// NewUDPClient create a new udp client.
func NewUDPClient(cfg *Config) (u Upstream, err error) {
	c := &client{config: cfg}
	c.c = &dns.Client{}
	return c, nil
}
