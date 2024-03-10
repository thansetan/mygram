package routes

import (
	"database/sql"
	"final-project/controller"
	"final-project/middleware"
	socialmediarepository "final-project/repository/socialmedia"
	socialmediaservice "final-project/service/socialmedia"
	"log/slog"
	"net/http"
)

func InitSocialMediaRoutes(r *http.ServeMux, db *sql.DB, logger *slog.Logger) {
	socialMediaRepo := socialmediarepository.New(db)
	socialMediaService := socialmediaservice.New(socialMediaRepo, logger)
	socialMediaController := controller.NewSocialMediaController(socialMediaService)

	r.Handle("POST /socialmedias", middleware.AllowedContentType(middleware.Auth(middleware.RateLimit(http.HandlerFunc(socialMediaController.Create)))))
	r.Handle("GET /socialmedias", middleware.Auth(middleware.RateLimit(http.HandlerFunc(socialMediaController.GetAll))))
	r.Handle("PUT /socialmedias/{socialMediaID}", middleware.AllowedContentType(middleware.Auth(middleware.RateLimit(http.HandlerFunc(socialMediaController.Update)))))
	r.Handle("DELETE /socialmedias/{socialMediaID}", middleware.Auth(middleware.RateLimit(http.HandlerFunc(socialMediaController.Delete))))
	r.Handle("GET /socialmedias/{socialMediaID}", middleware.Auth(middleware.RateLimit(http.HandlerFunc(socialMediaController.GetByID))))
	r.Handle("GET /socialmedias/my", middleware.Auth(middleware.RateLimit(http.HandlerFunc(socialMediaController.GetMine))))
}
