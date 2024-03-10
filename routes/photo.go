package routes

import (
	"database/sql"
	"final-project/controller"
	"final-project/middleware"
	photorepository "final-project/repository/photo"
	userrepository "final-project/repository/user"
	photoservice "final-project/service/photo"
	"log/slog"
	"net/http"
)

func InitPhotoRoutes(r *http.ServeMux, db *sql.DB, logger *slog.Logger) {
	userRepo := userrepository.New(db)
	photoRepo := photorepository.New(db)
	service := photoservice.New(userRepo, photoRepo, logger)
	controller := controller.NewPhotoController(service)

	r.Handle("POST /photos", middleware.AllowedContentType(middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.Create)))))
	r.Handle("GET /photos", middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.GetAll))))
	r.Handle("PUT /photos/{photoID}", middleware.AllowedContentType(middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.Update)))))
	r.Handle("DELETE /photos/{photoID}", middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.Delete))))
	r.Handle("GET /photos/{photoID}", middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.GetByID))))
	r.Handle("GET /photos/my", middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.GetMine))))
	r.Handle("GET /users/{username}/photos", middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.GetByUsername))))
}
