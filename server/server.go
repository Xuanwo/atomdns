package server

import (
	"fmt"
	"log"

	"github.com/miekg/dns"

	"github.com/Xuanwo/atomdns/config"
	"github.com/Xuanwo/atomdns/match"
	"github.com/Xuanwo/atomdns/pkg/request"
	"github.com/Xuanwo/atomdns/upstream"
)

// Server is the dns server.
type Server struct {
	upsrteams map[string]upstream.Upstream // upstream_name -> upstream
	matchers  []match.Match
	rules     map[string]upstream.Upstream // match rules -> upstream
}

// New will create a new dns server
func New(cfg *config.Config) (s *Server, err error) {
	s = &Server{}

	// Setup streams
	s.upsrteams = make(map[string]upstream.Upstream)
	for _, v := range cfg.Upstreams {
		up, err := upstream.New(v)
		if err != nil {
			return nil, fmt.Errorf("upsteam new: %w", err)
		}
		s.upsrteams[up.Name()] = up
	}

	// Setup matches.
	s.matchers = make([]match.Match, 0, len(cfg.Matches))
	for _, v := range cfg.Matches {
		m, err := match.New(v)
		if err != nil {
			return nil, fmt.Errorf("match new: %w", err)
		}
		s.matchers = append(s.matchers, m)
	}

	// Setup rules.
	s.rules = make(map[string]upstream.Upstream)
	for k, v := range cfg.Rules {
		s.rules[k] = s.upsrteams[v]
	}
	return s, nil
}

// ServeDNS implements dns.Handler
func (s *Server) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	log.Print(r.Question)

	req := &request.Request{R: r, W: w}

	for _, m := range s.matchers {
		if !m.IsMatch(req) {
			continue
		}
		log.Printf("rule %s matched, served via %s", m.Name(), s.rules[m.Name()].Name())
		s.rules[m.Name()].ServeDNS(req)
		return
	}

	log.Printf("no rules matched, served via %s", s.rules["default"].Name())
	s.rules["default"].ServeDNS(req)
	return
}
