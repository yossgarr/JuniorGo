import { useState, useEffect } from 'react';

// Sesuaikan dengan setup Anda, jika pakai Vite proxy biarkan '/api'
const API_URL = '/api'; 

export default function App() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  const [isRegister, setIsRegister] = useState(false);
  const [usernameInput, setUsernameInput] = useState('');
  const [passwordInput, setPasswordInput] = useState('');

  const [bukuList, setBukuList] = useState([]);
  const [bukuId, setBukuId] = useState(null);
  const [judul, setJudul] = useState('');
  const [penulis, setPenulis] = useState('');
  const [stok, setStok] = useState(5); // Default stok untuk buku baru

  // 1. Cek sesi login saat halaman pertama kali dibuka lewat Cookie
  useEffect(() => {
    cekSesi();
  }, []);

  const cekSesi = async () => {
    try {
      const res = await fetch(`${API_URL}/me`, {
        credentials: 'include', 
      });
      if (res.ok) {
        const data = await res.json();
        setUser(data.data.user);
        ambilBuku();
      } else {
        setUser(null);
      }
    } catch {
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  // 2. Ambil data buku
  const ambilBuku = async () => {
    try {
      const res = await fetch(`${API_URL}/buku`, {
        credentials: 'include',
      });
      if (res.ok) {
        const data = await res.json();
        setBukuList(data.data || []);
      }
    } catch {
      alert('Gagal mengambil data buku');
    }
  };

  // 3. Login / Register
  const handleAuth = async (e) => {
    e.preventDefault();
    const endpoint = isRegister ? `${API_URL}/register` : `${API_URL}/login`;

    try {
      const res = await fetch(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ username: usernameInput, password: passwordInput }),
      });
      const data = await res.json();

      if (res.ok) {
        if (isRegister) {
          alert('Registrasi berhasil! Silakan login.');
          setIsRegister(false);
        } else {
          setUser(data.data.username);
          ambilBuku();
        }
      } else {
        alert(data.message || 'Gagal login/register');
      }
    } catch {
      alert('Terjadi kesalahan jaringan');
    }
  };

  // 4. Logout
  const handleLogout = async () => {
    await fetch(`${API_URL}/logout`, {
      method: 'POST',
      credentials: 'include',
    });
    setUser(null);
    setBukuList([]);
  };

  // 5. Simpan Buku (POST / PUT) - Sekarang Mengirim Stok
  const handleSimpanBuku = async (e) => {
    e.preventDefault();
    const method = bukuId ? 'PUT' : 'POST';
    const endpoint = bukuId ? `${API_URL}/buku?id=${bukuId}` : `${API_URL}/buku`;

    try {
      const res = await fetch(endpoint, {
        method,
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ judul, penulis, stok: parseInt(stok) || 0 }),
      });

      if (res.ok) {
        resetForm();
        ambilBuku();
      } else {
        const data = await res.json();
        alert(data.message || 'Gagal menyimpan buku');
      }
    } catch {
      alert('Gagal menghubungi server');
    }
  };

  // 6. Hapus Buku (DELETE)
  const handleHapus = async (id) => {
    if (!window.confirm('Hapus buku ini?')) return;
    try {
      const res = await fetch(`${API_URL}/buku?id=${id}`, {
        method: 'DELETE',
        credentials: 'include',
      });
      if (res.ok) ambilBuku();
    } catch {
      alert('Gagal menghapus');
    }
  };

  // 7. Pinjam Buku (Transaksi ACID)
  const handlePinjam = async (id) => {
    try {
      const res = await fetch(`${API_URL}/buku/pinjam`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ buku_id: id }),
      });
      
      const data = await res.json();
      
      if (res.ok) {
        alert('Berhasil meminjam buku!');
        ambilBuku(); // Refresh tabel agar stok terbaru (berkurang) langsung terlihat
      } else {
        alert(data.message || 'Gagal meminjam buku');
      }
    } catch {
      alert('Gagal menghubungi server untuk meminjam buku');
    }
  };

  // Utilitas reset form
  const resetForm = () => {
    setBukuId(null);
    setJudul('');
    setPenulis('');
    setStok(5);
  };

  if (loading) {
    return <div style={{ textAlign: 'center', marginTop: 50 }}>Memeriksa sesi login...</div>;
  }

  // --- TAMPILAN BELUM LOGIN ---
  if (!user) {
    return (
      <div style={{ maxWidth: 400, margin: '60px auto', padding: 20, fontFamily: 'Arial', border: '1px solid #ccc', borderRadius: 8 }}>
        <h2>{isRegister ? 'Daftar Akun' : 'Login Sistem (Cookie Auth)'}</h2>
        <form onSubmit={handleAuth}>
          <div style={{ marginBottom: 10 }}>
            <label>Username</label>
            <input
              type="text"
              required
              value={usernameInput}
              onChange={(e) => setUsernameInput(e.target.value)}
              style={{ width: '100%', padding: 8, marginTop: 4, boxSizing: 'border-box' }}
            />
          </div>
          <div style={{ marginBottom: 15 }}>
            <label>Password</label>
            <input
              type="password"
              required
              value={passwordInput}
              onChange={(e) => setPasswordInput(e.target.value)}
              style={{ width: '100%', padding: 8, marginTop: 4, boxSizing: 'border-box' }}
            />
          </div>
          <button type="submit" style={{ width: '100%', padding: 10, background: '#007bff', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>
            {isRegister ? 'Daftar' : 'Masuk'}
          </button>
        </form>

        <button
          onClick={() => setIsRegister(!isRegister)}
          style={{ width: '100%', marginTop: 8, padding: 8, background: '#6c757d', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}
        >
          {isRegister ? 'Sudah punya akun? Login' : 'Belum punya akun? Daftar'}
        </button>

        <hr style={{ margin: '20px 0' }} />

        <a
          href={`${API_URL}/auth/google/login`}
          style={{ display: 'block', textAlign: 'center', padding: 10, background: '#4285F4', color: 'white', textDecoration: 'none', borderRadius: 4, fontWeight: 'bold' }}
        >
          Masuk dengan Google
        </a>
      </div>
    );
  }

  // --- TAMPILAN DASHBOARD BUKU ---
  return (
    <div style={{ maxWidth: 800, margin: '40px auto', padding: 20, fontFamily: 'Arial' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '2px solid #ddd', paddingBottom: 10 }}>
        <h2>Katalog Buku</h2>
        <div>
          <span>Halo, <b>{user}</b></span>
          <button onClick={handleLogout} style={{ marginLeft: 10, padding: '6px 12px', background: '#dc3545', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>
            Logout
          </button>
        </div>
      </div>

      {/* FORM INPUT BUKU */}
      <form onSubmit={handleSimpanBuku} style={{ margin: '20px 0', padding: 15, background: '#f8f9fa', borderRadius: 6 }}>
        <h3>{bukuId ? 'Edit Buku' : 'Tambah Buku Baru'}</h3>
        <div style={{ display: 'flex', gap: '10px', marginBottom: 10 }}>
          <input
            type="text"
            placeholder="Judul Buku"
            required
            value={judul}
            onChange={(e) => setJudul(e.target.value)}
            style={{ flex: 2, padding: 8, boxSizing: 'border-box' }}
          />
          <input
            type="text"
            placeholder="Penulis"
            required
            value={penulis}
            onChange={(e) => setPenulis(e.target.value)}
            style={{ flex: 2, padding: 8, boxSizing: 'border-box' }}
          />
          <input
            type="number"
            placeholder="Stok"
            required
            min="0"
            value={stok}
            onChange={(e) => setStok(e.target.value)}
            style={{ flex: 1, padding: 8, boxSizing: 'border-box' }}
          />
        </div>
        
        <button type="submit" style={{ padding: '8px 16px', background: '#28a745', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>
          {bukuId ? 'Perbarui' : 'Simpan'}
        </button>
        {bukuId && (
          <button type="button" onClick={resetForm} style={{ marginLeft: 8, padding: '8px 16px', background: '#6c757d', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>
            Batal
          </button>
        )}
      </form>

      {/* TABEL BUKU */}
      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead>
          <tr style={{ background: '#eee', textAlign: 'left' }}>
            <th style={{ padding: 8, border: '1px solid #ddd', width: '5%' }}>ID</th>
            <th style={{ padding: 8, border: '1px solid #ddd' }}>Judul</th>
            <th style={{ padding: 8, border: '1px solid #ddd' }}>Penulis</th>
            <th style={{ padding: 8, border: '1px solid #ddd', width: '10%', textAlign: 'center' }}>Stok</th>
            <th style={{ padding: 8, border: '1px solid #ddd', width: '25%', textAlign: 'center' }}>Aksi</th>
          </tr>
        </thead>
        <tbody>
          {bukuList.length === 0 ? (
            <tr>
              <td colSpan="5" style={{ textAlign: 'center', padding: 16 }}>Belum ada buku.</td>
            </tr>
          ) : (
            bukuList.map((b) => (
              <tr key={b.id}>
                <td style={{ padding: 8, border: '1px solid #ddd', textAlign: 'center' }}>{b.id}</td>
                <td style={{ padding: 8, border: '1px solid #ddd' }}>{b.judul}</td>
                <td style={{ padding: 8, border: '1px solid #ddd' }}>{b.penulis}</td>
                <td style={{ padding: 8, border: '1px solid #ddd', textAlign: 'center', fontWeight: 'bold', color: b.stok > 0 ? 'black' : 'red' }}>
                  {b.stok}
                </td>
                <td style={{ padding: 8, border: '1px solid #ddd', textAlign: 'center' }}>
                  
                  {/* TOMBOL PINJAM - LOGIKA STOK */}
                  <button 
                    onClick={() => handlePinjam(b.id)} 
                    disabled={b.stok < 1}
                    style={{ 
                      marginRight: 6, 
                      padding: '4px 8px', 
                      background: b.stok > 0 ? '#17a2b8' : '#e9ecef', 
                      color: b.stok > 0 ? 'white' : '#6c757d',
                      border: 'none', 
                      borderRadius: 4, 
                      cursor: b.stok > 0 ? 'pointer' : 'not-allowed'
                    }}>
                    {b.stok > 0 ? 'Pinjam' : 'Habis'}
                  </button>

                  <button 
                    onClick={() => { setBukuId(b.id); setJudul(b.judul); setPenulis(b.penulis); setStok(b.stok); }} 
                    style={{ marginRight: 6, padding: '4px 8px', background: '#ffc107', border: 'none', borderRadius: 4, cursor: 'pointer', color: 'black' }}>
                    Edit
                  </button>
                  <button 
                    onClick={() => handleHapus(b.id)} 
                    style={{ padding: '4px 8px', background: '#dc3545', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>
                    Hapus
                  </button>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}