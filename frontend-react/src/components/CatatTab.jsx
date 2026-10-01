import { useState } from 'react';
import { formatRupiah } from '../format';

// Klik barang = langsung mencatat pengeluaran sebesar harga x qty.
export default function CatatTab({ barang, onCatat, onKeBarang }) {
  const [qty, setQty] = useState(1);
  const [jumlah, setJumlah] = useState('');
  const [keterangan, setKeterangan] = useState('');
  const [busyId, setBusyId] = useState(null);

  const klikBarang = async (b) => {
    setBusyId(b.id);
    const ok = await onCatat({ tipe: 'pengeluaran', barang_id: b.id, qty }, `-${formatRupiah(b.harga * qty)} · ${b.nama}`);
    if (ok) setQty(1);
    setBusyId(null);
  };

  const simpanPemasukan = async (e) => {
    e.preventDefault();
    const ok = await onCatat(
      { tipe: 'pemasukan', jumlah: parseInt(jumlah, 10), keterangan },
      `+${formatRupiah(parseInt(jumlah, 10))} · ${keterangan || 'Pemasukan'}`,
    );
    if (ok) {
      setJumlah('');
      setKeterangan('');
    }
  };

  return (
    <div className="two-col">
      <section className="panel">
        <div className="panel-head">
          <h2>Pengeluaran</h2>
          <label className="qty">
            Qty
            <input type="number" min="1" value={qty} onChange={(e) => setQty(Math.max(1, parseInt(e.target.value, 10) || 1))} />
          </label>
        </div>
        {barang.length === 0 ? (
          <p className="empty">
            Belum ada barang. <button className="btn link" onClick={onKeBarang}>Tambah barang dulu</button>
          </p>
        ) : (
          <div className="item-grid">
            {barang.map((b) => (
              <button key={b.id} className="item-btn" disabled={busyId === b.id} onClick={() => klikBarang(b)}>
                <span className="nama">{b.nama}</span>
                <span className="harga">-{formatRupiah(b.harga)}</span>
              </button>
            ))}
          </div>
        )}
      </section>

      <section className="panel">
        <h2>Pemasukan</h2>
        <form onSubmit={simpanPemasukan}>
          <label>
            Jumlah (Rp)
            <input required type="number" min="1" value={jumlah} onChange={(e) => setJumlah(e.target.value)} placeholder="100000" />
          </label>
          <label>
            Keterangan
            <input value={keterangan} onChange={(e) => setKeterangan(e.target.value)} placeholder="Gaji, uang saku, ..." />
          </label>
          <button className="btn success block">Tambah Pemasukan</button>
        </form>
      </section>
    </div>
  );
}
