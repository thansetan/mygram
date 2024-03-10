package routes

import (
	"database/sql"
	"final-project/controller"
	"final-project/middleware"
	commentrepository "final-project/repository/comment"
	photorepository "final-project/repository/photo"
	commentservice "final-project/service/comment"
	"log/slog"
	"net/http"
)

func InitCommentRoutes(r *http.ServeMux, db *sql.DB, logger *slog.Logger) {
	commentRepo := commentrepository.New(db)
	photoRepo := photorepository.New(db)
	service := commentservice.New(commentRepo, photoRepo, logger)
	controller := controller.NewCommentController(service)

	r.Handle("POST /photos/{photoID}/comments", middleware.AllowedContentType(middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.Create)))))
	r.Handle("GET /comments", middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.GetAll))))
	r.Handle("PUT /comments/{commentID}", middleware.AllowedContentType(middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.Update)))))
	r.Handle("DELETE /comments/{commentID}", middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.Delete))))
	r.Handle("GET /comments/{commentID}", middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.GetByID))))
	r.Handle("GET /photos/{photoID}/comments", middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.GetByPhotoID))))
	r.Handle("GET /comments/my", middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.GetMine))))
}
