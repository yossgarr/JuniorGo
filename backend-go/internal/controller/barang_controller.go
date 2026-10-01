package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"backend-go/internal/middlewares"
	"backend-go/internal/models"
	"backend-go/internal/repo"
	"backend-go/internal/service"
	"backend-go/pkg/response"
)

type BarangController struct {
	svc *service.BarangService
}

func NewBarangController(s *service.BarangService) *BarangController {
	return &BarangController{svc: s}
}

func (c *BarangController) HandleBarang(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.UserID(r)
	if !ok {
		response.Gagal(w, http.StatusUnauthorized, "User tidak dikenali")
		return
	}

	switch r.Method {
	case http.MethodGet:
		list, err := c.svc.DapatkanSemua(userID)
		if err != nil {
			response.Gagal(w, http.StatusInternalServerError, "Gagal mengambil data barang")
			return
		}
		response.Sukses(w, http.StatusOK, "Berhasil memuat daftar barang", list)

	case http.MethodPost:
		var b models.Barang
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			response.Gagal(w, http.StatusBadRequest, "Format JSON tidak valid")
			return
		}
		if err := c.svc.Buat(userID, &b); err != nil {
			response.Gagal(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Sukses(w, http.StatusCreated, "Barang berhasil ditambahkan", b)

	case http.MethodPut:
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			response.Gagal(w, http.StatusBadRequest, "ID tidak valid atau kosong")
			return
		}
		var b models.Barang
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			response.Gagal(w, http.StatusBadRequest, "Format JSON tidak valid")
			return
		}
		b.ID = id
		if err := c.svc.Perbarui(userID, &b); err != nil {
			response.Gagal(w, statusDariError(err), err.Error())
			return
		}
		response.Sukses(w, http.StatusOK, "Barang berhasil diperbarui", b)

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
		response.Sukses(w, http.StatusOK, "Barang berhasil dihapus", nil)

	default:
		response.Gagal(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}

// statusDariError memetakan error "tidak ditemukan" ke 404, selain itu 400.
func statusDariError(err error) int {
	if errors.Is(err, repo.ErrBarangTidakDitemukan) || errors.Is(err, repo.ErrTransaksiTidakDitemukan) {
		return http.StatusNotFound
	}
	return http.StatusBadRequest
}
