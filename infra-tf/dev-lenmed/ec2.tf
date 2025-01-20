resource "aws_instance" "golang_app" {
  ami           = "ami-0c55b159cbfafe1f0" 
  instance_type = "t2.micro"  # Choose the instance type based on your needs
  security_groups = [aws_security_group.sg-lenmed-dev.name]
  subnet_id      = aws_subnet.sn2-lenmed-dev.id
  associate_public_ip_address = true

  user_data = <<-EOF
              #!/bin/bash
              sudo apt-get update -y
              sudo apt-get install -y golang-go
              # Add your deployment script here (e.g., clone repo, build app)
              EOF

  tags = {
    Name = "golang-app-instance"
  }
}
