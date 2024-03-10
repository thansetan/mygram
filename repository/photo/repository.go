package photorepository

import (
	"context"
	"database/sql"
	"final-project/model"
	"fmt"
)

type photoRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *photoRepository {
	return &photoRepository{db}
}

func (r *photoRepository) Save(ctx context.Context, data model.Photo) (model.Photo, error) {
	var (
		photo model.Photo
		stmt  = `
		INSERT INTO 
			photo(title, caption, url, user_id)
			VALUES($1, $2, $3, $4)
		RETURNING
			id, 
			title, 
			caption, 
			url, 
			user_id, 
			created_at
		`
	)

	row := r.db.QueryRowContext(ctx, stmt, data.Title, data.Caption, data.URL, data.UserID)
	if err := row.Err(); err != nil {
		return photo, fmt.Errorf("photoRepository.Create: %w", err)
	}

	err := row.Scan(&photo.ID, &photo.Title, &photo.Caption, &photo.URL, &photo.UserID, &photo.CreatedAt)
	if err != nil {
		return photo, fmt.Errorf("photoRepository.Create: %w", err)
	}

	return photo, nil
}

func (r *photoRepository) FindAll(ctx context.Context) ([]model.Photo, error) {
	var (
		photos []model.Photo
		stmt   = `
		SELECT
			p.id,
			p.title,
			p.caption,
			p.url,
			p.user_id,
			p.created_at,
			p.updated_at,
			u.email,
			u.username
		FROM photo p
		INNER JOIN user_ u ON p.user_id=u.id
		ORDER BY p.created_at DESC
		`
	)

	rows, err := r.db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("photoRepository.FindAll: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var photo model.Photo

		err := rows.Scan(&photo.ID, &photo.Title, &photo.Caption, &photo.URL, &photo.UserID, &photo.CreatedAt, &photo.UpdatedAt, &photo.User.Email, &photo.User.Username)
		if err != nil {
			return nil, fmt.Errorf("photoRepository.FindAll: %w", err)
		}

		photos = append(photos, photo)
	}

	return photos, nil
}

func (r *photoRepository) Update(ctx context.Context, data model.Photo) (model.Photo, error) {
	var (
		photo model.Photo
		stmt  = `
		UPDATE 
			photo
		SET 
			title=$1,
			caption=$2,
			url=$3,
			updated_at=NOW()
		WHERE id=$4 AND updated_at=$5
		RETURNING 
			id, 
			title, 
			caption, 
			url, 
			user_id, 
			updated_at
		`
	)

	row := r.db.QueryRowContext(ctx, stmt, data.Title, data.Caption, data.URL, data.ID, data.UpdatedAt)
	if err := row.Err(); err != nil {
		return photo, fmt.Errorf("photoRepository.Update: %w", err)
	}

	err := row.Scan(&photo.ID, &photo.Title, &photo.Caption, &photo.URL, &photo.UserID, &photo.UpdatedAt)
	if err != nil {
		return photo, fmt.Errorf("photoRepository.Update: %w", err)
	}

	return photo, nil
}

func (r *photoRepository) Delete(ctx context.Context, data model.Photo) error {
	var (
		stmt = `
		DELETE FROM
			photo
		WHERE id=$1 AND user_id=$2
		`
	)

	res, err := r.db.ExecContext(ctx, stmt, data.ID, data.UserID)
	if err != nil {
		return fmt.Errorf("photoRepository.Delete: %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("photoRepository.Delete: %w", err)
	} else if n == 0 {
		return fmt.Errorf("photoRepository.Delete: %w", sql.ErrNoRows)
	}

	return nil
}

func (r *photoRepository) FindByID(ctx context.Context, id uint64) (model.Photo, error) {
	var (
		photo model.Photo
		stmt  = `
		SELECT
			p.id,
			p.title,
			p.caption,
			p.url,
			p.user_id,
			p.created_at,
			p.updated_at,
			u.email,
			u.username
		FROM photo p
		INNER JOIN user_ u ON p.user_id=u.id
		WHERE p.id=$1`
	)

	row := r.db.QueryRowContext(ctx, stmt, id)
	if err := row.Err(); err != nil {
		return photo, fmt.Errorf("photoRepository.FindByID: %w", err)
	}

	err := row.Scan(&photo.ID, &photo.Title, &photo.Caption, &photo.URL, &photo.UserID, &photo.CreatedAt, &photo.UpdatedAt, &photo.User.Email, &photo.User.Username)
	if err != nil {
		return photo, fmt.Errorf("photoRepository.FindByID: %w", err)
	}

	return photo, nil
}

func (r *photoRepository) FindByUserID(ctx context.Context, userID uint64) ([]model.Photo, error) {
	var (
		photos []model.Photo
		stmt   = `
		SELECT
			p.id,
			p.title,
			p.caption,
			p.url,
			p.user_id,
			p.created_at,
			p.updated_at,
			u.email,
			u.username
		FROM photo p
		INNER JOIN user_ u ON p.user_id=u.id
		WHERE p.user_id=$1
		ORDER BY p.created_at DESC
		`
	)

	rows, err := r.db.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, fmt.Errorf("photoRepository.FindByUserID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var photo model.Photo

		err := rows.Scan(&photo.ID, &photo.Title, &photo.Caption, &photo.URL, &photo.UserID, &photo.CreatedAt, &photo.UpdatedAt, &photo.User.Email, &photo.User.Username)
		if err != nil {
			return nil, fmt.Errorf("photoRepository.FindByUserID: %w", err)
		}

		photos = append(photos, photo)
	}

	return photos, nil
}

func (r *photoRepository) FindByUsername(ctx context.Context, username string) ([]model.Photo, error) {
	var (
		photos []model.Photo
		stmt   = `
		SELECT
			p.id,
			p.title,
			p.caption,
			p.url,
			p.user_id,
			p.created_at,
			p.updated_at,
			u.email,
			u.username
		FROM photo p
		INNER JOIN user_ u ON p.user_id=u.id
		WHERE u.username=$1
		ORDER BY p.created_at DESC
		`
	)

	rows, err := r.db.QueryContext(ctx, stmt, username)
	if err != nil {
		return nil, fmt.Errorf("photoRepository.FindByUsername: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var photo model.Photo

		err := rows.Scan(&photo.ID, &photo.Title, &photo.Caption, &photo.URL, &photo.UserID, &photo.CreatedAt, &photo.UpdatedAt, &photo.User.Email, &photo.User.Username)
		if err != nil {
			return nil, fmt.Errorf("photoRepository.FindByUsername: %w", err)
		}

		photos = append(photos, photo)
	}

	return photos, nil
}
