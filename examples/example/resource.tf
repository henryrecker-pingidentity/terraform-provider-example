terraform {
  required_version = ">=1.1"
  required_providers {
    example = {
      version = "~> 0.0.1"
      source  = "example/example"
    }
  }
}

provider "example" {
}

import {
  to = example_import_error.generated
  id = "idVal"
}

# resource "example_import_error" "ex" {
# }
