package database

import (
	"context"
	"fmt"

	"github.com/aakash-tyagi/linmed/models"
)

func (db *Database) GetAllNumbers(ctx context.Context) (models.Dashboard, error) {
	var dashboard models.Dashboard

	err := db.Conn.QueryRow(ctx,
		`SELECT
			(SELECT COUNT(*) FROM customers) as customer_count,
			(SELECT COUNT(*) FROM stations) as station_count,
			(SELECT COUNT(*) FROM products) as product_count,
			(SELECT COUNT(*) FROM station_products) as station_product_count;`,
	).Scan(&dashboard.CustomerCount, &dashboard.StationCount, &dashboard.ProductCount, &dashboard.StationProductCount)
	if err != nil {
		return dashboard, err
	}

	return dashboard, nil
}

func (db *Database) GetExpiringProducts(ctx context.Context, startDate, endDate, customerId string, page, limit int) ([]models.StationProduct, int, error) {
	// Calculate offset for pagination
	offset := (page - 1) * limit

	// Base SQL query for fetching expiring products
	query := `SELECT sp.id, sp.station_id, sp.product_id, sp.installation_date, sp.expiry_date, sp.inspection_date, sp.created_at, sp.updated_at,
					 sp.child_product_1_id, sp.child_product_1_qty, sp.child_product_2_id, sp.child_product_2_qty,
					 c.name AS customer_name, p.name AS product_name
			  FROM station_products sp
			  JOIN customers c ON sp.station_id = c.id
			  JOIN products p ON sp.product_id = p.id
			  WHERE sp.expiry_date BETWEEN $1 AND $2`

	var args []interface{}
	args = append(args, startDate, endDate)

	// Add customer filter if customerId is provided
	if customerId != "" {
		query += " AND c.id = $3"
		args = append(args, customerId)
	}

	query += fmt.Sprintf(" ORDER BY sp.expiry_date ASC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	// Use db.Conn instead of conn
	rows, err := db.Conn.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close() // Ensure the rows are closed when done

	var expiringProducts []models.StationProduct

	// Iterate through the result set
	for rows.Next() {
		var product models.StationProduct
		if err := rows.Scan(
			&product.ID,
			&product.StationID,
			&product.ProductID,
			&product.InstalledDate,
			&product.ExpiryDate,
			&product.InspectionDate,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.ChildProduct1ID,
			&product.ChildProduct1Qty,
			&product.ChildProduct2ID,
			&product.ChildProduct2Qty,
			&product.CustomerName,
			&product.ProductName,
		); err != nil {
			return nil, 0, err
		}
		expiringProducts = append(expiringProducts, product)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// Total count query
	countQuery := `SELECT COUNT(*) FROM station_products sp
				   JOIN customers c ON sp.station_id = c.id
				   WHERE sp.expiry_date BETWEEN $1 AND $2`

	if customerId != "" {
		countQuery += " AND c.id = $3"
	}

	var totalCount int
	if customerId != "" {
		err = db.Conn.QueryRow(ctx, countQuery, startDate, endDate, customerId).Scan(&totalCount)
	} else {
		err = db.Conn.QueryRow(ctx, countQuery, startDate, endDate).Scan(&totalCount)
	}
	if err != nil {
		return nil, 0, err
	}

	// Return the list of expiring products and the total count
	return expiringProducts, totalCount, nil
}

func (db *Database) GetInspectionTasks(ctx context.Context, startDate, endDate string, page, limit int) ([]models.StationProduct, int, error) {
	var total int
	var tasks []models.StationProduct

	// Get total count
	err := db.Conn.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM station_products 
		WHERE inspection_date BETWEEN $1 AND $2`,
		startDate, endDate).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count query failed: %v", err)
	}

	// Get paginated results
	query := `SELECT id, station_id, product_id, installation_date, expiry_date, inspection_date, created_at, updated_at 
			 FROM station_products 
			 WHERE inspection_date BETWEEN $1 AND $2
			 ORDER BY inspection_date ASC
			 LIMIT $3 OFFSET $4`

	offset := (page - 1) * limit
	rows, err := db.Conn.Query(ctx, query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task models.StationProduct
		err := rows.Scan(
			&task.ID,
			&task.StationID,
			&task.ProductID,
			&task.InstalledDate,
			&task.ExpiryDate,
			&task.InspectionDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	return tasks, total, nil
}
