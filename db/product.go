package database

import (
	"context"

	"github.com/aakash-tyagi/linmed/models"
)

func (db *Database) AddCategory(ctx context.Context, category models.Category) (int, error) {

	// return id of the category

	var id int

	err := db.Conn.QueryRow(ctx,
		`INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id;`,
		category.Name,
		category.Description,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Database) GetCategory(ctx context.Context, id string) (models.Category, error) {
	var category models.Category

	err := db.Conn.QueryRow(ctx,
		`SELECT id, name, description,created_at, updated_at
		FROM categories
		WHERE id = $1;`,
		id,
	).Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (db *Database) UpdateCategory(ctx context.Context, id string, category models.Category) error {
	_, err := db.Conn.Exec(ctx,
		`UPDATE categories
		SET name = $1, description = $2
		WHERE id = $3;`,
		category.Name,
		category.Description,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteCategory(ctx context.Context, id string) error {
	_, err := db.Conn.Exec(ctx,
		`DELETE FROM categories
		WHERE id = $1;`,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetCategories(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category

	rows, err := db.Conn.Query(ctx,
		`SELECT id, name, description,created_at, updated_at
		FROM categories;`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (db *Database) AddProduct(ctx context.Context, product models.Product) (uint, error) {

	// return id of the product

	var id uint

	err := db.Conn.QueryRow(ctx,
		`INSERT INTO products (name, category_id, price, description, image_url, parent_id, coverage_amount, age_limit, children)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id;`,
		product.Name,
		product.CategoryID,
		product.Price,
		product.Description,
		product.ImageURL,
		product.ParentID,
		product.CoverageAmount,
		product.AgeLimit,
		product.Children,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil

}

func (db *Database) GetProduct(ctx context.Context, id string) (models.Product, error) {

	var product models.Product

	err := db.Conn.QueryRow(ctx,
		`SELECT id, name, category_id, price, description, image_url, parent_id, coverage_amount, age_limit, children, created_at, updated_at
		FROM products
		WHERE id = $1;`,
		id,
	).Scan(&product.ID, &product.Name, &product.CategoryID, &product.Price, &product.Description, &product.ImageURL, &product.ParentID, &product.CoverageAmount, &product.AgeLimit, &product.Children, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (db *Database) GetProducts(ctx context.Context, page, limit int) ([]models.Product, int, error) {
	var products []models.Product

	// get category name also with product

	rows, err := db.Conn.Query(ctx,
		`SELECT p.id, p.name, p.category_id, p.price, p.description, p.image_url, p.parent_id, p.coverage_amount, p.age_limit, p.children, p.created_at, p.updated_at, c.name
		FROM products p
		INNER JOIN categories c ON p.category_id = c.id
		LIMIT $1 OFFSET $2;`,
		limit, page,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.CategoryID, &product.Price, &product.Description, &product.ImageURL, &product.ParentID, &product.CoverageAmount, &product.AgeLimit, &product.Children, &product.CreatedAt, &product.UpdatedAt, &product.CategoryName); err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	return products, len(products), nil
}

func (db *Database) UpdateProduct(ctx context.Context, product models.Product) error {

	_, err := db.Conn.Exec(ctx,
		`UPDATE products
		SET name = $1, category_id = $2, price = $3, description = $4, image_url = $5, parent_id = $6, coverage_amount = $7, age_limit = $8, children = $9
		WHERE id = $10;`,
		product.Name,
		product.CategoryID,
		product.Price,
		product.Description,
		product.ImageURL,
		product.ParentID,
		product.CoverageAmount,
		product.AgeLimit,
		product.Children,
		product.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteProduct(ctx context.Context, id string) error {

	_, err := db.Conn.Exec(ctx,
		`DELETE FROM products
		WHERE id = $1;`,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}
