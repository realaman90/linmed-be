package database

import (
	"context"

	"github.com/aakash-tyagi/linmed/models"
)

func (s *Database) AddFloorPlan(ctx context.Context, floorPlan models.FloorPlan) (int, error) {

	var id int

	err := s.Conn.QueryRow(ctx,
		`INSERT INTO floor_plans (
		name,
		layout,
		customer_id,
		created_at,
		updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;`,
		floorPlan.Name, floorPlan.Layout, floorPlan.CustomerID, floorPlan.CreatedAt, floorPlan.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil

}
func (s *Database) GetFloorPlan(ctx context.Context, id string) (models.FloorPlan, error) {
	var floorPlan models.FloorPlan

	err := s.Conn.QueryRow(ctx,
		`SELECT id, name, layout, created_at, updated_at
		FROM floor_plans
		WHERE id = $1;`,
		id,
	).Scan(&floorPlan.ID, &floorPlan.Name, &floorPlan.Layout, &floorPlan.CreatedAt, &floorPlan.UpdatedAt)
	if err != nil {
		return floorPlan, err
	}

	return floorPlan, nil

}

func (s *Database) UpdateFloorPlan(ctx context.Context, ID string, floorPlan models.FloorPlan) error {

	_, err := s.Conn.Exec(ctx,
		`UPDATE floor_plans
		SET name = $1, layout = $3
		WHERE id = $4;`,
		floorPlan.Name, floorPlan.Layout, ID,
	)
	if err != nil {
		return err
	}

	// update floor plan in customer

	_, err = s.Conn.Exec(ctx,
		`UPDATE customers
		SET floor_plans = jsonb_set(
			floor_plans::jsonb,
			ARRAY[CAST($1 AS TEXT)],
			$2::jsonb,
			true
		)
		WHERE id = $3;`,
		ID, floorPlan.CustomerID,
	)

	return nil
}

func (s *Database) DeleteFloorPlan(ctx context.Context, id string) error {

	_, err := s.Conn.Exec(ctx,
		`DELETE FROM floor_plans
		WHERE id = $1;`,
		id,
	)
	if err != nil {
		return err
	}

	// remove floor plan from customer

	_, err = s.Conn.Exec(ctx,
		`UPDATE customers
		SET floor_plans = array_remove(floor_plans, $1)
		WHERE id = $2;`,
		id, id,
	)
	if err != nil {
		return err
	}

	return nil

}

func (s *Database) GetFloorPlans(ctx context.Context, customerId uint, page, limit int) ([]models.FloorPlan, int, error) {
	var floorPlans []models.FloorPlan

	rows, err := s.Conn.Query(ctx,
		`SELECT id, name, layout, created_at, updated_at
		FROM floor_plans
		WHERE customer_id = $1
		ORDER BY id
		LIMIT $2 OFFSET $3;`,
		customerId, limit, (page-1)*limit,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var floorPlan models.FloorPlan
		if err := rows.Scan(&floorPlan.ID, &floorPlan.Name, &floorPlan.Layout, &floorPlan.CreatedAt, &floorPlan.UpdatedAt); err != nil {
			return nil, 0, err
		}
		floorPlans = append(floorPlans, floorPlan)
	}

	return floorPlans, len(floorPlans), nil

}
