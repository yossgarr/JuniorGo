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

	// 1. Modul Barang & Transaksi (manajemen uang)
	barangRepo := repo.NewBarangRepository()
	transaksiRepo := repo.NewTransaksiRepository()
	barangCtrl := controller.NewBarangController(service.NewBarangService(barangRepo))
	transaksiCtrl := controller.NewTransaksiController(service.NewTransaksiService(transaksiRepo, barangRepo))

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

	// Endpoint diproteksi cookie/session
	mux.HandleFunc("/barang", middlewares.RequireAuth(barangCtrl.HandleBarang))
	mux.HandleFunc("/transaksi", middlewares.RequireAuth(transaksiCtrl.HandleTransaksi))
	mux.HandleFunc("/ringkasan", middlewares.RequireAuth(transaksiCtrl.Ringkasan))

	// Endpoint OAuth Google
	mux.HandleFunc("/auth/google/login", authCtrl.GoogleLogin)
	mux.HandleFunc("/auth/google/callback", authCtrl.GoogleCallback)

	// Inisialisasi token repository & middleware
	tokenRepo := repo.NewTokenRepository()
	tokenCtrl := controller.NewTokenController(tokenRepo)
	tokenAuthMiddleware := middlewares.NewTokenAuthMiddleware(tokenRepo)

	// Mengelola API token butuh login (token terikat ke user yang membuatnya)
	mux.HandleFunc("/api/token/create", middlewares.RequireAuth(tokenCtrl.BuatTokenBaru))
	mux.HandleFunc("/api/token/revoke", middlewares.RequireAuth(tokenCtrl.CabutToken))

	// Endpoint yang sama untuk mesin/otomasi, diproteksi API Token
	mux.HandleFunc("/api/barang", tokenAuthMiddleware.RequireAPIToken(barangCtrl.HandleBarang))
	mux.HandleFunc("/api/transaksi", tokenAuthMiddleware.RequireAPIToken(transaksiCtrl.HandleTransaksi))
	mux.HandleFunc("/api/ringkasan", tokenAuthMiddleware.RequireAPIToken(transaksiCtrl.Ringkasan))

	return middlewares.EnableCORS(middlewares.Logger(mux))
}
