terraform {
  required_version = "~> 1.0"
  backend "local" {}
  required_providers {
    opensearch = {
      source  = "moia-oss/opensearch-dashboards"
      version = "~> 0.10"
    }
  }
}

provider "opensearch" {
  base_url = "https://endpoint.com"
}
