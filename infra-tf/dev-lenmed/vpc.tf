resource "aws_vpc" "lenmed-dev" {
  cidr_block = "192.168.0.0/20"
  instance_tenancy = "default"
  enable_dns_support = true
}

resource "aws_subnet" "sn1-lenmed-dev" {
  cidr_block = "192.168.0.0/24"
  vpc_id = aws_vpc.lenmed-dev
  availability_zone = "eu-north-1a"
  map_public_ip_on_launch = true
}

resource "aws_subnet" "sn2-lenmed-dev" {
    cidr_block = "192.168.0.0/24"
    vpc_id = aws_vpc.lenmed-dev
    availability_zone = "eu-north-1a"
    map_public_ip_on_launch = true
} 


resource "aws_internet_gateway" "gw-lenmed-dev" {
  vpc_id = aws_vpc.lenmed-dev.id
}

resource "aws_route_table" "rt-lenmed-dev" {
  vpc_id = aws_vpc.lenmed-dev.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gw-lenmed-dev.id
  }

  route {
    ipv6_cidr_block = "::/0"
    gateway_id      = aws_internet_gateway.gw-lenmed-dev.id
  }
}