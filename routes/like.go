package routes

import (
	"database/sql"
	"final-project/controller"
	"final-project/middleware"
	likerepository "final-project/repository/like"
	photorepository "final-project/repository/photo"
	likeservice "final-project/service/like"
	"log/slog"
	"net/http"
)

func InitLikeRoutes(r *http.ServeMux, db *sql.DB, logger *slog.Logger) {
	photoRepo := photorepository.New(db)
	likeRepo := likerepository.New(db)
	likeService := likeservice.New(likeRepo, photoRepo, logger)
	controller := controller.NewLikeController(likeService)

	r.Handle("POST /photos/{photoID}/likes", middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.Create))))
	r.Handle("GET /photos/{photoID}/likes", middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.FindByPhotoID))))
	r.Handle("DELETE /photos/{photoID}/likes", middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.Delete))))
	r.Handle("GET /likes/my", middleware.Auth(middleware.RateLimit(http.HandlerFunc(controller.GetMine))))
}
