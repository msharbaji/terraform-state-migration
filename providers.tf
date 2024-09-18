terraform {
#   backend "local" {
#     path = "terraform.tfstate"
#   }
    backend "pg" {
      conn_str    = "postgres://localhost:5432/terraform_backend?sslmode=disable"
      schema_name = "terraform_remote_state"
      skip_index_creation = true
      skip_schema_creation = true
      skip_table_creation  = true
    }

  required_providers {
    null = {
      source = "hashicorp/null"
    }
    local = {
      source = "hashicorp/local"
    }
  }
}

# providers.tf
provider "null" {}
provider "local" {}