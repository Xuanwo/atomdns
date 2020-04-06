package server

import (
	"fmt"
	"log"
	"time"

	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"

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

	c *cache.Cache
}

// New will create a new dns server
func New(cfg *config.Config) (s *Server, err error) {
	s = &Server{
		c: cache.New(600*time.Second, 30*time.Minute),
	}

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

	req := &request.Request{R: r}

	if v, ok := s.c.Get(req.ID()); ok {
		m := v.(*dns.Msg)
		m.Id = r.Id
		err := w.WriteMsg(m)
		if err != nil {
			log.Printf("write msg: %v", err)
			return
		}

		log.Printf("cache hit for %s", req.ID())
		return
	}

	var up upstream.Upstream

	for _, m := range s.matchers {
		if !m.IsMatch(req) {
			continue
		}
		log.Printf("rule %s matched, served via %s", m.Name(), s.rules[m.Name()].Name())

		up = s.rules[m.Name()]
		break
	}
	if up == nil {
		log.Printf("no rules matched, served via %s", s.rules["default"].Name())

		up = s.rules["default"]
	}

	m, err := up.ServeDNS(req)
	if err != nil {
		m = new(dns.Msg)
		m.SetRcode(r, dns.RcodeServerFailure)
	}

	err = w.WriteMsg(m)
	if err != nil {
		log.Printf("dns response write failed: %v", err)
		return
	}

	if m.Rcode == dns.RcodeServerFailure {
		return
	}

	if len(m.Answer) > 0 {
		s.c.Set(req.ID(), m.Copy(), time.Duration(m.Answer[0].Header().Ttl)*time.Second)
	} else {
		s.c.Set(req.ID(), m.Copy(), 0)
	}
	return
}
