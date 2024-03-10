package photoservice

import (
	"context"
	"database/sql"
	"errors"
	"final-project/dto"
	"final-project/helper"
	"final-project/model"
	"final-project/repository"
	"log/slog"
	"net/http"
)

type photoService struct {
	userRepo  repository.UserRepository
	photoRepo repository.PhotoRepository
	logger    *slog.Logger
}

func New(userRepo repository.UserRepository, photoRepo repository.PhotoRepository, logger *slog.Logger) *photoService {
	return &photoService{userRepo, photoRepo, logger}
}

func (s *photoService) Create(ctx context.Context, data dto.PhotoRequest) (dto.PhotoCreateResponse, error) {
	var (
		resp dto.PhotoCreateResponse
		err  error
	)

	userID, ok := ctx.Value(helper.UserIDKey).(float64)
	if !ok {
		s.logger.ErrorContext(ctx, "ctx.Value(helper.UserIDKey).(float64): userID is not float64")
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	photo := model.Photo{
		Title:  data.Title,
		URL:    data.URL,
		UserID: uint64(userID),
	}

	if data.Caption != "" {
		photo.Caption.String = data.Caption
		photo.Caption.Valid = true
	}

	photo, err = s.photoRepo.Save(ctx, photo)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = dto.PhotoCreateResponse{
		ID:        photo.ID,
		Title:     photo.Title,
		URL:       photo.URL,
		UserID:    photo.UserID,
		CreatedAt: photo.CreatedAt,
	}

	if photo.Caption.Valid {
		resp.Caption = photo.Caption.String
	}

	return resp, nil
}

func (s *photoService) GetAll(ctx context.Context) ([]dto.PhotoResponse, error) {
	var resp []dto.PhotoResponse

	photos, err := s.photoRepo.FindAll(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = make([]dto.PhotoResponse, 0, len(photos))

	for _, photo := range photos {
		item := dto.PhotoResponse{
			ID:        photo.ID,
			Title:     photo.Title,
			URL:       photo.URL,
			UserID:    photo.UserID,
			CreatedAt: photo.CreatedAt,
			UpdatedAt: photo.UpdatedAt,
			User: dto.User{
				ID:       photo.UserID,
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
		}

		if photo.Caption.Valid {
			item.Caption = photo.Caption.String
		}

		resp = append(resp, item)
	}

	return resp, nil
}

func (s *photoService) Update(ctx context.Context, id uint64, data dto.PhotoRequest) (resp dto.PhotoUpdateResponse, err error) {
	userID, ok := ctx.Value(helper.UserIDKey).(float64)
	if !ok {
		s.logger.ErrorContext(ctx, "ctx.Value(helper.UserIDKey).(float64): userID is not float64")
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	photo, err := s.photoRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrNotAllowed, http.StatusForbidden)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	if photo.UserID != uint64(userID) {
		s.logger.ErrorContext(ctx, "photo.UserID != uint64(userID): user is not the owner of the photo")
		return resp, helper.NewResponseError(helper.ErrNotAllowed, http.StatusForbidden)
	}

	photo.Title = data.Title
	photo.URL = data.URL
	photo.Caption.String = data.Caption
	photo.Caption.Valid = true

	photo, err = s.photoRepo.Update(ctx, photo)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrUpdateConflict, http.StatusConflict)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = dto.PhotoUpdateResponse{
		ID:        photo.ID,
		Title:     photo.Title,
		Caption:   photo.Caption.String,
		URL:       photo.URL,
		UserID:    photo.UserID,
		UpdatedAt: photo.UpdatedAt,
	}

	return resp, nil
}

func (s *photoService) Delete(ctx context.Context, id uint64) (err error) {
	userID, ok := ctx.Value(helper.UserIDKey).(float64)
	if !ok {
		s.logger.ErrorContext(ctx, "ctx.Value(helper.UserIDKey).(float64): userID is not float64")
		return helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	err = s.photoRepo.Delete(ctx, model.Photo{
		ID:     id,
		UserID: uint64(userID),
	})
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewResponseError(helper.ErrNotAllowed, http.StatusForbidden)
		}
		return helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	return nil
}

func (s *photoService) GetByID(ctx context.Context, id uint64) (dto.PhotoResponse, error) {
	var resp dto.PhotoResponse

	photo, err := s.photoRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrPhotoNotFound, http.StatusNotFound)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = dto.PhotoResponse{
		ID:        photo.ID,
		Title:     photo.Title,
		URL:       photo.URL,
		UserID:    photo.UserID,
		CreatedAt: photo.CreatedAt,
		UpdatedAt: photo.UpdatedAt,
		User: dto.User{
			ID:       photo.UserID,
			Email:    photo.User.Email,
			Username: photo.User.Username,
		},
	}

	if photo.Caption.Valid {
		resp.Caption = photo.Caption.String
	}

	return resp, nil
}

func (s *photoService) GetByUserID(ctx context.Context, userID uint64) ([]dto.PhotoResponse, error) {
	var resp []dto.PhotoResponse

	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrUserNotFound, http.StatusNotFound)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	photos, err := s.photoRepo.FindByUserID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = make([]dto.PhotoResponse, 0, len(photos))

	for _, photo := range photos {
		item := dto.PhotoResponse{
			ID:        photo.ID,
			Title:     photo.Title,
			URL:       photo.URL,
			UserID:    photo.UserID,
			CreatedAt: photo.CreatedAt,
			UpdatedAt: photo.UpdatedAt,
			User: dto.User{
				ID:       photo.UserID,
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
		}

		if photo.Caption.Valid {
			item.Caption = photo.Caption.String
		}

		resp = append(resp, item)
	}

	return resp, nil
}

func (s *photoService) GetByUsername(ctx context.Context, username string) ([]dto.PhotoResponse, error) {
	var resp []dto.PhotoResponse

	_, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrUserNotFound, http.StatusNotFound)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	photos, err := s.photoRepo.FindByUsername(ctx, username)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = make([]dto.PhotoResponse, 0, len(photos))

	for _, photo := range photos {
		item := dto.PhotoResponse{
			ID:        photo.ID,
			Title:     photo.Title,
			URL:       photo.URL,
			UserID:    photo.UserID,
			CreatedAt: photo.CreatedAt,
			UpdatedAt: photo.UpdatedAt,
			User: dto.User{
				ID:       photo.UserID,
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
		}

		if photo.Caption.Valid {
			item.Caption = photo.Caption.String
		}

		resp = append(resp, item)
	}

	return resp, nil
}
