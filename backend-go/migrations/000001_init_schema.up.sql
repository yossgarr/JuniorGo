CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS barang (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nama VARCHAR(255) NOT NULL,
    harga BIGINT NOT NULL CHECK (harga > 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transaksi (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tipe VARCHAR(20) NOT NULL CHECK (tipe IN ('pemasukan', 'pengeluaran')),
    jumlah BIGINT NOT NULL CHECK (jumlah > 0),
    keterangan VARCHAR(255) NOT NULL DEFAULT '',
    barang_id INT REFERENCES barang(id) ON DELETE SET NULL,
    nama_barang VARCHAR(255) NOT NULL DEFAULT '',
    qty INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_barang_user ON barang(user_id);
CREATE INDEX IF NOT EXISTS idx_transaksi_user_waktu ON transaksi(user_id, created_at DESC);
