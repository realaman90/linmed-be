package database

import (
	"context"

	"github.com/aakash-tyagi/linmed/models"
)

func (db *Database) AddCategory(ctx context.Context, category models.Category) error {

	_, err := db.Conn.Exec(ctx,
		`INSERT INTO categories (name, created_at, updated_at)
		VALUES ($1, $2);`,
		category.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetCategory(ctx context.Context, id string) (models.Category, error) {
	var category models.Category

	err := db.Conn.QueryRow(ctx,
		`SELECT id, name, created_at, updated_at
		FROM categories
		WHERE id = $1;`,
		id,
	).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (db *Database) GetCategories(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category

	rows, err := db.Conn.Query(ctx,
		`SELECT id, name, created_at, updated_at
		FROM categories;`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
