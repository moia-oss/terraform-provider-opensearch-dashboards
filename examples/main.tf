terraform {
  required_version = "~> 1.0"
  backend "local" {}
  required_providers {
    opensearch = {
      source  = "hashicorp.com/edu/opensearch"
      version = "0.3.1"
    }
  }
}

provider "opensearch" {
  base_url = "https://endpoint.com"
}
