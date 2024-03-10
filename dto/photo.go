package dto

import (
	"errors"
	"final-project/helper"
	"time"
)

type PhotoRequest struct {
	Title   string `json:"title"`
	Caption string `json:"caption"`
	URL     string `json:"photo_url"`
}

func (p PhotoRequest) ValidateCreate() error {
	var errs error

	if p.Title == "" {
		errs = errors.Join(errs, helper.ErrEmptyTitle)
	} else if len(p.Title) > 100 {
		errs = errors.Join(errs, helper.ErrTitleTooLong)
	}

	if p.URL == "" {
		errs = errors.Join(errs, helper.ErrEmptyPhotoURL)
	} else if !helper.IsValidURL(p.URL) {
		errs = errors.Join(errs, helper.ErrInvalidPhotoURL)
	}

	return errs
}

type PhotoCreateResponse struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	URL       string    `json:"photo_url"`
	UserID    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type PhotoResponse struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	URL       string    `json:"photo_url"`
	UserID    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `json:"user"`
}

func (p PhotoRequest) ValidateUpdate() error {
	var errs error

	if p.Title == "" {
		errs = errors.Join(errs, helper.ErrEmptyTitle)
	} else if len(p.Title) > 100 {
		errs = errors.Join(errs, helper.ErrTitleTooLong)
	}

	if p.URL == "" {
		errs = errors.Join(errs, helper.ErrEmptyPhotoURL)
	} else if !helper.IsValidURL(p.URL) {
		errs = errors.Join(errs, helper.ErrInvalidPhotoURL)
	}

	return errs
}

type PhotoUpdateResponse struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	URL       string    `json:"photo_url"`
	UserID    uint64    `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Photo struct {
	ID      uint64 `json:"id"`
	Title   string `json:"title"`
	Caption string `json:"caption"`
	URL     string `json:"photo_url"`
	UserID  uint64 `json:"user_id"`
}
