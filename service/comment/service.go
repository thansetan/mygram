package commentservice

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

type commentService struct {
	commentRepo repository.CommentRepository
	photoRepo   repository.PhotoRepository
	logger      *slog.Logger
}

func New(commentRepo repository.CommentRepository, photoRepo repository.PhotoRepository, logger *slog.Logger) *commentService {
	return &commentService{commentRepo, photoRepo, logger}
}

func (s *commentService) Create(ctx context.Context, data dto.CommentRequest) (dto.CommentCreateResponse, error) {
	var (
		resp dto.CommentCreateResponse
		err  error
	)

	userID, ok := ctx.Value(helper.UserIDKey).(float64)
	if !ok {
		s.logger.ErrorContext(ctx, "userID is not float64", "cause", "ctx.Value(helper.UserIDKey).(float64)")
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	_, err = s.photoRepo.FindByID(ctx, data.PhotoID)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error(), "cause", "s.photoRepo.FindByID")
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrPhotoNotFound, http.StatusNotFound)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	comment := model.Comment{
		PhotoID: data.PhotoID,
		UserID:  uint64(userID),
		Message: data.Message,
	}

	comment, err = s.commentRepo.Save(ctx, comment)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error(), "cause", "s.commentRepo.Create")
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = dto.CommentCreateResponse{
		ID:        comment.ID,
		PhotoID:   comment.PhotoID,
		UserID:    comment.UserID,
		Message:   comment.Message,
		CreatedAt: comment.CreatedAt,
	}

	return resp, nil
}

func (s *commentService) GetAll(ctx context.Context) ([]dto.CommentResponse, error) {
	var resp []dto.CommentResponse

	comments, err := s.commentRepo.FindAll(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error(), "cause", "s.commentRepo.FindAll")
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = make([]dto.CommentResponse, 0, len(comments))

	for _, comment := range comments {
		resp = append(resp, dto.CommentResponse{
			ID:        comment.ID,
			PhotoID:   comment.PhotoID,
			UserID:    comment.UserID,
			Message:   comment.Message,
			CreatedAt: comment.CreatedAt,
			UpdateAt:  comment.UpdatedAt,
			User: dto.User{
				ID:       comment.UserID,
				Username: comment.User.Username,
				Email:    comment.User.Email,
			},
			Photo: dto.Photo{
				ID:      comment.PhotoID,
				Title:   comment.Photo.Title,
				Caption: comment.Photo.Caption.String,
				URL:     comment.Photo.URL,
				UserID:  comment.Photo.UserID,
			},
		})
	}

	return resp, nil
}

func (s *commentService) Update(ctx context.Context, commentID uint64, data dto.CommentRequest) (resp dto.CommentUpdateResponse, err error) {

	userID, ok := ctx.Value(helper.UserIDKey).(float64)
	if !ok {
		s.logger.ErrorContext(ctx, "userID is not float64", "cause", "ctx.Value(helper.UserIDKey).(float64)")
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	comment, err := s.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error(), "cause", "s.commentRepo.FindByID")
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrNotAllowed, http.StatusForbidden)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	if comment.UserID != uint64(userID) {
		s.logger.ErrorContext(ctx, "user is not the owner of the comment", "cause", "comment.UserID != uint64(userID)")
		return resp, helper.NewResponseError(helper.ErrNotAllowed, http.StatusForbidden)
	}

	comment.Message = data.Message

	comment, err = s.commentRepo.Update(ctx, comment)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error(), "cause", "s.commentRepo.Update")
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrUpdateConflict, http.StatusConflict)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = dto.CommentUpdateResponse{
		ID:        comment.ID,
		PhotoID:   comment.PhotoID,
		UserID:    comment.UserID,
		Message:   comment.Message,
		UpdatedAt: comment.UpdatedAt,
	}

	return resp, nil
}

