package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aakash-tyagi/linmed/config"
	database "github.com/aakash-tyagi/linmed/db"
	"github.com/aakash-tyagi/linmed/server"
	"github.com/sirupsen/logrus"
)

func main() {

	ctx := context.TODO()
	// intialise logger
	logger := logrus.New()

	// intialise config
	config, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}

	fmt.Println(config.DBUrl)

	// create database connection
	db, err := database.New(ctx, config.DBUrl)
	if err != nil {
		panic(err)
	}

	defer db.Conn.Close(ctx)
	// Proceed with table creation or other logic
	if err := db.CreateTabels(ctx); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	s := server.New(
		config,
		logger,
		db,
	)

	s.Start()
}

/*
	1. intialise config
	2. intialise postgres
	3. intialise logger
	4. intialise server
	5. start server
*/

/*
	4.1 Admin Portal
4.1.1 User Management
Manage users (staff and Super Admins)
Send invitation emails to new users
Manage user roles and permissions
4.1.2 Product Management
Manage products
Organise products by categories
Enable parent-child relationships between products
Track product statistics
4.1.3 Customer Management
Manage customers
Assign workers to customers
Upload and manage customer floor plans
Track customer statistics
4.1.4 Station Management
Manage stations for each customer
Map stations on customer floor plans
Track station statistics and alerts
Maintain station logs
4.1.5 Installed Product Management
Add products to stations
Manage product inventory at each station
Track installation dates, expiration dates, and inspection dates
Create and manage parent-child relationships between installed products
4.1.6 Dashboard
Display summary of expiring products
Show recent alerts
List pending orders
Highlight stations due for inspection
Allow customization of dashboard layout
Enable setting of alert thresholds
4.1.7 Reporting
Generate work reports (filterable by date range, worker, and customer)
Produce product statistics reports (filterable by date range and product category)
4.1.8 Data Import
Import customer data from CSV/Excel files
Provide field mapping functionality to match imported data with database fields
4.1.9 Notification Management
Configure email/SMS notifications
Select recipients for different types of alerts
4.2 Customer Portal
4.2.1 Order Management
Create order requests for products at specific stations
View order history and status
4.2.2 Station Information
View station details and logs
Access interactive floor plan with station locations
4.2.3 Reporting
Access customer-specific reports on installed products and service history

*/
