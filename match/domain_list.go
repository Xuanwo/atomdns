package match

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/miekg/dns"

	"github.com/Xuanwo/atomdns/pkg/request"
)

type domainList struct {
	config *Config

	Path string `hcl:"path"`

	m  map[string]struct{}
	fn func(domain string) bool
}

func newDomainList(cfg *Config) (d *domainList, err error) {
	d = &domainList{config: cfg}

	var diags hcl.Diagnostics
	diags = gohcl.DecodeBody(cfg.Options, nil, d)
	if diags.HasErrors() {
		return nil, fmt.Errorf("new domain list: %w", diags)
	}

	file, err := os.Open(d.Path)
	if err != nil {
		return nil, fmt.Errorf("new domain list: %w", err)
	}

	s := bufio.NewScanner(file)
	d.m = make(map[string]struct{})

	var domain string
	for s.Scan() {
		domain = dns.Fqdn(s.Text())

		d.m[domain] = struct{}{}
	}

	return d, nil
}

func (d *domainList) IsMatch(r *request.Request) bool {
	return d.fn(r.Name())
}

func (d *domainList) Name() string {
	return d.config.Name
}

func (d *domainList) in(domain string) bool {
	s := strings.Split(domain, ".")
	for i := 0; i < len(s)-1; i++ {
		_, ok := d.m[strings.Join(s[i:], ".")]
		if ok {
			return true
		}
	}
	return false
}

// NewInDomainList will create a match which return true while request domain is in domain list.
func NewInDomainList(cfg *Config) (m Match, err error) {
	d, err := newDomainList(cfg)
	if err != nil {
		return nil, fmt.Errorf("new in domain list: %w", err)
	}

	d.fn = func(domain string) bool {
		return d.in(domain)
	}
	return d, nil
}
