import { useState } from 'react';
import { formatBertanda, formatWaktu } from '../format';

const FILTER = [
  ['', 'Semua'],
  ['pemasukan', 'Pemasukan'],
  ['pengeluaran', 'Pengeluaran'],
];

export default function RiwayatTab({ transaksi, onHapus }) {
  const [filter, setFilter] = useState('');
  const tampil = filter ? transaksi.filter((t) => t.tipe === filter) : transaksi;

  return (
    <div className="panel">
      <div className="panel-head">
        <h2>Riwayat</h2>
        <div className="seg">
          {FILTER.map(([val, label]) => (
            <button key={val} className={filter === val ? 'active' : ''} onClick={() => setFilter(val)}>{label}</button>
          ))}
        </div>
      </div>

      {tampil.length === 0 ? (
        <p className="empty">Belum ada transaksi.</p>
      ) : (
        <ul className="history">
          {tampil.map((t) => (
            <li key={t.id}>
              <div>
                <div className="ket">
                  {t.keterangan}
                  {t.qty > 1 && <span className="muted"> × {t.qty}</span>}
                </div>
                <div className="muted small-text">{formatWaktu(t.created_at)}</div>
              </div>
              <div className="right">
                <strong className={t.tipe === 'pemasukan' ? 'pos' : 'neg'}>{formatBertanda(t.tipe, t.jumlah)}</strong>
                <button className="btn small danger" onClick={() => onHapus(t)} aria-label="Hapus transaksi">✕</button>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
