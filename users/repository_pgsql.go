package users

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresSQLRepository struct {
	db *pgx.Conn
}

func NewPostgresSQLRepository(db *pgx.Conn) *PostgresSQLRepository {
	return &PostgresSQLRepository{
		db: db,
	}
}

func (r *PostgresSQLRepository) Create(ctx context.Context, user User) (*User, error) {
	var newUser User
	err := r.db.QueryRow(ctx,
		"INSERT INTO users(id, name) values($1, $2) RETURNING id, created_at, updated_at, name",
		user.ID, user.Name).Scan(&newUser.ID, &newUser.CreatedAt, &newUser.UpdatedAt, &newUser.Name)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}
	return &newUser, nil
}

func (r *PostgresSQLRepository) GetById(ctx context.Context, ID uuid.UUID) (*User, error) {
	row := r.db.QueryRow(ctx,
		"SELECT id, name, created_at, updated_at from users where id = $1", ID)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotExist
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresSQLRepository) GetByName(ctx context.Context, name string) (*User, error) {
	row := r.db.QueryRow(ctx,
		"SELECT id, name, created_at, updated_at from users where name = $1", name)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotExist
		}
		return nil, err
	}
	return &user, nil
}
