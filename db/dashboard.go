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

func (db *Database) GetExpiringProducts(ctx context.Context, startDate, endDate string, page, limit int) ([]models.ExpiringProductResponse, int, error) {
	var total int

	res := []models.ExpiringProductResponse{}

	// Get total count
	err := db.Conn.QueryRow(ctx, `
		SELECT COUNT(*) 
		FROM station_products 
		WHERE expiry_date BETWEEN $1 AND $2`,
		startDate, endDate).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count query failed: %v", err)
	}

	// Get paginated results with customer and product details
	query := `SELECT sp.id, sp.station_id, sp.product_id, sp.installation_date, sp.expiry_date, sp.inspection_date, sp.created_at, sp.updated_at,
					 c.name AS customer_name, p.name AS product_name
			  FROM station_products sp
			  JOIN customers c ON sp.station_id = c.id
			  JOIN products p ON sp.product_id = p.id
			  WHERE sp.expiry_date BETWEEN $1 AND $2
			  ORDER BY sp.expiry_date ASC
			  LIMIT $3 OFFSET $4`

	offset := (page - 1) * limit
	rows, err := db.Conn.Query(ctx, query, startDate, endDate, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task models.ExpiringProductResponse
		err := rows.Scan(
			&task.ID,
			&task.StationID,
			&task.ProductID,
			&task.InstalledDate,
			&task.ExpiryDate,
			&task.InspectionDate,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.CustomerName,
			&task.ProductName,
		)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, task)
	}

	return res, total, nil
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
