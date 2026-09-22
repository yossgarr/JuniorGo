package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend-go/internal/models"
	"backend-go/internal/service"
	"backend-go/pkg/response"
)

type BukuController struct {
	svc *service.BukuService
}

func NewBukuController(s *service.BukuService) *BukuController {
	return &BukuController{svc: s}
}

func (c *BukuController) HandleBuku(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.AmbilDaftarBuku(w, r)
	case http.MethodPost:
		c.TambahBuku(w, r)
	case http.MethodPut:
		c.UbahBuku(w, r)
	case http.MethodDelete:
		c.HapusBuku(w, r)
	default:
		response.Gagal(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
	}
}

func (c *BukuController) AmbilDaftarBuku(w http.ResponseWriter, r *http.Request) {
	data, err := c.svc.DapatkanSemua()
	if err != nil {
		response.Gagal(w, http.StatusInternalServerError, "Gagal mengambil data")
		return
	}
	response.Sukses(w, http.StatusOK, "Berhasil memuat daftar buku", data)
}

func (c *BukuController) TambahBuku(w http.ResponseWriter, r *http.Request) {
	var bukuBaru models.Buku
	if err := json.NewDecoder(r.Body).Decode(&bukuBaru); err != nil {
		response.Gagal(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}
	if err := c.svc.BuatBuku(&bukuBaru); err != nil {
		response.Gagal(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Sukses(w, http.StatusCreated, "Buku berhasil ditambahkan", bukuBaru)
}

func (c *BukuController) UbahBuku(w http.ResponseWriter, r *http.Request) {
	// Membaca ID dari query param: /buku?id=...
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Gagal(w, http.StatusBadRequest, "ID tidak valid atau kosong")
		return
	}

	var buku models.Buku
	if err := json.NewDecoder(r.Body).Decode(&buku); err != nil {
		response.Gagal(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}
	buku.ID = id

	if err := c.svc.PerbaruiBuku(&buku); err != nil {
		response.Gagal(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Sukses(w, http.StatusOK, "Buku berhasil diperbarui", buku)
}

func (c *BukuController) HapusBuku(w http.ResponseWriter, r *http.Request) {
	// Membaca ID dari query param: /buku?id=...
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Gagal(w, http.StatusBadRequest, "ID tidak valid atau kosong")
		return
	}

	if err := c.svc.HapusBuku(id); err != nil {
		response.Gagal(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.Sukses(w, http.StatusOK, "Buku berhasil dihapus", nil)
}