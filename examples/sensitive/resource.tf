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

resource "example_sensitive" "sensitive" {
  tables = [
    {
      rows = [
        {
          fields = [
            {
              name  = "Key ID"
              value = "jwtSymmetricKey1"
            },
            {
              name  = "Encoding"
              value = "b64u"
            }
          ]
          sensitive_fields = [
            {
              name  = "Key"
              value = "Asdf"
            },
          ]
          # If this attribute is uncommented, the unexpected plans do not occur
          # default_row = false
        }
      ]
    },
  ]
}
