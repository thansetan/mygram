package dto

import (
	"errors"
	"final-project/helper"
	"time"
)

type CommentRequest struct {
	Message string `json:"message"`
	PhotoID uint64 `json:"-"`
}

func (c CommentRequest) ValidateCreate() error {
	var errs error

	if c.Message == "" {
		errs = errors.Join(errs, helper.ErrEmptyMessage)
	}

	return errs
}

type CommentCreateResponse struct {
	ID        uint64    `json:"id"`
	Message   string    `json:"message"`
	PhotoID   uint64    `json:"photo_id"`
	UserID    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentResponse struct {
	ID        uint64    `json:"id"`
	Message   string    `json:"message"`
	PhotoID   uint64    `json:"photo_id"`
	UserID    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
	User      User      `json:"user"`
	Photo     Photo     `json:"photo"`
}

func (c CommentRequest) ValidateUpdate() error {
	var errs error

	if c.Message == "" {
		errs = errors.Join(errs, helper.ErrEmptyMessage)
	}

	return errs
}

type CommentUpdateResponse struct {
	ID        uint64    `json:"id"`
	Message   string    `json:"message"`
	PhotoID   uint64    `json:"photo_id"`
	UserID    uint64    `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentGetByPhotoIDResponse struct {
	ID        uint64    `json:"id"`
	Message   string    `json:"message"`
	PhotoID   uint64    `json:"photo_id"`
	UserID    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
	User      User      `json:"user"`
}

type CommentGetByUserIDResponse struct {
	ID        uint64    `json:"id"`
	Message   string    `json:"message"`
	PhotoID   uint64    `json:"photo_id"`
	UserID    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"updated_at"`
	Photo     Photo     `json:"photo"`
}
