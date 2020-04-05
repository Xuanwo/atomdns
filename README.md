# atomdns

AtomDNS is built for myself.

## Install

### Archlinux with AUR

```go
yay -S atomdns
```

## Usage

```bash
atomdns /path/to/atomdns.hcl
```

## Config

### Forward all dns query to upstream

```hcl
listen = "127.0.0.1:53"

upstream "default" {
  type = "dot"
  addr = "185.222.222.222:853"
  tls_server_name = "public-dns-a.dns.sb"
}

rules = {
  default: "default"
}
```

### Accelerate chinese domains

```hcl
listen = "127.0.0.1:53"


upstream "oversea" {
  type = "dot"
  addr = "185.222.222.222:853"
  tls_server_name = "public-dns-a.dns.sb"
}

upstream "mainland" {
  type = "udp"
  addr = "114.114.114.114:53"
}

match "to_mainland" {
  type = "in_domain_list"
  # get this file from https://github.com/felixonmars/dnsmasq-china-list
  path = "/etc/atomdns/accelerated-domains.china.raw.txt"
}

rules = {
  to_mainland: "mainland",
  default: "oversea"
}
```