func (s *commentService) Delete(ctx context.Context, commentID uint64) (err error) {
	userID, ok := ctx.Value(helper.UserIDKey).(float64)
	if !ok {
		s.logger.ErrorContext(ctx, "userID is not float64", "cause", "ctx.Value(helper.UserIDKey).(float64)")
		return helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	err = s.commentRepo.Delete(ctx, model.Comment{
		ID:     commentID,
		UserID: uint64(userID),
	})
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error(), "cause", "s.commentRepo.Delete")
		if errors.Is(err, sql.ErrNoRows) {
			return helper.NewResponseError(helper.ErrNotAllowed, http.StatusForbidden)
		}
		return helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	return nil
}

func (s *commentService) GetByID(ctx context.Context, commentID uint64) (dto.CommentResponse, error) {
	var resp dto.CommentResponse

	comment, err := s.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error(), "cause", "s.commentRepo.FindByID")
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrCommentNotFound, http.StatusNotFound)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = dto.CommentResponse{
		ID:        comment.ID,
		PhotoID:   comment.PhotoID,
		UserID:    comment.UserID,
		Message:   comment.Message,
		CreatedAt: comment.CreatedAt,
		UpdateAt:  comment.UpdatedAt,
		User: dto.User{
			ID:       comment.UserID,
			Username: comment.User.Username,
			Email:    comment.User.Email,
		},
		Photo: dto.Photo{
			ID:      comment.PhotoID,
			Title:   comment.Photo.Title,
			Caption: comment.Photo.Caption.String,
			URL:     comment.Photo.URL,
			UserID:  comment.Photo.UserID,
		},
	}

	return resp, nil
}

func (s *commentService) GetByPhotoID(ctx context.Context, photoID uint64) ([]dto.CommentGetByPhotoIDResponse, error) {
	var resp []dto.CommentGetByPhotoIDResponse

	_, err := s.photoRepo.FindByID(ctx, photoID)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error(), "cause", "s.photoRepo.FindByID")
		if errors.Is(err, sql.ErrNoRows) {
			return resp, helper.NewResponseError(helper.ErrPhotoNotFound, http.StatusNotFound)
		}
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	comments, err := s.commentRepo.FindByPhotoID(ctx, model.Photo{ID: photoID})
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error(), "cause", "s.commentRepo.FindByPhotoID")
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = make([]dto.CommentGetByPhotoIDResponse, 0, len(comments))

	for _, comment := range comments {
		resp = append(resp, dto.CommentGetByPhotoIDResponse{
			ID:        comment.ID,
			PhotoID:   comment.PhotoID,
			UserID:    comment.UserID,
			Message:   comment.Message,
			CreatedAt: comment.CreatedAt,
			UpdateAt:  comment.UpdatedAt,
			User: dto.User{
				ID:       comment.UserID,
				Username: comment.User.Username,
				Email:    comment.User.Email,
			},
		})
	}

	return resp, nil
}

func (s *commentService) GetByUserID(ctx context.Context, userID uint64) ([]dto.CommentGetByUserIDResponse, error) {
	var resp []dto.CommentGetByUserIDResponse

	comments, err := s.commentRepo.FindByUserID(ctx, userID)
	if err != nil {
		s.logger.ErrorContext(ctx, err.Error(), "cause", "s.commentRepo.FindByUserID")
		return resp, helper.NewResponseError(helper.ErrInternal, http.StatusInternalServerError)
	}

	resp = make([]dto.CommentGetByUserIDResponse, 0, len(comments))

	for _, comment := range comments {
		resp = append(resp, dto.CommentGetByUserIDResponse{
			ID:        comment.ID,
			PhotoID:   comment.PhotoID,
			UserID:    comment.UserID,
			Message:   comment.Message,
			CreatedAt: comment.CreatedAt,
			UpdateAt:  comment.UpdatedAt,
			Photo: dto.Photo{
				ID:      comment.PhotoID,
				Title:   comment.Photo.Title,
				Caption: comment.Photo.Caption.String,
				URL:     comment.Photo.URL,
				UserID:  comment.Photo.UserID,
			},
		})
	}

	return resp, nil
}
