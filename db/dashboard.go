package database

import (
	"context"

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

func (db *Database) GetUpcomingStaionWithTask(ctx context.Context, period string) ([]models.StationProduct, error) {

	var staionProduct []models.StationProduct

	rows, err := db.Conn.Query(ctx,
		`SELECT
			sp.id,
			sp.station_id,
			sp.product_id,
			sp.installation_date,
			sp.expiry_date,
			sp.inspection_date,
			sp.created_at,
			sp.updated_at
		FROM station_products sp
		WHERE sp.inspection_date BETWEEN NOW() AND NOW() + INTERVAL $1;`,
		period,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stationProduct models.StationProduct
		if err := rows.Scan(&stationProduct.ID, &stationProduct.StationID, &stationProduct.ProductID, &stationProduct.InstalledDate, &stationProduct.ExpiryDate, &stationProduct.InspectionDate, &stationProduct.CreatedAt, &stationProduct.UpdatedAt); err != nil {
			return nil, err
		}
		staionProduct = append(staionProduct, stationProduct)
	}

	return staionProduct, nil
}
