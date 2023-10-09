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

resource "example_example" "myExample" {
  string_val = "whatever"
  // Always get a dirty plan when this attribute is commented out
  /*attribute_mapping = {
    attribute_contract_fulfillment = {
      "entryUUID" = {
        source = {
          type = "ADAPTER"
        },
        value = "entryUUID"
      }
      "policy.action" = {
        source = {
          type = "ADAPTER"
        },
        value = "policy.action"
      },
      "username" = {
        source = {
          type = "ADAPTER"
        },
        value = "username"
      }
    },
    attribute_sources = [],
    issuance_criteria = {
      conditional_criteria = []
    }
  }*/
}
