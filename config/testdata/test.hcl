listen = "127.0.0.1:5353"


upstream "mainland" {
  type = "udp"
  addr = "114.114.114.114:53"
}

upstream "oversea" {
  type = "udp"
  addr = "127.0.0.1:53"
}

match "to_mainland" {
  type = "in_domain_list"
  path = "/path/to/accelerated-domains.china.raw.txt"
}

rules = {
  to_mainland: "mainland",
  default: "oversea"
}