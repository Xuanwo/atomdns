package main

import (
	"log"
	"os"

	"github.com/miekg/dns"

	"github.com/Xuanwo/atomdns/config"
	"github.com/Xuanwo/atomdns/server"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no config input")
	}

	cfg, err := config.Load(os.Args[1])
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	s, err := server.New(cfg)
	if err != nil {
		log.Fatalf("server new failed: %v", err)
	}

	err = dns.ListenAndServe(cfg.Listen, "udp", s)
	if err != nil {
		log.Fatalf("dns server exited: %v", err)
	}
}
