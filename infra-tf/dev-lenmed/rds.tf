resource "aws_db_instance" "lenmed_db" {
  identifier        = "lenmed-db-instance"
  engine            = "postgres"
  engine_version    = "16"  
  instance_class    = "db.t3.micro"  # Choose an instance type
  allocated_storage = 5  # GB, adjust as necessary
  storage_type      = "gp2"
  db_name           = "lenmed_db"  # Your database name
  username          = "lenmed"  # Your admin username
  password          = "lenmed123"  # Your admin password
  parameter_group_name = "default.postgres13"  # PostgreSQL version
  multi_az          = false  # Set to true for multi-AZ deployment for high availability
  publicly_accessible = false  # Set to true if the DB should be publicly accessible
  vpc_security_group_ids = [aws_security_group.sg-lenmed-dev.id]
  db_subnet_group_name = aws_db_subnet_group.lenmed-db-subnet-group.name
  tags = {
    Name = "lenmed-db-instance"
  }

  # Automatically backups and maintenance options
  backup_retention_period = 7  # Retain backups for 7 days
  maintenance_window       = "Mon:00:00-Mon:03:00"
}

resource "aws_db_subnet_group" "lenmed-db-subnet-group" {
  name       = "lenmed-db-subnet-group"
  subnet_ids = [aws_subnet.sn1-lenmed-dev.id, aws_subnet.sn2-lenmed-dev.id]  # Ensure subnets are private

  tags = {
    Name = "lenmed-db-subnet-group"
  }
}
