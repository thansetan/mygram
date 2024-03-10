package repository

import (
	"context"
	"final-project/model"
)

type UserRepository interface {
	Save(context.Context, model.User) (model.User, error)
	FindByEmail(context.Context, string) (model.User, error)
	Update(context.Context, model.User) (model.User, error)
	Delete(context.Context, uint64) error
	FindByID(context.Context, uint64) (model.User, error)
	FindByUsername(context.Context, string) (model.User, error)
}

type PhotoRepository interface {
	Save(context.Context, model.Photo) (model.Photo, error)
	FindAll(context.Context) ([]model.Photo, error)
	Update(context.Context, model.Photo) (model.Photo, error)
	Delete(context.Context, model.Photo) error
	FindByID(context.Context, uint64) (model.Photo, error)
	FindByUserID(context.Context, uint64) ([]model.Photo, error)
	FindByUsername(context.Context, string) ([]model.Photo, error)
}

type CommentRepository interface {
	Save(context.Context, model.Comment) (model.Comment, error)
	FindAll(context.Context) ([]model.Comment, error)
	FindByPhotoID(context.Context, model.Photo) ([]model.Comment, error)
	Update(context.Context, model.Comment) (model.Comment, error)
	Delete(context.Context, model.Comment) error
	FindByID(context.Context, uint64) (model.Comment, error)
	FindByUserID(context.Context, uint64) ([]model.Comment, error)
}

type LikeRepository interface {
	Save(context.Context, model.Like) (model.Like, error)
	FindByPhotoID(context.Context, uint64) ([]model.Like, error)
	Delete(context.Context, model.Like) error
	FindByUserID(context.Context, uint64) ([]model.Like, error)
}

type SocialMediaRepository interface {
	Save(context.Context, model.SocialMedia) (model.SocialMedia, error)
	FindAll(context.Context) ([]model.SocialMedia, error)
	Update(context.Context, model.SocialMedia) (model.SocialMedia, error)
	Delete(context.Context, model.SocialMedia) error
	FindByID(context.Context, uint64) (model.SocialMedia, error)
	FindByUserID(context.Context, uint64) ([]model.SocialMedia, error)
}
