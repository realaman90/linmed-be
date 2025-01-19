package main

import (
	"context"
	"log"

	"github.com/aakash-tyagi/linmed/aws"
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

	// intialise s3 bucket

	s3Client, err := aws.NewS3Client(config.Region, config.AccessId, config.AcessKey)
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}

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
		s3Client,
	)

	s.Start()
}
