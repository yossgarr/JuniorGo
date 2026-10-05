import { useCallback, useEffect, useRef, useState } from 'react';
import { api } from './api';
import AuthForm from './components/AuthForm';
import BarangTab from './components/BarangTab';
import CatatTab from './components/CatatTab';
import RiwayatTab from './components/RiwayatTab';
import Ringkasan from './components/Ringkasan';
import './index.css';

const TABS = [
  ['catat', 'Catat'],
  ['barang', 'Barang'],
  ['riwayat', 'Riwayat'],
];

export default function App() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [tab, setTab] = useState('catat');

  const [ringkasan, setRingkasan] = useState(null);
  const [barang, setBarang] = useState([]);
  const [transaksi, setTransaksi] = useState([]);

  const [toast, setToast] = useState(null);
  const toastTimer = useRef(null);

  const tampilToast = useCallback((text, type = 'ok') => {
    clearTimeout(toastTimer.current);
    setToast({ text, type });
    toastTimer.current = setTimeout(() => setToast(null), 2500);
  }, []);

  const muatSemua = useCallback(async () => {
    const [r, b, t] = await Promise.all([api.ringkasan(), api.daftarBarang(), api.daftarTransaksi()]);
    setRingkasan(r);
    setBarang(b || []);
    setTransaksi(t || []);
  }, []);

  // Menangani error API: sesi habis -> kembali ke login, selain itu tampilkan pesan
  const tangani = useCallback(
    (err) => {
      if (err.status === 401) setUser(null);
      else tampilToast(err.message, 'error');
    },
    [tampilToast],
  );

  // Cek sesi (cookie) saat halaman dibuka
  useEffect(() => {
    api
      .me()
      .then((data) => {
        setUser(data.user);
        // Gagal memuat data tidak boleh dianggap sesi habis
        muatSemua().catch(tangani);
      })
      .catch(() => setUser(null))
      .finally(() => setLoading(false));
  }, [muatSemua, tangani]);

  const handleLogin = async (username) => {
    setUser(username);
    try {
      await muatSemua();
    } catch (err) {
      tangani(err);
    }
  };

  const handleLogout = async () => {
    try {
      await api.logout();
    } finally {
      setUser(null);
      setRingkasan(null);
      setBarang([]);
      setTransaksi([]);
      setTab('catat');
    }
  };

  // Mengembalikan true jika berhasil, agar form bisa dikosongkan
  const jalankan = async (aksi, pesan) => {
    try {
      await aksi();
      await muatSemua();
      if (pesan) tampilToast(pesan);
      return true;
    } catch (err) {
      tangani(err);
      return false;
    }
  };

  const catat = (payload, pesan) => jalankan(() => api.catatTransaksi(payload), pesan);

  const simpanBarang = (id, data) =>
    jalankan(() => (id ? api.ubahBarang(id, data) : api.tambahBarang(data)), id ? 'Barang diperbarui' : 'Barang ditambahkan');

  const hapusBarang = (b) => {
    if (!window.confirm(`Hapus "${b.nama}"? Riwayat transaksinya tetap tersimpan.`)) return;
    jalankan(() => api.hapusBarang(b.id), 'Barang dihapus');
  };

  const hapusTransaksi = (t) => {
    if (!window.confirm(`Hapus transaksi "${t.keterangan}"? Saldo akan dihitung ulang.`)) return;
    jalankan(() => api.hapusTransaksi(t.id), 'Transaksi dihapus');
  };

  if (loading) return <p className="center-note">Memeriksa sesi login...</p>;

  if (!user) return <AuthForm onLogin={handleLogin} />;

  return (
    <div className="app">
      <header className="topbar">
        <h1>Manage Money By YossTakke</h1>
        <div className="user">
          <span>Halo, <b>{user}</b></span>
          <button className="btn small" onClick={handleLogout}>Keluar</button>
        </div>
      </header>

      <Ringkasan data={ringkasan} />

      <nav className="tabs">
        {TABS.map(([key, label]) => (
          <button key={key} className={tab === key ? 'active' : ''} onClick={() => setTab(key)}>{label}</button>
        ))}
      </nav>

      {tab === 'catat' && <CatatTab barang={barang} onCatat={catat} onKeBarang={() => setTab('barang')} />}
      {tab === 'barang' && <BarangTab barang={barang} onSimpan={simpanBarang} onHapus={hapusBarang} />}
      {tab === 'riwayat' && <RiwayatTab transaksi={transaksi} onHapus={hapusTransaksi} />}

      {toast && <div className={`toast ${toast.type}`} role="status">{toast.text}</div>}
    </div>
  );
}
