package models

import "time"

const (
	TipePemasukan   = "pemasukan"
	TipePengeluaran = "pengeluaran"
)

// Barang adalah katalog item beserta harganya (acuan pengeluaran).
// Harga dalam rupiah penuh (tanpa desimal).
type Barang struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int64  `json:"harga"`
}

// Transaksi adalah satu catatan pemasukan atau pengeluaran.
// Nama barang & jumlah disimpan sebagai snapshot supaya riwayat tidak
// berubah ketika harga barang diedit atau barang dihapus.
type Transaksi struct {
	ID         int       `json:"id"`
	Tipe       string    `json:"tipe"`
	Jumlah     int64     `json:"jumlah"`
	Keterangan string    `json:"keterangan"`
	BarangID   *int      `json:"barang_id"`
	NamaBarang string    `json:"nama_barang"`
	Qty        int       `json:"qty"`
	CreatedAt  time.Time `json:"created_at"`
}

// TransaksiRequest: jika BarangID diisi (pengeluaran), jumlah dihitung server
// dari harga barang x qty. Jika tidak, Jumlah wajib diisi manual.
type TransaksiRequest struct {
	Tipe       string `json:"tipe"`
	Jumlah     int64  `json:"jumlah"`
	Keterangan string `json:"keterangan"`
	BarangID   int    `json:"barang_id"`
	Qty        int    `json:"qty"`
}

type Ringkasan struct {
	Saldo            int64 `json:"saldo"`
	TotalPemasukan   int64 `json:"total_pemasukan"`
	TotalPengeluaran int64 `json:"total_pengeluaran"`
}
