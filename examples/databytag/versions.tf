terraform {
  required_version = ">= 0.13"
  required_providers {

    terraform = {
      source = "terraform.io/builtin/terraform"
    }

    netbox = {
      source  = "terraform.cloud.blizzard.net/cf/netbox"
      version = "~> 0.1.8"
    }
  }
}
