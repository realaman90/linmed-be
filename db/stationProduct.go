package database

import (
	"context"

	"github.com/aakash-tyagi/linmed/models"
)

func (db *Database) AddStationProduct(ctx context.Context, stationProduct models.StationProduct) (int, error) {

	// return id of the stationProduct

	var id int

	err := db.Conn.QueryRow(ctx,
		`INSERT INTO station_products (
		station_id,
		product_id,
		installation_date,
		expiry_date,
		inspection_date,
		created_at,
		updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;`,
		stationProduct.StationID, stationProduct.ProductID, stationProduct.InstalledDate, stationProduct.ExpiryDate, stationProduct.InspectionDate, stationProduct.CreatedAt, stationProduct.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Database) GetStationProductById(ctx context.Context, id string) (models.StationProduct, error) {
	var stationProduct models.StationProduct

	err := db.Conn.QueryRow(ctx,
		`SELECT id, station_id, product_id, installation_date, expiry_date, inspection_date, created_at, updated_at
		FROM station_products
		WHERE id = $1;`,
		id,
	).Scan(&stationProduct.ID, &stationProduct.StationID, &stationProduct.ProductID, &stationProduct.InstalledDate, &stationProduct.ExpiryDate, &stationProduct.InspectionDate, &stationProduct.CreatedAt, &stationProduct.UpdatedAt)
	if err != nil {
		return stationProduct, err
	}

	return stationProduct, nil
}

func (db *Database) UpdateStationProduct(ctx context.Context, ID string, stationProduct models.StationProduct) error {

	_, err := db.Conn.Exec(ctx,
		`UPDATE station_products
		SET station_id = $1, product_id = $2, installation_date = $3, expiry_date = $4, inspection_date = $5
		WHERE id = $6;`,
		stationProduct.StationID, stationProduct.ProductID, stationProduct.InstalledDate, stationProduct.ExpiryDate, stationProduct.InspectionDate, ID,
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
