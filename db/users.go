package database

import (
	"context"

	"github.com/aakash-tyagi/linmed/models"
)

func (db *Database) AddUser(ctx context.Context, user models.User) (uint, error) {

	var id uint

	err := db.Conn.QueryRow(ctx,
		`INSERT INTO users (
		username,
		email,
		first_name,
		last_name, 
		phone_number,
		is_active,
		last_login,
		password_hash,
		role,
		created_at,
		updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id;`,
		user.Username, user.Email, user.FirstName, user.LastName, user.PhoneNumber, user.IsActive, user.LastLogin, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *Database) GetUser(ctx context.Context, id string) (models.User, error) {
	var user models.User

	err := db.Conn.QueryRow(ctx,
		`SELECT id, username, email, first_name, last_name, created_at, updated_at, phone_number, is_active, last_login, role
		FROM users
		WHERE id = $1;`,
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt, &user.PhoneNumber, &user.IsActive, &user.LastLogin, &user.Role)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (db *Database) UpdateUser(ctx context.Context, ID string, user models.User) error {

	_, err := db.Conn.Exec(ctx,
		`UPDATE users
		SET username = $1, email = $2, first_name = $3, last_name = $4
		WHERE id = $5;`,
		user.Username, user.Email, user.FirstName, user.LastName, ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetUsers(ctx context.Context, page, limit int) ([]models.User, int, error) {
	var users []models.User

	rows, err := db.Conn.Query(ctx,
		`SELECT id, username, email, first_name, last_name, created_at, updated_at, phone_number, is_active, last_login, role
		FROM users
		LIMIT $1 OFFSET $2;`,
		limit, page,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt, &user.PhoneNumber, &user.IsActive, &user.LastLogin, &user.Role); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, len(users), nil
}
