package helper

import (
	"errors"
)

var (
	ErrInternal              = errors.New("there's something wrong. it's our fault, not yours")
	ErrNotAllowed            = errors.New("you're not allowed to perform this action")
	ErrDuplicateUserEmail    = errors.New("user with given username or email already exists")
	ErrInvalidLogin          = errors.New("invalid email or password")
	ErrUserNotFound          = errors.New("user not found")
	ErrPhotoNotFound         = errors.New("photo with given id not found")
	ErrLikeNotFound          = errors.New("you haven't liked this photo yet")
	ErrMultipleLikes         = errors.New("you've liked this photo")
	ErrCommentNotFound       = errors.New("comment with given id not found")
	ErrSocialMediaNotFound   = errors.New("social media with given id not found")
	ErrNotLoggedIn           = errors.New("you're not logged in")
	ErrInvalidID             = errors.New("id must be a positive integer")
	ErrInvalidContentType    = errors.New("invalid content type")
	ErrEmptyEmail            = errors.New("email can't be empty")
	ErrInvalidEmail          = errors.New("invalid email format")
	ErrEmptyUsername         = errors.New("username can't be empty")
	ErrUsernameTooLong       = errors.New("username can't be more than 100 characters")
	ErrEmptyPassword         = errors.New("password can't be empty")
	ErrPasswordTooShort      = errors.New("password must be at least 6 characters")
	ErrPasswordTooLong       = errors.New("password can't be more than 72 characters")
	ErrAgeTooYoung           = errors.New("age must be at least 8 years old")
	ErrEmptyTitle            = errors.New("title can't be empty")
	ErrTitleTooLong          = errors.New("title can't be more than 100 characters")
	ErrEmptyPhotoURL         = errors.New("photo_url can't be empty")
	ErrInvalidPhotoURL       = errors.New("invalid photo_url format")
	ErrEmptyMessage          = errors.New("message can't be empty")
	ErrEmptyPhotoID          = errors.New("photo_id can't be empty")
	ErrEmptyName             = errors.New("name can't be empty")
	ErrEmptySocialMediaURL   = errors.New("social_media_url can't be empty")
	ErrInvalidSocialMediaURL = errors.New("invalid social_media_url format")
	ErrInvalidJWT            = errors.New("invalid JWT token")
	ErrInvalidBasePath       = errors.New("base_path must start and end with a single '/' and can't contain any special characters except '-' and '/'")
	ErrInvalidDuration       = errors.New("invalid duration format")
	ErrUpdateConflict        = errors.New("the data you're trying to update has been modified by someone else")
)

type ResponseError struct {
	err  error
	code int
}

func NewResponseError(err error, code int) error {
	return &ResponseError{
		err:  err,
		code: code,
	}
}

func (r *ResponseError) Error() string {
	return r.err.Error()
}

func (r *ResponseError) Code() int {
	return r.code
}
