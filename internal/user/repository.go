package user

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user User) error {
	query := `
		INSERT INTO users
			(id, email, name, password)
		VALUES
			(:id, :email, :name, :password)
	`

	updatedQuery, args, err := sqlx.Named(query, user)
	if err != nil {
		return err
	}

	// since we won't be using the returned data, leave it blank
	_, err = r.db.ExecContext(ctx, sqlx.Rebind(sqlx.DOLLAR, updatedQuery), args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (User, error) {
	var result User

	query := `
		SELECT
			id,
			email,
			name,
			password
		FROM
			users
		WHERE
			email = $1
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &result, query, email)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id string) (User, error) {
	var result User

	query := `
		SELECT
			id,
			email,
			name,
			password
		FROM
			users
		WHERE
			id = $1
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &result, query, id)
	if err != nil {
		return result, err
	}

	return result, nil
}
