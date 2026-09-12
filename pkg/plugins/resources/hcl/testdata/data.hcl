resource "person" "john" {
  first_name  = "John"
  middle_name = ""
  surname     = "Doe"
  age         = 30

  nested {
    attr1 = "val1"
  }
}

locals {
  lets_encrypt_dns_challenged_domains = {
    "trusted.ci.jenkins.io" = "2024-04-03T20:00:00Z"
    "cert.ci.jenkins.io"    = "2024-04-03T21:00:00Z"
  }
}
