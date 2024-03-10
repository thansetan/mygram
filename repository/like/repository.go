package likerepository

import (
	"context"
	"database/sql"
	"final-project/model"
	"fmt"
)

type likeRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *likeRepository {
	return &likeRepository{db}
}

func (r *likeRepository) Save(ctx context.Context, data model.Like) (model.Like, error) {
	var (
		like model.Like
		stmt = `
		INSERT INTO
			like_(user_id, photo_id)
			VALUES($1, $2)
		RETURNING
			id,
			user_id,
			photo_id,
			created_at
		`
	)

	row := r.db.QueryRowContext(ctx, stmt, data.UserID, data.PhotoID)
	if err := row.Err(); err != nil {
		return like, fmt.Errorf("likeRepository.Create: %w", err)
	}

	err := row.Scan(&like.ID, &like.UserID, &like.PhotoID, &like.CreatedAt)
	if err != nil {
		return like, fmt.Errorf("likeRepository.Create: %w", err)
	}

	return like, nil
}

func (r *likeRepository) FindByPhotoID(ctx context.Context, photoID uint64) ([]model.Like, error) {
	var (
		likes []model.Like
		stmt  = `
		SELECT
			l.id,
			l.user_id,
			l.photo_id,
			l.created_at,
			u.id,
			u.username,
			u.email
		FROM like_ l
		INNER JOIN user_ u ON l.user_id = u.id
		WHERE l.photo_id = $1
		ORDER BY l.created_at DESC
		`
	)

	rows, err := r.db.QueryContext(ctx, stmt, photoID)
	if err != nil {
		return nil, fmt.Errorf("likeRepository.FindByPhotoID: %w", err)
	}

	for rows.Next() {
		var like model.Like

		err := rows.Scan(&like.ID, &like.UserID, &like.PhotoID, &like.CreatedAt, &like.User.ID, &like.User.Username, &like.User.Email)
		if err != nil {
			return nil, fmt.Errorf("likeRepository.FindByPhotoID: %w", err)
		}

		likes = append(likes, like)
	}

	return likes, nil
}

func (r *likeRepository) Delete(ctx context.Context, data model.Like) error {
	var (
		stmt = `
		DELETE FROM like_
		WHERE user_id = $1 AND photo_id = $2
		`
	)

	res, err := r.db.ExecContext(ctx, stmt, data.UserID, data.PhotoID)
	if err != nil {
		return fmt.Errorf("likeRepository.Delete: %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("likeRepository.Delete: %w", err)
	} else if n == 0 {
		return fmt.Errorf("likeRepository.Delete: %w", sql.ErrNoRows)
	}

	return nil
}

func (r *likeRepository) FindByUserID(ctx context.Context, userID uint64) ([]model.Like, error) {
	var (
		likes []model.Like
		stmt  = `
		SELECT
			l.id,
			l.photo_id,
			l.user_id,
			l.created_at,
			p.id,
			p.title,
			p.caption,
			p.url,
			p.user_id
		FROM like_ l
		LEFT JOIN photo p ON l.photo_id=p.id
		WHERE l.user_id = $1
		ORDER BY l.created_at DESC
		`
	)

	rows, err := r.db.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, fmt.Errorf("likeRepository.FindByUserID: %w", err)
	}

	for rows.Next() {
		var like model.Like

		err := rows.Scan(&like.ID, &like.PhotoID, &like.UserID, &like.CreatedAt, &like.Photo.ID, &like.Photo.Title, &like.Photo.Caption, &like.Photo.URL, &like.Photo.UserID)
		if err != nil {
			return nil, fmt.Errorf("likeRepository.FindByUserID: %w", err)
		}

		likes = append(likes, like)
	}

	return likes, nil
}
