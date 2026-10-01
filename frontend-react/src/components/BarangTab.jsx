import { useState } from 'react';
import { formatRupiah } from '../format';

export default function BarangTab({ barang, onSimpan, onHapus }) {
  const [editId, setEditId] = useState(null);
  const [nama, setNama] = useState('');
  const [harga, setHarga] = useState('');

  const reset = () => {
    setEditId(null);
    setNama('');
    setHarga('');
  };

  const submit = async (e) => {
    e.preventDefault();
    const ok = await onSimpan(editId, { nama, harga: parseInt(harga, 10) });
    if (ok) reset();
  };

  const mulaiEdit = (b) => {
    setEditId(b.id);
    setNama(b.nama);
    setHarga(String(b.harga));
  };

  return (
    <div>
      <form className="panel inline-form" onSubmit={submit}>
        <h2>{editId ? 'Edit Barang' : 'Tambah Barang'}</h2>
        <div className="row">
          <input required placeholder="Nama (mis. Indomie)" value={nama} onChange={(e) => setNama(e.target.value)} />
          <input required type="number" min="1" placeholder="Harga (mis. 2000)" value={harga} onChange={(e) => setHarga(e.target.value)} />
          <button className="btn primary">{editId ? 'Perbarui' : 'Simpan'}</button>
          {editId && <button type="button" className="btn" onClick={reset}>Batal</button>}
        </div>
      </form>

      <div className="panel">
        {barang.length === 0 ? (
          <p className="empty">Belum ada barang.</p>
        ) : (
          <table>
            <thead>
              <tr><th>Nama</th><th className="num">Harga</th><th /></tr>
            </thead>
            <tbody>
              {barang.map((b) => (
                <tr key={b.id}>
                  <td>{b.nama}</td>
                  <td className="num">{formatRupiah(b.harga)}</td>
                  <td className="actions">
                    <button className="btn small" onClick={() => mulaiEdit(b)}>Edit</button>
                    <button className="btn small danger" onClick={() => onHapus(b)}>Hapus</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
