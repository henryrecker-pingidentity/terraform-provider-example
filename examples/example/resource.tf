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
  # Causes infinite plans
  float64 =  242.08120431461208
  # Does not cause infinite plans - last digit changed
  # float64 =  242.08120431461209
}
