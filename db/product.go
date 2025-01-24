package database

import (
	"context"
	"database/sql"

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

func (db *Database) GetCategories(ctx context.Context) (models.CategoriesResponse, error) {
	var categories []models.Category

	rows, err := db.Conn.Query(ctx, `
		SELECT c.id, c.name, c.description, c.color, c.created_at, c.updated_at, COUNT(p.id) as product_count
		FROM categories c
		LEFT JOIN products p ON p.category_id = c.id
		GROUP BY c.id`)
	if err != nil {
		return models.CategoriesResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.Color,
			&category.CreatedAt, &category.UpdatedAt, &category.ProductCount); err != nil {
			return models.CategoriesResponse{}, err
		}
		categories = append(categories, category)
	}

	return models.CategoriesResponse{Categories: categories}, nil
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
		`SELECT id, name, category_id,price, description, image_url, parent_id, coverage_amount, age_limit, created_at, updated_at
		FROM products
		WHERE id = $1;`,
		id,
	).Scan(&product.ID, &product.Name, &product.CategoryID, &product.Price, &product.Description,
		&product.ImageURL, &product.ParentID, &product.CoverageAmount, &product.AgeLimit,
		&product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return product, err
	}

	// Fetch the children (products that have the current product as their parent)
	rows, err := db.Conn.Query(ctx,
		`SELECT id, name, category_id,price, description, image_url, parent_id, coverage_amount, age_limit, created_at, updated_at
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
		err := rows.Scan(&child.ID, &child.Name, &child.CategoryID, &child.Price, &child.Description,
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

	var total int

	// Get total count of products
	err := db.Conn.QueryRow(ctx, "SELECT COUNT(*) FROM products").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch all products, including parent-child relations
	rows, err := db.Conn.Query(ctx,
		`SELECT p.id, p.name, p.category_id, c.name as category_name, p.price, p.description, 
                p.image_url, p.parent_id, p.coverage_amount, p.age_limit, 
                p.created_at, p.updated_at
		 FROM products p
		 LEFT JOIN categories c ON p.category_id = c.id
		 ORDER BY p.id
		 LIMIT $1 OFFSET $2`, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Map to store products by ID
	productMap := make(map[uint]*models.Product)

	// Process query results
	for rows.Next() {
		var p models.Product
		var parentID sql.NullInt64

		if err := rows.Scan(&p.ID, &p.Name, &p.CategoryID, &p.CategoryName, &p.Price, &p.Description,
			&p.ImageURL, &parentID, &p.CoverageAmount, &p.AgeLimit, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}

		// Set ParentID properly
		if parentID.Valid {
			p.ParentID = uintPtr(uint(parentID.Int64))
		}

		// Store the product in the map
		productMap[p.ID] = &p
	}

	// Construct parent-child relationships
	var result []models.Product
	for _, product := range productMap {
		if product.ParentID != nil {
			// If the product has a parent, add it to its parent's children
			parentProduct := productMap[*product.ParentID]
			if parentProduct != nil {
				parentProduct.Children = append(parentProduct.Children, *product)
			}
		} else {
			// If no parent, add to the main list
			result = append(result, *product)
		}
	}

	// Ensure to include nested children
	for i := range result {
		result[i].Children = getNestedChildren(productMap, result[i])
	}

	return result, total, nil
}

// Helper function to create a pointer from uint
func uintPtr(i uint) *uint {
	return &i
}

// Helper function to get nested children
func getNestedChildren(productMap map[uint]*models.Product, parent models.Product) []models.Product {
	var children []models.Product
	for _, product := range productMap {
		if product.ParentID != nil && *product.ParentID == parent.ID {
			children = append(children, *product)
			product.Children = getNestedChildren(productMap, *product) // Recursively get nested children
		}
	}
	return children
}

// Helper function to get child products
// func (db *Database) getChildren(ctx context.Context, parentID uint) ([]models.Product, error) {
// 	var children []models.Product

// 	rows, err := db.Conn.Query(ctx,
// 		`SELECT id, name, category_id, price, description, image_url, parent_id, coverage_amount, age_limit, created_at, updated_at
// 		 FROM products
// 		 WHERE parent_id = $1
// 		 ORDER BY id`, parentID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var p models.Product
// 		if err := rows.Scan(&p.ID, &p.Name, &p.CategoryID, &p.Price, &p.Description, &p.ImageURL,
// 			&p.ParentID, &p.CoverageAmount, &p.AgeLimit, &p.CreatedAt, &p.UpdatedAt); err != nil {
// 			return nil, err
// 		}
// 		children = append(children, p)
// 	}

// 	return children, nil
// }

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
