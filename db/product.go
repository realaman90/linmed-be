package database

import (
	"context"

	"github.com/aakash-tyagi/linmed/models"
)

func (db *Database) AddCategory(ctx context.Context, category models.Category) (int, error) {

	// return id of the category

	var id int

	err := db.Conn.QueryRow(ctx,
		`INSERT INTO categories (name, description, color)
		VALUES ($1, $2, $3)
		RETURNING id;`,
		category.Name,
		category.Description,
		category.Color,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Database) GetCategory(ctx context.Context, id string) (models.Category, error) {
	var category models.Category

	err := db.Conn.QueryRow(ctx,
		`SELECT id, name, description,color,created_at, updated_at
		FROM categories
		WHERE id = $1;`,
		id,
	).Scan(&category.ID, &category.Name, &category.Description, &category.Color, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return category, err
	}

	return category, nil
}

func (db *Database) UpdateCategory(ctx context.Context, id string, category models.Category) error {
	_, err := db.Conn.Exec(ctx,
		`UPDATE categories
		SET name = $1, description = $2, color = $3
		WHERE id = $4;`,
		category.Name,
		category.Description,
		category.Color,
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
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.Color, &category.CreatedAt, &category.UpdatedAt); err != nil {
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

	// Fetch the main product (parent)
	err := db.Conn.QueryRow(ctx,
		`SELECT id, name, category_id, category_name,price, description, image_url, parent_id, coverage_amount, age_limit, created_at, updated_at
		FROM products
		WHERE id = $1;`,
		id,
	).Scan(&product.ID, &product.Name, &product.CategoryID, &product.CategoryName, &product.Price, &product.Description,
		&product.ImageURL, &product.ParentID, &product.CoverageAmount, &product.AgeLimit,
		&product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return product, err
	}

	// Fetch the children (products that have the current product as their parent)
	rows, err := db.Conn.Query(ctx,
		`SELECT id, name, category_id, category_name,price, description, image_url, parent_id, coverage_amount, age_limit, created_at, updated_at
		FROM products
		WHERE parent_id = $1;`,
		product.ID,
	)
	if err != nil {
		return product, err
	}
	defer rows.Close()

	for rows.Next() {
		var child models.Product
		err := rows.Scan(&child.ID, &child.Name, &child.CategoryID, &child.CategoryName, &child.Price, &child.Description,
			&child.ImageURL, &child.ParentID, &child.CoverageAmount, &child.AgeLimit,
			&child.CreatedAt, &child.UpdatedAt)
		if err != nil {
			return product, err
		}
		product.Children = append(product.Children, child)
	}

	return product, nil
}

func (db *Database) GetProducts(ctx context.Context, page, limit int) ([]models.Product, int, error) {
	var products []models.Product

	// Query parent products with category names
	rows, err := db.Conn.Query(ctx,
		`SELECT p.id, p.name, p.category_id, p.price, p.description, p.image_url, p.parent_id, 
		        p.coverage_amount, p.age_limit, p.created_at, p.updated_at, c.name as category_name
		FROM products p
		INNER JOIN categories c ON p.category_id = c.id
		WHERE p.parent_id IS NULL
		LIMIT $1 OFFSET $2;`,
		limit, (page-1)*limit,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
			&product.ID, &product.Name, &product.CategoryID, &product.Price,
			&product.Description, &product.ImageURL, &product.ParentID, &product.CoverageAmount,
			&product.AgeLimit, &product.CreatedAt, &product.UpdatedAt, &product.CategoryName,
		); err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	// Retrieve children for each product
	for i := range products {
		childRows, err := db.Conn.Query(ctx,
			`SELECT id, name, category_id, price, description, image_url, parent_id, 
			        coverage_amount, age_limit, created_at, updated_at
			FROM products 
			WHERE parent_id = $1;`,
			products[i].ID,
		)
		if err != nil {
			return nil, 0, err
		}
		defer childRows.Close()

		for childRows.Next() {
			var child models.Product
			if err := childRows.Scan(
				&child.ID, &child.Name, &child.CategoryID, &child.Price,
				&child.Description, &child.ImageURL, &child.ParentID, &child.CoverageAmount,
				&child.AgeLimit, &child.CreatedAt, &child.UpdatedAt,
			); err != nil {
				return nil, 0, err
			}
			products[i].Children = append(products[i].Children, child)
		}
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
