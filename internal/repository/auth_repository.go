package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/aldngrha/ecommerce-be/internal/entity"
)

type IAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	InsertUser(ctx context.Context, user *entity.User) error
}

type authRepository struct {
	db *sql.DB
}

func (ar *authRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := ar.db.QueryRowContext(ctx, "SELECT id, email, password, full_name, role_code FROM users WHERE email = $1", email)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var user entity.User

	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.RoleCode,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no user found, return nil and no error
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (as *authRepository) InsertUser(ctx context.Context, user *entity.User) error {
	_, err := as.db.ExecContext(ctx,
		"INSERT INTO users (id, full_name, email, role_code, password, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, is_deleted) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		user.Id,
		user.FullName,
		user.Email,
		user.RoleCode,
		user.Password,
		user.CreatedAt,
		user.CreatedBy,
		user.UpdatedAt,
		user.UpdatedBy,
		user.DeletedAt,
		user.DeletedBy,
		user.IsDeleted,
	)
	if err != nil {
		return err
	}

	return nil
}

func NewAuthRepository(db *sql.DB) IAuthRepository {
	return &authRepository{
		db: db,
	}
}
