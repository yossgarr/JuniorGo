package controller

import (
	"crypto/md5"
	"encoding/json"
	"encoding/hex"
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
	// 1. Ambil data (yang sudah dioptimasi Redis sebelumnya)
	bukuList, err := c.svc.DapatkanSemua()
	if err != nil {
		response.Gagal(w, http.StatusInternalServerError, "Gagal mengambil data")
		return
	}

	// 2. Buat "sidik jari" (ETag) dari isi data JSON
	dataBytes, _ := json.Marshal(bukuList)
	hash := md5.Sum(dataBytes)
	etag := hex.EncodeToString(hash[:])

	// 3. Pasang Header HTTP Caching
	// - max-age=30: simpan di browser selama 30 detik
	// - must-revalidate: setelah 30 detik, wajib konfirmasi ke server apakah data berubah
	w.Header().Set("Cache-Control", "public, max-age=30, must-revalidate")
	w.Header().Set("ETag", etag)

	// 4. Validasi ETag dari browser
	// Jika browser mengirimkan ETag yang sama, berarti data di browser masih valid
	if r.Header.Get("If-None-Match") == etag {
		// Balas 304 Not Modified (tanpa body data) -> sangat hemat kuota!
		w.WriteHeader(http.StatusNotModified)
		return
	}

	// 5. Kirim data normal (Status 200) jika data memang baru atau pertama kali dibuka
	response.Sukses(w, http.StatusOK, "Berhasil memuat daftar buku", bukuList)
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