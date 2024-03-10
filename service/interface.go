package service

import (
	"context"
	"final-project/dto"
)

type UserService interface {
	Create(context.Context, dto.UserRequest) (dto.UserCreateResponse, error)
	Login(context.Context, dto.UserRequest) (dto.UserLoginResponse, error)
	Update(context.Context, dto.UserRequest) (dto.UserUpdateResponse, error)
	Delete(context.Context) error
}

type PhotoService interface {
	Create(context.Context, dto.PhotoRequest) (dto.PhotoCreateResponse, error)
	GetAll(context.Context) ([]dto.PhotoResponse, error)
	Update(context.Context, uint64, dto.PhotoRequest) (dto.PhotoUpdateResponse, error)
	Delete(context.Context, uint64) error
	GetByID(context.Context, uint64) (dto.PhotoResponse, error)
	GetByUserID(context.Context, uint64) ([]dto.PhotoResponse, error)
	GetByUsername(context.Context, string) ([]dto.PhotoResponse, error)
}

type LikeService interface {
	Create(context.Context, dto.LikeRequest) (dto.LikeCreateResponse, error)
	GetByPhotoID(context.Context, uint64) ([]dto.LikeResponse, error)
	Delete(context.Context, uint64) error
	GetByUserID(context.Context, uint64) ([]dto.GetLikeByUserIDResponse, error)
}

type CommentService interface {
	Create(context.Context, dto.CommentRequest) (dto.CommentCreateResponse, error)
	GetAll(context.Context) ([]dto.CommentResponse, error)
	Update(context.Context, uint64, dto.CommentRequest) (dto.CommentUpdateResponse, error)
	Delete(context.Context, uint64) error
	GetByID(context.Context, uint64) (dto.CommentResponse, error)
	GetByPhotoID(context.Context, uint64) ([]dto.CommentGetByPhotoIDResponse, error)
	GetByUserID(context.Context, uint64) ([]dto.CommentGetByUserIDResponse, error)
}

type SocialMediaService interface {
	Create(context.Context, dto.SocialMediaRequest) (dto.SocialMediaCreateResponse, error)
	GetAll(context.Context) ([]dto.SocialMediaResponse, error)
	Update(context.Context, uint64, dto.SocialMediaRequest) (dto.SocialMediaUpdateResponse, error)
	Delete(context.Context, uint64) error
	GetByID(context.Context, uint64) (dto.SocialMediaResponse, error)
	GetByUserID(context.Context, uint64) ([]dto.SocialMediaGetByUserIDResponse, error)
}
