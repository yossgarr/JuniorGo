import { useState, useEffect } from 'react';

const API_URL = 'http://localhost:8080';

// --- KOMPONEN KELOLA TOKEN ---
function TokenManager() {
  const [namaToken, setNamaToken] = useState('');
  const [masaBerlaku, setMasaBerlaku] = useState(30);
  const [tokenDibuat, setTokenDibuat] = useState(null);

  const handleBuatToken = async (e) => {
    e.preventDefault();
    try {
      const res = await fetch(`${API_URL}/api/token/create`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          nama_token: namaToken,
          masa_berlaku_hari: Number(masaBerlaku),
        }),
      });

      const hasil = await res.json();
      if (res.ok) {
        setTokenDibuat(hasil.data.token);
        setNamaToken('');
      } else {
        alert(hasil.message || 'Gagal membuat token');
      }
    } catch {
      alert('Terjadi kesalahan jaringan');
    }
  };

  const handleRevoke = async () => {
    if (!tokenDibuat) return;
    try {
      const res = await fetch(`${API_URL}/api/token/revoke`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ token: tokenDibuat }),
      });

      if (res.ok) {
        alert('Token berhasil dicabut!');
        setTokenDibuat(null);
      }
    } catch {
      alert('Gagal mencabut token');
    }
  };

  return (
    <div style={{ background: '#fff', padding: 20, borderRadius: 8, marginTop: 24, border: '1px solid #ddd' }}>
      <h3>Kelola Akses API (Developer Token)</h3>
      <p style={{ fontSize: 13, color: '#666' }}>
        Gunakan token ini untuk menghubungkan script otomasi atau perangkat IoT tanpa perlu login browser.
      </p>

      <form onSubmit={handleBuatToken} style={{ display: 'flex', gap: 10, marginTop: 12 }}>
        <input
          type="text"
          placeholder="Nama Penggunaan (misal: Sensor / Script Python)"
          required
          value={namaToken}
          onChange={(e) => setNamaToken(e.target.value)}
          style={{ flex: 2, padding: 8 }}
        />
        <select
          value={masaBerlaku}
          onChange={(e) => setMasaBerlaku(e.target.value)}
          style={{ flex: 1, padding: 8 }}
        >
          <option value={7}>Aktif 7 Hari</option>
          <option value={30}>Aktif 30 Hari</option>
          <option value={0}>Tanpa Batas Waktu</option>
        </select>
        <button type="submit" style={{ padding: '8px 16px', background: '#007bff', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer' }}>
          Buat Token
        </button>
      </form>

      {tokenDibuat && (
        <div style={{ marginTop: 16, padding: 12, background: '#e8f5e9', border: '1px solid #c8e6c9', borderRadius: 4 }}>
          <b style={{ color: '#2e7d32' }}>Token Akses Berhasil Dibuat!</b>
          <p style={{ margin: '6px 0', fontSize: 12, color: '#555' }}>
            Salin token ini sekarang:
          </p>
          <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
            <input
              type="text"
              readOnly
              value={tokenDibuat}
              style={{ width: '100%', padding: 8, fontFamily: 'monospace', background: '#fff' }}
            />
            <button
              type="button"
              onClick={() => navigator.clipboard.writeText(tokenDibuat)}
              style={{ padding: '8px 12px', cursor: 'pointer' }}
            >
              Salin
            </button>
            <button
              type="button"
              onClick={handleRevoke}
              style={{ padding: '8px 12px', background: '#dc3545', color: '#fff', border: 'none', borderRadius: 4, cursor: 'pointer' }}
            >
              Cabut
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

// --- KOMPONEN UTAMA (WAJIB export default App) ---
export default function App() {
  const [token, setToken] = useState(localStorage.getItem('jwt_token') || '');
  const [user, setUser] = useState(localStorage.getItem('jwt_user') || '');

  const [isRegister, setIsRegister] = useState(false);
  const [usernameInput, setUsernameInput] = useState('');
  const [passwordInput, setPasswordInput] = useState('');

  const [bukuList, setBukuList] = useState([]);
  const [bukuId, setBukuId] = useState(null);
  const [judul, setJudul] = useState('');
  const [penulis, setPenulis] = useState('');

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const tokenUrl = params.get('token');
    const userUrl = params.get('user');

    if (tokenUrl && userUrl) {
      localStorage.setItem('jwt_token', tokenUrl);
      localStorage.setItem('jwt_user', userUrl);
      setToken(tokenUrl);
      setUser(userUrl);
      window.history.replaceState({}, document.title, window.location.pathname);
    }
  }, []);

  useEffect(() => {
    if (token) {
      ambilBuku();
    }
  }, [token]);

  const ambilBuku = async () => {
    try {
      const res = await fetch(`${API_URL}/buku`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      if (res.status === 401) {
        handleLogout();
        return;
      }
      const data = await res.json();
      setBukuList(data.data || []);
    } catch {
      alert('Gagal mengambil data buku');
    }
  };

  const handleAuth = async (e) => {
    e.preventDefault();
    const endpoint = isRegister ? `${API_URL}/register` : `${API_URL}/login`;

    try {
      const res = await fetch(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: usernameInput, password: passwordInput })
      });
      const data = await res.json();

      if (res.ok) {
        if (isRegister) {
          alert('Registrasi berhasil! Silakan login.');
          setIsRegister(false);
        } else {
          localStorage.setItem('jwt_token', data.data.token);
          localStorage.setItem('jwt_user', data.data.username);
          setToken(data.data.token);
          setUser(data.data.username);
        }
      } else {
        alert(data.message || 'Gagal login/register');
      }
    } catch {
      alert('Terjadi kesalahan koneksi');
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('jwt_token');
    localStorage.removeItem('jwt_user');
    setToken('');
    setUser('');
    setBukuList([]);
  };

  const handleSimpanBuku = async (e) => {
    e.preventDefault();
    const method = bukuId ? 'PUT' : 'POST';
    const endpoint = bukuId ? `${API_URL}/buku?id=${bukuId}` : `${API_URL}/buku`;

    try {
      const res = await fetch(endpoint, {
        method,
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`
        },
        body: JSON.stringify({ judul, penulis })
      });

      if (res.ok) {
        setJudul('');
        setPenulis('');
        setBukuId(null);
        ambilBuku();
      } else {
        const data = await res.json();
        alert(data.message || 'Gagal menyimpan buku');
      }
    } catch {
      alert('Gagal menghubungi server');
    }
  };

  const handleHapus = async (id) => {
    if (!confirm('Hapus buku ini?')) return;
    try {
      const res = await fetch(`${API_URL}/buku?id=${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` }
      });
      if (res.ok) ambilBuku();
    } catch {
      alert('Gagal menghapus');
    }
  };

  if (!token) {
    return (
      <div style={{ maxWidth: 400, margin: '60px auto', padding: 20, fontFamily: 'Arial', border: '1px solid #ccc', borderRadius: 8 }}>
        <h2>{isRegister ? 'Daftar Akun' : 'Login Sistem'}</h2>
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

  return (
    <div style={{ maxWidth: 700, margin: '40px auto', padding: 20, fontFamily: 'Arial' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', borderBottom: '2px solid #ddd', paddingBottom: 10 }}>
        <h2>Katalog Buku</h2>
        <div>
          <span>Halo, <b>{user}</b> </span>
          <button onClick={handleLogout} style={{ marginLeft: 10, padding: '6px 12px', background: '#dc3545', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>
            Logout
          </button>
        </div>
      </div>

      <form onSubmit={handleSimpanBuku} style={{ margin: '20px 0', padding: 15, background: '#f8f9fa', borderRadius: 6 }}>
        <h3>{bukuId ? 'Edit Buku' : 'Tambah Buku Baru'}</h3>
        <div style={{ marginBottom: 10 }}>
          <input
            type="text"
            placeholder="Judul Buku"
            required
            value={judul}
            onChange={(e) => setJudul(e.target.value)}
            style={{ width: '100%', padding: 8, boxSizing: 'border-box' }}
          />
        </div>
        <div style={{ marginBottom: 10 }}>
          <input
            type="text"
            placeholder="Penulis"
            required
            value={penulis}
            onChange={(e) => setPenulis(e.target.value)}
            style={{ width: '100%', padding: 8, boxSizing: 'border-box' }}
          />
        </div>
        <button type="submit" style={{ padding: '8px 16px', background: '#28a745', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>
          {bukuId ? 'Perbarui' : 'Simpan'}
        </button>
        {bukuId && (
          <button type="button" onClick={() => { setBukuId(null); setJudul(''); setPenulis(''); }} style={{ marginLeft: 8, padding: '8px 16px', background: '#6c757d', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>
            Batal
          </button>
        )}
      </form>

      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead>
          <tr style={{ background: '#eee', textAlign: 'left' }}>
            <th style={{ padding: 8, border: '1px solid #ddd' }}>ID</th>
            <th style={{ padding: 8, border: '1px solid #ddd' }}>Judul</th>
            <th style={{ padding: 8, border: '1px solid #ddd' }}>Penulis</th>
            <th style={{ padding: 8, border: '1px solid #ddd', width: 140 }}>Aksi</th>
          </tr>
        </thead>
        <tbody>
          {bukuList.length === 0 ? (
            <tr>
              <td colSpan="4" style={{ textAlign: 'center', padding: 16 }}>Belum ada buku.</td>
            </tr>
          ) : (
            bukuList.map((b) => (
              <tr key={b.id}>
                <td style={{ padding: 8, border: '1px solid #ddd' }}>${b.id}</td>
                <td style={{ padding: 8, border: '1px solid #ddd' }}>${b.judul}</td>
                <td style={{ padding: 8, border: '1px solid #ddd' }}>${b.penulis}</td>
                <td style={{ padding: 8, border: '1px solid #ddd' }}>
                  <button onClick={() => { setBukuId(b.id); setJudul(b.judul); setPenulis(b.penulis); }} style={{ marginRight: 6, padding: '4px 8px', background: '#ffc107', border: 'none', borderRadius: 4, cursor: 'pointer' }}>Edit</button>
                  <button onClick={() => handleHapus(b.id)} style={{ padding: '4px 8px', background: '#dc3545', color: 'white', border: 'none', borderRadius: 4, cursor: 'pointer' }}>Hapus</button>
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>

      {/* Bagian Token Manager */}
      <TokenManager />
    </div>
  );
}