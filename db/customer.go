package database

import (
	"context"

	"github.com/aakash-tyagi/linmed/models"
)

func (db *Database) AddCustomer(ctx context.Context, customer models.Customer) (int, error) {

	var id int

	err := db.Conn.QueryRow(ctx,
		`INSERT INTO customers (
		name,
		email,
		phone,
		address,
		created_at,
		updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;`,
		customer.Name, customer.Email, customer.Phone, customer.Address, customer.CreatedAt, customer.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
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
func (db *Database) GetCustomers(ctx context.Context, page, limit int) ([]models.Customer, int, error) {
	var customers []models.Customer

	rows, err := db.Conn.Query(ctx,
		`SELECT id, name, email, phone, address, created_at, updated_at
		FROM customers
		ORDER BY id
		LIMIT $1 OFFSET $2;`,
		limit, (page-1)*limit,
	)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var customer models.Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone, &customer.Address, &customer.CreatedAt, &customer.UpdatedAt); err != nil {
			return nil, 0, err
		}
		customers = append(customers, customer)
	}

	return customers, len(customers), nil

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
