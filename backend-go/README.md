# Catatan Uang — Backend (Go)

Aplikasi manajemen uang: pemasukan, pengeluaran, dan katalog barang (acuan pengeluaran).

## Endpoint (cookie auth; versi `/api/...` memakai API token)
| Method | Path | Fungsi |
|---|---|---|
| GET/POST/PUT/DELETE | `/barang` (`?id=`) | CRUD barang + harga |
| GET/POST/DELETE | `/transaksi` (`?id=`, `?tipe=`) | Riwayat / catat / hapus transaksi |
| GET | `/ringkasan` | Saldo, total pemasukan, total pengeluaran |

`POST /transaksi`:
- Pengeluaran dari barang: `{"tipe":"pengeluaran","barang_id":1,"qty":2}` — nominal = harga × qty, dihitung server.
- Manual: `{"tipe":"pemasukan","jumlah":100000,"keterangan":"Gaji"}`

Tabel dibuat otomatis saat server start (lihat `internal/initialize/db.go`, salinan SQL di `migrations/`).
Contoh request: `api.http`.
