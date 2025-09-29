provider "aws" {
  region  = "ap-southeast-2"
  profile = "default"
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    values = ["ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-arm64-server-*"]
    name   = "name"
  }
  owners = ["099720109477"]
}

resource "aws_instance" "hpa_server" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t4g.micro"
  key_name      = aws_key_pair.this.key_name

  subnet_id                   = module.vpc.public_subnets[0]
  associate_public_ip_address = true
  vpc_security_group_ids      = [aws_security_group.web_access.id]

  tags = {
    Name = "hpa_server"
  }
}
output "ec2_public_dns" {
  value = aws_instance.hpa_server.public_dns
}

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "6.3.0"

  name = "main_vpc"
  cidr = "10.0.0.0/16"

  azs             = ["ap-southeast-2a"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  public_subnets  = ["10.0.101.0/24"]

  enable_dns_hostnames = true
}

resource "aws_security_group" "web_access" {
  name   = "web_access"
  vpc_id = module.vpc.vpc_id

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
resource "tls_private_key" "this" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "this" {
  key_name   = "priv_key"
  public_key = tls_private_key.this.public_key_openssh
}

resource "local_file" "private_key" {
  content  = tls_private_key.this.private_key_pem
  filename = "priv_key"
}
