package database

import (
	"context"

	"github.com/aakash-tyagi/linmed/models"
)

func (db *Database) AddCustomer(ctx context.Context, customer models.Customer) error {

	_, err := db.Conn.Exec(ctx,
		`INSERT INTO customers (name, email, phone, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6);`,
		customer.Name, customer.Email, customer.Phone, customer.Address, customer.CreatedAt, customer.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetCustomer(ctx context.Context, id string) (models.Customer, error) {
	var customer models.Customer

	err := db.Conn.QueryRow(ctx,
		`SELECT id, name, email, phone, address, created_at, updated_at
		FROM customers
		WHERE id = $1;`,
		id,
	).Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.Address, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		return customer, err
	}

	return customer, nil
}

// get all customers
func (db *Database) GetCustomers(ctx context.Context) ([]models.Customer, error) {
	var customers []models.Customer

	rows, err := db.Conn.Query(ctx,
		`SELECT id, name, email, phone, address, created_at, updated_at
		FROM customers;`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var customer models.Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.Address, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

// update customer
func (db *Database) UpdateCustomer(ctx context.Context, customer models.Customer) error {

	_, err := db.Conn.Exec(ctx,
		`UPDATE customers
		SET name = $1, email = $2, phone = $3, address = $4, updated_at = $5
		WHERE id = $6;`,
		customer.Name, customer.Email, customer.Phone, customer.Address, customer.UpdatedAt, customer.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
