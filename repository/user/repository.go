package userrepository

import (
	"context"
	"database/sql"
	"final-project/model"
	"fmt"
)

type userRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Save(ctx context.Context, data model.User) (model.User, error) {
	var (
		user model.User
		stmt = `
		INSERT INTO
			user_(username, email, password, age)
			VALUES($1, $2, $3, $4)
		RETURNING
			id,
			username,
			email,
			age
		`
	)

	row := r.db.QueryRowContext(ctx, stmt, data.Username, data.Email, data.Password, data.Age)
	if err := row.Err(); err != nil {
		return user, fmt.Errorf("userRepository.Create: %w", err)
	}

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Age)
	if err != nil {
		return user, fmt.Errorf("userRepository.Create: %w", err)
	}

	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	var (
		user model.User
		stmt = `
		SELECT
			id,
			password
		FROM user_
		WHERE email=$1
		`
	)

	row := r.db.QueryRowContext(ctx, stmt, email)
	if err := row.Err(); err != nil {
		return user, fmt.Errorf("userRepository.FindByEmail: %w", err)
	}

	err := row.Scan(&user.ID, &user.Password)
	if err != nil {
		return user, fmt.Errorf("userRepository.FindByEmail: %w", err)
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, data model.User) (model.User, error) {
	var (
		user model.User
		stmt = `
		UPDATE
			user_
		SET
			email=$1,
			username=$2,
			updated_at=NOW()
		WHERE id=$3 AND updated_at=$4
		RETURNING
			id,
			username,
			email,
			age,
			updated_at
		`
	)

	row := r.db.QueryRowContext(ctx, stmt, data.Email, data.Username, data.ID, data.UpdatedAt)
	if err := row.Err(); err != nil {
		return user, fmt.Errorf("userRepository.Update: %w", err)
	}

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Age, &user.UpdatedAt)
	if err != nil {
		return user, fmt.Errorf("userRepository.Update: %w", err)
	}

	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, userID uint64) error {
	var (
		stmt = `
		DELETE FROM
			user_
		WHERE id=$1
		`
	)
	res, err := r.db.ExecContext(ctx, stmt, userID)
	if err != nil {
		return fmt.Errorf("userRepository.Delete: %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("userRepository.Delete: %w", err)
	} else if n == 0 {
		return fmt.Errorf("userRepository.Delete: %w", sql.ErrNoRows)
	}

	return nil
}

func (r *userRepository) FindByID(ctx context.Context, userID uint64) (model.User, error) {
	var (
		user model.User
		stmt = `
		SELECT
			id,
			username,
			email,
			age,
			created_at,
			updated_at
		FROM user_
		WHERE id=$1
		`
	)

	row := r.db.QueryRowContext(ctx, stmt, userID)
	if err := row.Err(); err != nil {
		return user, fmt.Errorf("userRepository.FindByID: %w", err)
	}

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Age, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, fmt.Errorf("userRepository.FindByID: %w", err)
	}

	return user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (model.User, error) {
	var (
		user model.User
		stmt = `
		SELECT
			id,
			email,
			age
		FROM user_
		WHERE username=$1
		`
	)

	row := r.db.QueryRowContext(ctx, stmt, username)
	if err := row.Err(); err != nil {
		return user, fmt.Errorf("userRepository.FindByUsername: %w", err)
	}

	err := row.Scan(&user.ID, &user.Email, &user.Age)
	if err != nil {
		return user, fmt.Errorf("userRepository.FindByUsername: %w", err)
	}

	return user, nil
}
