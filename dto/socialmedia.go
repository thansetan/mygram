package dto

import (
	"errors"
	"final-project/helper"
	"time"
)

type SocialMediaRequest struct {
	Name string `json:"name"`
	URL  string `json:"social_media_url"`
}

func (s SocialMediaRequest) ValidateCreate() error {
	var errs error

	if s.Name == "" {
		errs = errors.Join(errs, helper.ErrEmptyName)
	}

	if s.URL == "" {
		errs = errors.Join(errs, helper.ErrEmptySocialMediaURL)
	} else if !helper.IsValidURL(s.URL) {
		errs = errors.Join(errs, helper.ErrInvalidSocialMediaURL)
	}

	return errs
}

type SocialMediaCreateResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"social_media_url"`
	UserID    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type SocialMediaResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"social_media_url"`
	UserID    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user"`
}

func (s SocialMediaRequest) ValidateUpdate() error {
	var errs error

	if s.Name == "" {
		errs = errors.Join(errs, helper.ErrEmptyName)
	}

	if s.URL == "" {
		errs = errors.Join(errs, helper.ErrEmptySocialMediaURL)
	} else if !helper.IsValidURL(s.URL) {
		errs = errors.Join(errs, helper.ErrInvalidSocialMediaURL)
	}

	return errs
}

type SocialMediaUpdateResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"social_media_url"`
	UserID    uint64    `json:"user_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SocialMediaGetByUserIDResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"social_media_url"`
	UserID    uint64    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
