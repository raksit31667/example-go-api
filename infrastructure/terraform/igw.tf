resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.example_go_api.id

  tags = {
    Cluster = "example_go_api"
  }
}
