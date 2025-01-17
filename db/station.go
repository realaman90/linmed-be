package database

import (
	"context"

	"github.com/aakash-tyagi/linmed/models"
)

func (db *Database) AddStation(ctx context.Context, station models.Station) (int, error) {

	var id int

	err := db.Conn.QueryRow(ctx,
		`INSERT INTO stations (
		name,
		description,
		customer_id,
		floor_plan_id,
		created_at,
		updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;`,
		station.Name, station.Description, station.CustomerID, station.FloorPlanID, station.CreatedAt, station.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Database) GetStation(ctx context.Context, id string) (models.Station, error) {
	var station models.Station

	err := db.Conn.QueryRow(ctx,
		`SELECT id, name, description, customer_id, floor_plan_id, created_at, updated_at
		FROM stations
		WHERE id = $1;`,
		id,
	).Scan(&station.ID, &station.Name, &station.Description, &station.CustomerID, &station.FloorPlanID, &station.CreatedAt, &station.UpdatedAt)
	if err != nil {
		return station, err
	}

	return station, nil
}

func (db *Database) UpdateStation(ctx context.Context, ID string, station models.Station) error {

	_, err := db.Conn.Exec(ctx,
		`UPDATE stations
		SET name = $1, description = $2
		WHERE id = $3;`,
		station.Name, station.Description, ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteStation(ctx context.Context, ID string) error {

	_, err := db.Conn.Exec(ctx,
		`DELETE FROM stations
		WHERE id = $1;`,
		ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetStations(ctx context.Context, page, limit int, floorPlanId, customerId string) ([]models.Station, int, error) {

	// var stations

	return nil, 0, nil
}
