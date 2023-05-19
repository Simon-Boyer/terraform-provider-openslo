terraform {
  required_providers {
    openslo = {
      source = "registry.terraform.io/arctiq/openslo"
    }
  }
}

provider "hashicups" {}

data "hashicups_coffees" "example" {}