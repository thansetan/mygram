package userservice

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

type userService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

func New(userRepo repository.UserRepository, logger *slog.Logger) *userService {
	return &userService{userRepo, logger}
}

func (s *userService) Create(ctx context.Context, data dto.UserRequest) (dto.UserCreateResponse, error) {
	var (
		resp dto.UserCreateResponse
		user model.User
		err  error
	)

	user.Username = data.Username
	user.Email = data.Email
	user.Age = data.Age

	user.Password, err = helper.HashPassword(data.Password)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	user, err = s.userRepo.Save(ctx, user)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		pqErr := new(pq.Error)
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == "unique_violation" {
				return resp, helper.NewResponseError(helper.ErrDuplicateUserEmail, http.StatusConflict)
			}
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp.ID = user.ID
	resp.Username = user.Username
	resp.Email = user.Email
	resp.Age = user.Age

	return resp, nil
}

func (s *userService) Login(ctx context.Context, data dto.UserRequest) (dto.UserLoginResponse, error) {
	var resp dto.UserLoginResponse

	user, err := s.userRepo.FindByEmail(ctx, data.Email)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrInvalidLogin, http.StatusUnauthorized)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	if !helper.IsValidPassword(user.Password, data.Password) {
		s.logger.ErrorContext(ctx, "invalid password")
		return resp, helper.NewResponseError(helper.ErrInvalidLogin, http.StatusUnauthorized)
	}

	resp.Token, err = helper.GenerateJWT(user.ID)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	return resp, nil
}

func (s *userService) Update(ctx context.Context, data dto.UserRequest) (resp dto.UserUpdateResponse, err error) {

	userID, ok := ctx.Value(helper.UserIDKey).(float64)
	if !ok {
		s.logger.ErrorContext(ctx, "ctx.Value(helper.UserIDKey).(float64): userID is not float64")
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	user, err := s.userRepo.FindByID(ctx, uint64(userID))
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrNotAllowed, http.StatusForbidden)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	user.Email = data.Email
	user.Username = data.Username

	user, err = s.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		pqErr := new(pq.Error)
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == "unique_violation" {
				return resp, helper.NewResponseError(helper.ErrDuplicateUserEmail, http.StatusConflict)
			}
		}
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrUpdateConflict, http.StatusConflict)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = dto.UserUpdateResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Age:       user.Age,
		UpdatedAt: user.UpdatedAt,
	}

	return resp, nil
}

func (s *userService) Delete(ctx context.Context) error {
	userID, ok := ctx.Value(helper.UserIDKey).(float64)
	if !ok {
		s.logger.ErrorContext(ctx, "ctx.Value(helper.UserIDKey).(float64): userID is not float64")
		return helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	err := s.userRepo.Delete(ctx, uint64(userID))
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewResponseError(helper.ErrUserNotFound, http.StatusNotFound)
		}
		return helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	return nil
}
