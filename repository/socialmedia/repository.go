package socialmediarepository

import (
	"context"
	"database/sql"
	"final-project/model"
	"fmt"
)

type socialMediaRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *socialMediaRepository {
	return &socialMediaRepository{db}
}

func (r *socialMediaRepository) Save(ctx context.Context, data model.SocialMedia) (model.SocialMedia, error) {
	var (
		socialMedia model.SocialMedia
		stmt        = `
		INSERT INTO
			social_media(user_id, name, url)
			VALUES($1, $2, $3)
		RETURNING
			id,
			name,
			url,
			user_id,
			created_at
		`
	)

	row := r.db.QueryRowContext(ctx, stmt, data.UserID, data.Name, data.URL)
	if err := row.Err(); err != nil {
		return socialMedia, fmt.Errorf("socialMediaRepository.Create: %w", err)
	}

	err := row.Scan(&socialMedia.ID, &socialMedia.Name, &socialMedia.URL, &socialMedia.UserID, &socialMedia.CreatedAt)
	if err != nil {
		return socialMedia, fmt.Errorf("socialMediaRepository.Create: %w", err)
	}

	return socialMedia, nil
}

func (r *socialMediaRepository) FindAll(ctx context.Context) ([]model.SocialMedia, error) {
	var (
		socialMedias []model.SocialMedia
		stmt         = `
		SELECT
			s.id,
			s.name,
			s.url,
			s.user_id,
			s.created_at,
			s.updated_at,
			u.username,
			u.email
		FROM social_media s
		INNER JOIN user_ u ON s.user_id = u.id
		ORDER BY s.created_at DESC
		`
	)

	rows, err := r.db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("socialMediaRepository.FindAll: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var socialMedia model.SocialMedia

		err := rows.Scan(&socialMedia.ID, &socialMedia.Name, &socialMedia.URL, &socialMedia.UserID, &socialMedia.CreatedAt, &socialMedia.UpdatedAt, &socialMedia.User.Username, &socialMedia.User.Email)
		if err != nil {
			return nil, fmt.Errorf("socialMediaRepository.FindAll: %w", err)
		}

		socialMedias = append(socialMedias, socialMedia)
	}

	return socialMedias, nil
}

func (r *socialMediaRepository) Update(ctx context.Context, data model.SocialMedia) (model.SocialMedia, error) {
	var (
		socialMedia model.SocialMedia
		stmt        = `
		UPDATE
			social_media
		SET
			name=$1,
			url=$2,
			updated_at=NOW()
		WHERE id=$3 AND updated_at=$4
		RETURNING
			id,
			name,
			url,
			user_id,
			updated_at
	`
	)

	row := r.db.QueryRowContext(ctx, stmt, data.Name, data.URL, data.ID, data.UpdatedAt)
	if err := row.Err(); err != nil {
		return socialMedia, fmt.Errorf("socialMediaRepository.Update: %w", err)
	}

	err := row.Scan(&socialMedia.ID, &socialMedia.Name, &socialMedia.URL, &socialMedia.UserID, &socialMedia.UpdatedAt)
	if err != nil {
		return socialMedia, fmt.Errorf("socialMediaRepository.Update: %w", err)
	}

	return socialMedia, nil
}

func (r *socialMediaRepository) Delete(ctx context.Context, data model.SocialMedia) error {
	var (
		stmt = `
		DELETE FROM
			social_media
		WHERE id=$1 AND user_id=$2
		`
	)

	res, err := r.db.ExecContext(ctx, stmt, data.ID, data.UserID)
	if err != nil {
		return fmt.Errorf("socialMediaRepository.Delete: %w", err)
	}

	if n, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("socialMediaRepository.Delete: %w", err)
	} else if n == 0 {
		return fmt.Errorf("socialMediaRepository.Delete: %w", sql.ErrNoRows)
	}

	return nil
}

func (r *socialMediaRepository) FindByID(ctx context.Context, id uint64) (model.SocialMedia, error) {
	var (
		socialMedia model.SocialMedia
		stmt        = `
		SELECT
			s.id,
			s.name,
			s.url,
			s.user_id,
			s.created_at,
			s.updated_at,
			u.username,
			u.email
		FROM social_media s
		INNER JOIN user_ u ON s.user_id = u.id
		WHERE s.id=$1
		`
	)

	row := r.db.QueryRowContext(ctx, stmt, id)
	if err := row.Err(); err != nil {
		return socialMedia, fmt.Errorf("socialMediaRepository.FindByID: %w", err)
	}

	err := row.Scan(&socialMedia.ID, &socialMedia.Name, &socialMedia.URL, &socialMedia.UserID, &socialMedia.CreatedAt, &socialMedia.UpdatedAt, &socialMedia.User.Username, &socialMedia.User.Email)
	if err != nil {
		return socialMedia, fmt.Errorf("socialMediaRepository.FindByID: %w", err)
	}

	return socialMedia, nil
}

func (r *socialMediaRepository) FindByUserID(ctx context.Context, userID uint64) ([]model.SocialMedia, error) {
	var (
		socialMedias []model.SocialMedia
		stmt         = `
		SELECT
			s.id,
			s.name,
			s.url,
			s.user_id,
			s.created_at,
			s.updated_at
		FROM social_media s
		WHERE s.user_id=$1
		ORDER BY s.created_at DESC
		`
	)

	rows, err := r.db.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, fmt.Errorf("socialMediaRepository.FindByUserID: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var socialMedia model.SocialMedia

		err := rows.Scan(&socialMedia.ID, &socialMedia.Name, &socialMedia.URL, &socialMedia.UserID, &socialMedia.CreatedAt, &socialMedia.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("socialMediaRepository.FindByUserID: %w", err)
		}

		socialMedias = append(socialMedias, socialMedia)
	}

	return socialMedias, nil
}
