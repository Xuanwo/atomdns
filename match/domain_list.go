package match

import (
	"bufio"
	"fmt"
	"os"

	iradix "github.com/hashicorp/go-immutable-radix"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/miekg/dns"

	"github.com/Xuanwo/atomdns/pkg/request"
)

type domainList struct {
	config *Config

	Path string `hcl:"path"`

	r  *iradix.Tree
	fn func(domain []byte) bool
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
	d.r = iradix.New()

	for s.Scan() {
		domain := dns.Fqdn(s.Text())

		d.r, _, _ = d.r.Insert([]byte(domain), nil)
	}

	return d, nil
}

func (d *domainList) IsMatch(r *request.Request) bool {
	return d.fn([]byte(r.Name()))
}

func (d *domainList) Name() string {
	return d.config.Name
}

// NewInDomainList will create a match which return true while request domain is in domain list.
func NewInDomainList(cfg *Config) (m Match, err error) {
	d, err := newDomainList(cfg)
	if err != nil {
		return nil, fmt.Errorf("new in domain list: %w", err)
	}

	d.fn = func(domain []byte) bool {
		_, ok := d.r.Get(domain)
		return ok
	}
	return d, nil
}
