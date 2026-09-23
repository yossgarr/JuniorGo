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
	mux.HandleFunc("/logout", authCtrl.Logout)

	// Endpoint untuk cek sesi aktif saat halaman di-refresh
	mux.HandleFunc("/me", middlewares.RequireAuth(authCtrl.CekMe))

	// Endpoint buku tetap diproteksi
	mux.HandleFunc("/buku", middlewares.RequireAuth(bukuCtrl.HandleBuku))

	// Endpoint OAuth Google
	mux.HandleFunc("/auth/google/login", authCtrl.GoogleLogin)
	mux.HandleFunc("/auth/google/callback", authCtrl.GoogleCallback)

	// Inisialisasi token repository & middleware
	tokenRepo := repo.NewTokenRepository()
	tokenCtrl := controller.NewTokenController(tokenRepo)
	tokenAuthMiddleware := middlewares.NewTokenAuthMiddleware(tokenRepo)

	// Endpoint membuat token baru
	mux.HandleFunc("/api/token/create", tokenCtrl.BuatTokenBaru)
	mux.HandleFunc("/api/token/revoke", tokenCtrl.CabutToken)

	// Endpoint buku diproteksi menggunakan Token Authentication statis
	mux.HandleFunc("/api/buku", tokenAuthMiddleware.RequireAPIToken(bukuCtrl.HandleBuku))


	return middlewares.EnableCORS(middlewares.Logger(mux))
}