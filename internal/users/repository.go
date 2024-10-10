package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) (*sqliteRepository, error) {
	return &sqliteRepository{db: db}, nil
}

func (r *sqliteRepository) Create(ctx context.Context, user User) (User, error) {
	addressJSON, err := json.Marshal(user.Address)
	if err != nil {
		return User{}, fmt.Errorf("failed to marshal address to JSON: %v", err)
	}

	query := `INSERT INTO users (name, email, password, address) VALUES (?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, addressJSON)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return User{}, ErrMailAlreadyExists
			}
		}
		return user, fmt.Errorf("failed to create user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return user, fmt.Errorf("failed to retrieve last insert ID: %v", err)
	}

	user.ID = id
	return user, nil
}

func (r *sqliteRepository) Update(ctx context.Context, email, password string) error {
	query := `UPDATE users SET password = ? WHERE email = ?`
	res, err := r.db.ExecContext(ctx, query, password, email)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *sqliteRepository) GetByEMail(ctx context.Context, email string) (User, error) {
	query := `SELECT id, name, email, password FROM users WHERE email = ?`
	row := r.db.QueryRowContext(ctx, query, email)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrNotFound
	} else if err != nil {
		return user, fmt.Errorf("failed to retrieve user by email: %v", err)
	}
	return user, nil
}

func (r *sqliteRepository) GetAll(ctx context.Context, limit, offset int) ([]User, error) {
	query := `SELECT id, name, email, password FROM users LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}
	defer rows.Close()

	var usersList []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		usersList = append(usersList, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return usersList, nil
}

func (r *sqliteRepository) Close() error {
	return r.db.Close()
}
