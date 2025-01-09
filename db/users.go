package database

import (
	"context"

	"github.com/aakash-tyagi/linmed/models"
)

func (db *Database) AddUser(ctx context.Context, user models.User) error {

	_, err := db.Conn.Exec(ctx,
		`INSERT INTO users (username, email, password_hash, first_name, last_name)
		VALUES ($1, $2, $3, $4, $5);`,
		user.Username, user.Email, user.PasswordHash, user.FirstName, user.LastName,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetUser(ctx context.Context, id string) (models.User, error) {
	var user models.User

	err := db.Conn.QueryRow(ctx,
		`SELECT id, username, email, first_name, last_name, created_at, updated_at
		FROM users
		WHERE id = $1;`,
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
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

func (db *Database) GetUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	rows, err := db.Conn.Query(ctx,
		`SELECT id, username, email, first_name, last_name, created_at, updated_at
		FROM users;`,
	)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}
