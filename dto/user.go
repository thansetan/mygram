package dto

import (
	"errors"
	"final-project/helper"
	"time"
)

var isValidEmail = helper.IsValidEmailRegex(helper.StackOverflowEmailPattern)

type UserRequest struct {
	Username string `json:"username" example:"budiganteng"`
	Email    string `json:"email" example:"budi@rocketmail.com"`
	Password string `json:"password" example:"budigantengbanget123"`
	Age      uint64 `json:"age" example:"25"`
}

func (u UserRequest) ValidateCreate() error {
	var errs error

	if u.Email == "" {
		errs = errors.Join(errs, helper.ErrEmptyEmail)
	} else if !isValidEmail(u.Email) {
		errs = errors.Join(errs, helper.ErrInvalidEmail)
	}

	if u.Username == "" {
		errs = errors.Join(errs, helper.ErrEmptyUsername)
	} else if len(u.Username) > 100 {
		errs = errors.Join(errs, helper.ErrUsernameTooLong)
	}

	if u.Password == "" {
		errs = errors.Join(errs, helper.ErrEmptyPassword)
	} else if u.Password != "" && len(u.Password) < 6 {
		errs = errors.Join(errs, helper.ErrPasswordTooShort)
	} else if len(u.Password) > 72 {
		// bcrypt.GenerateFromPassword only accepts at most 72 characters
		errs = errors.Join(errs, helper.ErrPasswordTooLong)
	}

	if u.Age < 8 {
		errs = errors.Join(errs, helper.ErrAgeTooYoung)
	}

	return errs
}

type UserCreateResponse struct {
	ID       uint64 `json:"id"`
	Age      uint64 `json:"age"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (u UserRequest) ValidateLogin() error {
	var errs error

	if u.Email == "" {
		errs = errors.Join(errs, helper.ErrEmptyEmail)
	} else if !isValidEmail(u.Email) {
		errs = errors.Join(errs, helper.ErrInvalidEmail)
	}

	if u.Password == "" {
		errs = errors.Join(errs, helper.ErrEmptyPassword)
	} else if u.Password != "" && len(u.Password) < 6 {
		errs = errors.Join(errs, helper.ErrPasswordTooShort)
	}

	return errs
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

func (u UserRequest) ValidateUpdate() error {
	var errs error

	if u.Username == "" {
		errs = errors.Join(errs, helper.ErrEmptyUsername)
	} else if len(u.Username) > 100 {
		errs = errors.Join(errs, helper.ErrUsernameTooLong)
	}

	if u.Email == "" {
		errs = errors.Join(errs, helper.ErrEmptyEmail)
	} else if !isValidEmail(u.Email) {
		errs = errors.Join(errs, helper.ErrInvalidEmail)
	}

	return errs
}

type UserUpdateResponse struct {
	ID        uint64    `json:"id"`
	Age       uint64    `json:"age"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
