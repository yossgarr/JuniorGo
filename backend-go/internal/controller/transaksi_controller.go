package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend-go/internal/middlewares"
	"backend-go/internal/models"
	"backend-go/internal/service"
	"backend-go/pkg/response"
)

type TransaksiController struct {
	svc *service.TransaksiService
}

func NewTransaksiController(s *service.TransaksiService) *TransaksiController {
	return &TransaksiController{svc: s}
}

// HandleTransaksi: GET (riwayat, filter ?tipe=), POST (catat), DELETE (?id=).
func (c *TransaksiController) HandleTransaksi(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.UserID(r)
	if !ok {
		response.Gagal(w, http.StatusUnauthorized, "User tidak dikenali")
		return
	}

	switch r.Method {
	case http.MethodGet:
		list, err := c.svc.Daftar(userID, r.URL.Query().Get("tipe"))
		if err != nil {
			response.Gagal(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Sukses(w, http.StatusOK, "Berhasil memuat riwayat transaksi", list)

	case http.MethodPost:
		var req models.TransaksiRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Gagal(w, http.StatusBadRequest, "Format JSON tidak valid")
			return
		}
		t, err := c.svc.Catat(userID, req)
		if err != nil {
			response.Gagal(w, statusDariError(err), err.Error())
			return
		}
		response.Sukses(w, http.StatusCreated, "Transaksi berhasil dicatat", t)

	case http.MethodDelete:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			response.Gagal(w, http.StatusBadRequest, "ID tidak valid atau kosong")
			return
		}
		if err := c.svc.Hapus(userID, id); err != nil {
			response.Gagal(w, statusDariError(err), err.Error())
			return
		}
		response.Sukses(w, http.StatusOK, "Transaksi berhasil dihapus", nil)

	default:
		response.Gagal(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}

func (c *TransaksiController) Ringkasan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Gagal(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
		return
	}
	userID, ok := middlewares.UserID(r)
	if !ok {
		response.Gagal(w, http.StatusUnauthorized, "User tidak dikenali")
		return
	}
	s, err := c.svc.Ringkasan(userID)
	if err != nil {
		response.Gagal(w, http.StatusInternalServerError, "Gagal menghitung ringkasan")
		return
	}
	response.Sukses(w, http.StatusOK, "Ringkasan keuangan", s)
}
