package mailvalidation

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteMailValidationRepo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &sqliteMailValidationRepo{db: db}
}

func (repo *sqliteMailValidationRepo) CreateOrUpdate(ctx context.Context, mailValidation MailValidation) error {
	_, err := repo.db.ExecContext(ctx, `
		INSERT INTO mail_validations (email, code, expired_at) 
		VALUES (?, ?, ?)
		ON CONFLICT(email) DO UPDATE SET 
			code = excluded.code,
			expired_at = excluded.expired_at;
	`, mailValidation.Email, mailValidation.Code, mailValidation.ExpiredAt)
	if err != nil {
		return fmt.Errorf("failed to create or update mail validation: %w", err)
	}
	return nil
}

func (repo *sqliteMailValidationRepo) GetByEmail(ctx context.Context, email string) (MailValidation, error) {
	var mailValidation MailValidation
	err := repo.db.QueryRowContext(ctx, "SELECT email, code, expired_at FROM mail_validations WHERE email = ?", email).Scan(&mailValidation.Email, &mailValidation.Code, &mailValidation.ExpiredAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return MailValidation{}, ErrRecordNotFound
		}
		return MailValidation{}, fmt.Errorf("failed to get mail validation: %w", err)
	}
	return mailValidation, nil
}

func (repo *sqliteMailValidationRepo) Delete(ctx context.Context, email string) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM mail_validations WHERE email = ?", email)
	if err != nil {
		return fmt.Errorf("failed to delete mail validation: %w", err)
	}
	return nil
}
