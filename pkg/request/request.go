// Package request abstracts a client's request so that we can handle them in an unified way. This package borrows a lot from coredns/request.
package request

import (
	"strings"

	"github.com/miekg/dns"
)

// Request contains some connection state and is useful in plugin.
type Request struct {
	R *dns.Msg
	W dns.ResponseWriter

	name string // lowercase qname.
}

// Name will return the request's domain name
func (r *Request) Name() string {
	if r.name != "" {
		return r.name
	}
	if r.R == nil {
		r.name = "."
		return "."
	}
	if len(r.R.Question) == 0 {
		r.name = "."
		return "."
	}

	r.name = strings.ToLower(dns.Name(r.R.Question[0].Name).String())
	return r.name
}
