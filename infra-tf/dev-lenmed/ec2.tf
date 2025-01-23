# Get the latest Amazon Linux 2 AMI for eu-north-1
data "aws_ami" "latest_amazon_linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]  # Amazon Linux 2 AMI
  }
}

resource "aws_instance" "golang_app" {
  ami           = data.aws_ami.latest_amazon_linux.id
  instance_type = "t3.micro"  
  vpc_security_group_ids = [aws_security_group.sg-lenmed-dev.id]
  subnet_id      = aws_subnet.sn2-lenmed-dev.id
  associate_public_ip_address = true

  depends_on = [aws_security_group.sg-lenmed-dev]  # Ensure SG is created first

  user_data = <<-EOF
              #!/bin/bash
              # Update the system
              sudo yum update -y
              sudo yum upgrade -y

              # Install Golang
              sudo yum install -y golang

              # Install Git (if not installed)
              sudo yum install -y git

              # Clone the repository (Replace with your repository URL)
              cd /home/ec2-user
              git clone https://github.com/Fastlanedevs/linmed-be.git

              # Change permissions and ownership
              sudo chown -R ec2-user:ec2-user linmed-be
              cd linmed-be

              # Build and run the Go application
              go mod tidy
              go run main.go > app.log 2>&1 & 
              EOF

  tags = {
    Name = "golang-app-instance"
  }
}
