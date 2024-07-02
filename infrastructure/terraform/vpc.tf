resource "aws_vpc" "example_go_api" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Cluster = "example_go_api"
  }
}
