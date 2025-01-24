package database

import (
	"context"

	"github.com/aakash-tyagi/linmed/models"
)

func (db *Database) AddStationProduct(ctx context.Context, stationProduct models.StationProduct) (int, error) {

	var id int

	err := db.Conn.QueryRow(ctx,
		`INSERT INTO station_products (
		station_id,
		product_id,
		installation_date,
		expiry_date,
		inspection_date,
		child_product_1_id,
		child_product_1_qty,
		child_product_2_id,
		child_product_2_qty,
		created_at,
		updated_at,
		product_name,
		customer_name)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,$12,$13)
		RETURNING id;`,
		stationProduct.StationID, stationProduct.ProductID,
		stationProduct.InstalledDate, stationProduct.ExpiryDate,
		stationProduct.InspectionDate, stationProduct.ChildProduct1ID,
		stationProduct.ChildProduct1Qty, stationProduct.ChildProduct2ID,
		stationProduct.ChildProduct2Qty,
		stationProduct.CreatedAt,
		stationProduct.UpdatedAt,
		stationProduct.ProductName,
		stationProduct.CustomerName,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Database) GetStationProductById(ctx context.Context, id string) (models.StationProduct, error) {
	var stationProduct models.StationProduct

	err := db.Conn.QueryRow(ctx,
		`SELECT id, station_id, product_id, installation_date, expiry_date, inspection_date,
		child_product_1_id,
		child_product_1_qty,
		child_product_2_id,
		child_product_1_qty,
		created_at, updated_at
		FROM station_products
		WHERE id = $1;`,
		id,
	).Scan(&stationProduct.ID, &stationProduct.StationID, &stationProduct.ProductID, &stationProduct.InstalledDate, &stationProduct.ExpiryDate, &stationProduct.InspectionDate,
		&stationProduct.ChildProduct1ID,
		&stationProduct.ChildProduct1Qty,
		&stationProduct.ChildProduct2ID,
		&stationProduct.ChildProduct2Qty,
		&stationProduct.CreatedAt, &stationProduct.UpdatedAt)
	if err != nil {
		return stationProduct, err
	}

	return stationProduct, nil
}

func (db *Database) UpdateStationProduct(ctx context.Context, ID string, stationProduct models.StationProduct) error {

	_, err := db.Conn.Exec(ctx,
		`UPDATE station_products
		SET station_id = $1, product_id = $2, installation_date = $3, expiry_date = $4, inspection_date = $5
		child_product_1_id = $6,
		child_product_1_qty = $7,
		child_product_2_id = $8,
		child_product_1_qty = $9,
		WHERE id = $10;`,
		stationProduct.StationID, stationProduct.ProductID, stationProduct.InstalledDate, stationProduct.ExpiryDate, stationProduct.InspectionDate,
		stationProduct.ChildProduct1,
		stationProduct.ChildProduct1Qty,
		stationProduct.ChildProduct2ID,
		stationProduct.ChildProduct2Qty,
		ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteStationProduct(ctx context.Context, ID string) error {

	_, err := db.Conn.Exec(ctx,
		`DELETE FROM station_products
		WHERE id = $1;`,
		ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetStationProducts(ctx context.Context, page, limit int, customerId, stationId string) ([]models.StationProduct, int, error) {
	var stationProducts []models.StationProduct

	rows, err := db.Conn.Query(ctx,
		`SELECT id, station_id, product_id, installation_date, expiry_date, inspection_date,
		child_product_1_id,
		child_product_1_qty,
		child_product_2_id,
		child_product_1_qty,
		created_at, updated_at,
		product_name, customer_name
		FROM station_products
		WHERE station_id = $1
		ORDER BY id
		LIMIT $2 OFFSET $3;`,
		stationId, limit, (page-1)*limit,
	)

	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		var stationProduct models.StationProduct
		if err := rows.Scan(&stationProduct.ID, &stationProduct.StationID, &stationProduct.ProductID, &stationProduct.InstalledDate, &stationProduct.ExpiryDate, &stationProduct.InspectionDate,
			&stationProduct.ChildProduct1ID,
			&stationProduct.ChildProduct1Qty,
			&stationProduct.ChildProduct2ID,
			&stationProduct.ChildProduct2Qty,
			&stationProduct.CreatedAt, &stationProduct.UpdatedAt,
			&stationProduct.ProductName, &stationProduct.CustomerName); err != nil {
			return nil, 0, err
		}
		stationProducts = append(stationProducts, stationProduct)
	}

	return stationProducts, len(stationProducts), nil

}
