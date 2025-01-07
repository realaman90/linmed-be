package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	Conn *pgx.Conn
}

func New(ctx context.Context, dbUrl string) (*Database, error) {

	dbconn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database %s: %w", dbUrl, err)
	}

	if err := dbconn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &Database{Conn: dbconn}, nil
}

func (db *Database) CreateTabels(ctx context.Context) error {

	_, err := db.Conn.Exec(ctx,
		`CREATE TABLE categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);
		`)
	if err != nil {
		return err
	}

	// Example: Create Users Table
	_, err = db.Conn.Exec(ctx,
		`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR NOT NULL,
			email VARCHAR UNIQUE NOT NULL,
			password_hash VARCHAR NOT NULL,
			first_name VARCHAR,
			last_name VARCHAR,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
			);
		`)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec(ctx,
		`
		CREATE TABLE IF NOT EXISTS products (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			parent_id INT REFERENCES products(id),
			category_id INT REFERENCES categories(id) NOT NULL,
			price DECIMAL(10, 2) NOT NULL,
			coverage_amount DECIMAL(10, 2),
			age_limit INT,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec(ctx,
		`
		CREATE TABLE IF NOT EXISTS customers (
    		id SERIAL PRIMARY KEY,
    		name VARCHAR(100) NOT NULL,
    		email VARCHAR(100) NOT NULL UNIQUE,
    		phone VARCHAR(15),
    		address VARCHAR(255),
    		created_at TIMESTAMP DEFAULT NOW(),
    		updated_at TIMESTAMP DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec(ctx,
		`
		CREATE TABLE IF NOT EXISTS floor_plans (
    		 id SERIAL PRIMARY KEY,
   			 name VARCHAR(100) NOT NULL,
             layout TEXT,
   			 customer_id INT REFERENCES customers(id),
   			 created_at TIMESTAMP DEFAULT NOW(),
   			 updated_at TIMESTAMP DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec(ctx,
		`
		CREATE TABLE IF NOT EXISTS stations (
    		 id SERIAL PRIMARY KEY,
    		 name VARCHAR(100) NOT NULL,
    		 description TEXT,
    		 customer_id INT REFERENCES customers(id),
     		 floor_plan_id INT REFERENCES floor_plans(id),
    		 created_at TIMESTAMP DEFAULT NOW(),
   			 updated_at TIMESTAMP DEFAULT NOW()
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec(ctx,
		`
		CREATE TABLE IF NOT EXISTS station_products (
    		 station_id INT REFERENCES stations(id) ON DELETE CASCADE,
   			 product_id INT REFERENCES products(id) ON DELETE CASCADE,
   			 PRIMARY KEY (station_id, product_id)
		);
	`)
	if err != nil {
		return err
	}

	return nil
}
