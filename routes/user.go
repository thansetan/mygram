package routes

import (
	"database/sql"
	"final-project/controller"
	"final-project/middleware"
	userrepository "final-project/repository/user"
	userservice "final-project/service/user"
	"log/slog"
	"net/http"
)

func InitUserRoutes(r *http.ServeMux, db *sql.DB, logger *slog.Logger) {
	userRepo := userrepository.New(db)
	userService := userservice.New(userRepo, logger)
	userController := controller.NewUserController(userService)

	r.Handle("POST /users/register", middleware.AllowedContentType(http.HandlerFunc(userController.Register)))
	r.Handle("POST /users/login", middleware.AllowedContentType(http.HandlerFunc(userController.Login)))
	r.Handle("PUT /users", middleware.AllowedContentType(middleware.Auth(middleware.RateLimit(http.HandlerFunc(userController.Update)))))
	r.Handle("DELETE /users", middleware.Auth(middleware.RateLimit(http.HandlerFunc(userController.Delete))))
}
