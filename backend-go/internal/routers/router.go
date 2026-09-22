package routers

import (
	"net/http"

	"backend-go/internal/controller"
	"backend-go/internal/middlewares"
	"backend-go/internal/repo"
	"backend-go/internal/service"
)

func InitRouter() http.Handler {
	mux := http.NewServeMux()

	// 1. Modul Buku
	bukuRepo := repo.NewBukuRepository()
	bukuSvc := service.NewBukuService(bukuRepo)
	bukuCtrl := controller.NewBukuController(bukuSvc)

	// 2. Modul Auth
	userRepo := repo.NewUserRepository()
	userSvc := service.NewUserService(userRepo)
	authCtrl := controller.NewAuthController(userSvc)

	// Public Routes
	mux.HandleFunc("/register", authCtrl.Register)
	mux.HandleFunc("/login", authCtrl.Login)

	// Protected Routes (Wajib JWT)
	mux.HandleFunc("/buku", middlewares.RequireAuth(bukuCtrl.HandleBuku))

	// Endpoint OAuth Google
	mux.HandleFunc("/auth/google/login", authCtrl.GoogleLogin)
	mux.HandleFunc("/auth/google/callback", authCtrl.GoogleCallback)

	return middlewares.EnableCORS(middlewares.Logger(mux))
}