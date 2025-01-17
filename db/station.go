package database

import (
	"context"
	"fmt"

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

func (db *Database) GetStations(ctx context.Context, page, limit int, floorPlanID, customerID string) ([]models.Station, int, error) {
	var (
		stations []models.Station
		total    int
		offset   = (page - 1) * limit
	)

	// Query to get the total count of stations
	countQuery := `
		SELECT COUNT(*)
		FROM stations
		WHERE ($1::int IS NULL OR floor_plan_id = $1::int) 
		  AND ($2::int IS NULL OR customer_id = $2::int);
	`
	if err := db.Conn.QueryRow(ctx, countQuery, floorPlanID, customerID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to get total stations count: %w", err)
	}

	// Query to fetch stations with pagination
	stationsQuery := `
		SELECT id, name, description, customer_id, floor_plan_id, created_at, updated_at
		FROM stations
		WHERE ($1::int IS NULL OR floor_plan_id = $1::int) 
		  AND ($2::int IS NULL OR customer_id = $2::int)
		ORDER BY id
		LIMIT $3 OFFSET $4;
	`
	rows, err := db.Conn.Query(ctx, stationsQuery, floorPlanID, customerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch stations: %w", err)
	}
	defer rows.Close()

	// Scan rows into the stations slice
	for rows.Next() {
		var station models.Station
		if err := rows.Scan(&station.ID, &station.Name, &station.Description, &station.CustomerID, &station.FloorPlanID, &station.CreatedAt, &station.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan station row: %w", err)
		}
		stations = append(stations, station)
	}

	// Check for errors after iterating through rows
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating through rows: %w", err)
	}

	return stations, total, nil
}
