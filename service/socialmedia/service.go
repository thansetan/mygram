package socialmediaservice

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

type socialMediaService struct {
	socialMediaRepo repository.SocialMediaRepository
	logger          *slog.Logger
}

func New(socialMediaRepo repository.SocialMediaRepository, logger *slog.Logger) *socialMediaService {
	return &socialMediaService{socialMediaRepo, logger}
}

func (s *socialMediaService) Create(ctx context.Context, data dto.SocialMediaRequest) (dto.SocialMediaCreateResponse, error) {
	var (
		resp dto.SocialMediaCreateResponse
		err  error
	)

	userID, ok := ctx.Value(helper.UserIDKey).(float64)
	if !ok {
		s.logger.ErrorContext(ctx, "ctx.Value(helper.UserIDKey).(float64): userID is not float64")
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	socialMedia := model.SocialMedia{
		UserID: uint64(userID),
		Name:   data.Name,
		URL:    data.URL,
	}

	socialMedia, err = s.socialMediaRepo.Save(ctx, socialMedia)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = dto.SocialMediaCreateResponse{
		ID:        socialMedia.ID,
		UserID:    socialMedia.UserID,
		Name:      socialMedia.Name,
		URL:       socialMedia.URL,
		CreatedAt: socialMedia.CreatedAt,
	}

	return resp, nil
}

func (s *socialMediaService) GetAll(ctx context.Context) ([]dto.SocialMediaResponse, error) {
	var resp []dto.SocialMediaResponse

	socialMedias, err := s.socialMediaRepo.FindAll(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = make([]dto.SocialMediaResponse, 0, len(socialMedias))

	for _, socialMedia := range socialMedias {
		resp = append(resp, dto.SocialMediaResponse{
			ID:        socialMedia.ID,
			UserID:    socialMedia.UserID,
			Name:      socialMedia.Name,
			URL:       socialMedia.URL,
			CreatedAt: socialMedia.CreatedAt,
			UpdatedAt: socialMedia.UpdatedAt,
			User: dto.User{
				ID:       socialMedia.UserID,
				Username: socialMedia.User.Username,
				Email:    socialMedia.User.Email,
			},
		})
	}

	return resp, nil
}

func (s *socialMediaService) Update(ctx context.Context, id uint64, data dto.SocialMediaRequest) (resp dto.SocialMediaUpdateResponse, err error) {
	userID, ok := ctx.Value(helper.UserIDKey).(float64)
	if !ok {
		s.logger.ErrorContext(ctx, "ctx.Value(helper.UserIDKey).(float64): userID is not float64")
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	socialMedia, err := s.socialMediaRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrNotAllowed, http.StatusForbidden)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	if socialMedia.UserID != uint64(userID) {
		s.logger.ErrorContext(ctx, "socialMedia.UserID != userID: user is not the owner of the social media")
		return resp, helper.NewResponseError(helper.ErrNotAllowed, http.StatusForbidden)
	}

	socialMedia.Name = data.Name
	socialMedia.URL = data.URL

	socialMedia, err = s.socialMediaRepo.Update(ctx, socialMedia)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrUpdateConflict, http.StatusConflict)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = dto.SocialMediaUpdateResponse{
		ID:        socialMedia.ID,
		UserID:    socialMedia.UserID,
		Name:      socialMedia.Name,
		URL:       socialMedia.URL,
		UpdatedAt: socialMedia.UpdatedAt,
	}

	return resp, nil
}

func (s *socialMediaService) Delete(ctx context.Context, id uint64) (err error) {
	userID, ok := ctx.Value(helper.UserIDKey).(float64)
	if !ok {
		s.logger.ErrorContext(ctx, "ctx.Value(helper.UserIDKey).(float64): userID is not float64")
		return helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	err = s.socialMediaRepo.Delete(ctx, model.SocialMedia{
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

func (s *socialMediaService) GetByID(ctx context.Context, id uint64) (dto.SocialMediaResponse, error) {
	var resp dto.SocialMediaResponse

	socialMedia, err := s.socialMediaRepo.FindByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrSocialMediaNotFound, http.StatusNotFound)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = dto.SocialMediaResponse{
		ID:        socialMedia.ID,
		UserID:    socialMedia.UserID,
		Name:      socialMedia.Name,
		URL:       socialMedia.URL,
		CreatedAt: socialMedia.CreatedAt,
		UpdatedAt: socialMedia.UpdatedAt,
		User: dto.User{
			ID:       socialMedia.UserID,
			Username: socialMedia.User.Username,
			Email:    socialMedia.User.Email,
		},
	}

	return resp, nil
}

func (s *socialMediaService) GetByUserID(ctx context.Context, userID uint64) ([]dto.SocialMediaGetByUserIDResponse, error) {
	var resp []dto.SocialMediaGetByUserIDResponse

	socialMedias, err := s.socialMediaRepo.FindByUserID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = make([]dto.SocialMediaGetByUserIDResponse, 0, len(socialMedias))

	for _, socialMedia := range socialMedias {
		resp = append(resp, dto.SocialMediaGetByUserIDResponse{
			ID:        socialMedia.ID,
			UserID:    socialMedia.UserID,
			Name:      socialMedia.Name,
			URL:       socialMedia.URL,
			CreatedAt: socialMedia.CreatedAt,
			UpdatedAt: socialMedia.UpdatedAt,
		})
	}

	return resp, nil
}
