package likeservice

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

	"github.com/lib/pq"
)

type likeService struct {
	likeRepository  repository.LikeRepository
	photoRepository repository.PhotoRepository
	logger          *slog.Logger
}

func New(likeRepository repository.LikeRepository, photoRepository repository.PhotoRepository, logger *slog.Logger) *likeService {
	return &likeService{likeRepository, photoRepository, logger}
}

func (s *likeService) Create(ctx context.Context, data dto.LikeRequest) (dto.LikeCreateResponse, error) {
	var (
		resp dto.LikeCreateResponse
	)

	userID, ok := ctx.Value(helper.UserIDKey).(float64)
	if !ok {
		s.logger.ErrorContext(ctx, "userID is not float64", "cause", "ctx.Value(helper.UserIDKey).(float64)")
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	_, err := s.photoRepository.FindByID(ctx, data.PhotoID)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error(), "cause", "s.photoRepository.FindByID")
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrPhotoNotFound, http.StatusNotFound)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	like := model.Like{
		UserID:  uint64(userID),
		PhotoID: data.PhotoID,
	}

	like, err = s.likeRepository.Save(ctx, like)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error(), "cause", "s.repository.Create")
		pgErr := new(pq.Error)
		if errors.As(err, &pgErr) {
			if pgErr.Code.Name() == "unique_violation" {
				return resp, helper.NewResponseError(helper.ErrMultipleLikes, http.StatusConflict)
			}
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = dto.LikeCreateResponse{
		ID:        like.ID,
		UserID:    like.UserID,
		PhotoID:   like.PhotoID,
		CreatedAt: like.CreatedAt,
	}

	return resp, nil
}

func (s *likeService) GetByPhotoID(ctx context.Context, photoID uint64) ([]dto.LikeResponse, error) {
	var (
		resp []dto.LikeResponse
	)

	_, err := s.photoRepository.FindByID(ctx, photoID)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrPhotoNotFound, http.StatusNotFound)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	likes, err := s.likeRepository.FindByPhotoID(ctx, photoID)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = make([]dto.LikeResponse, 0, len(likes))

	for _, like := range likes {
		resp = append(resp, dto.LikeResponse{
			ID:        like.ID,
			UserID:    like.UserID,
			PhotoID:   like.PhotoID,
			CreatedAt: like.CreatedAt,
			User: dto.User{
				ID:       like.User.ID,
				Email:    like.User.Email,
				Username: like.User.Username,
			},
		})
	}

	return resp, nil
}

func (s *likeService) Delete(ctx context.Context, photoID uint64) error {
	userID, ok := ctx.Value(helper.UserIDKey).(float64)
	if !ok {
		s.logger.ErrorContext(ctx, "ctx.Value(helper.UserIDKey).(float64): userID is not float64")
		return helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	err := s.likeRepository.Delete(ctx, model.Like{
		UserID:  uint64(userID),
		PhotoID: photoID,
	})
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewResponseError(helper.ErrLikeNotFound, http.StatusNotFound)
		}
		return helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	return nil
}

func (s *likeService) GetByUserID(ctx context.Context, userID uint64) ([]dto.GetLikeByUserIDResponse, error) {
	var (
		resp []dto.GetLikeByUserIDResponse
	)

	likes, err := s.likeRepository.FindByUserID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = make([]dto.GetLikeByUserIDResponse, 0, len(likes))

	for _, like := range likes {
		resp = append(resp, dto.GetLikeByUserIDResponse{
			ID:        like.ID,
			UserID:    like.UserID,
			PhotoID:   like.PhotoID,
			CreatedAt: like.CreatedAt,
			Photo: dto.Photo{
				ID:      like.Photo.ID,
				Title:   like.Photo.Title,
				Caption: like.Photo.Caption.String,
				URL:     like.Photo.URL,
				UserID:  like.Photo.UserID,
			},
		})
	}

	return resp, nil
}
