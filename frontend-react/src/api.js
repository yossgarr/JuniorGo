// Jika pakai Vite proxy / rewrite Vercel, biarkan '/api'
const API_URL = '/api';

export class ApiError extends Error {
  constructor(message, status) {
    super(message);
    this.status = status;
  }
}

async function request(path, { method = 'GET', body } = {}) {
  let res;
  try {
    res = await fetch(`${API_URL}${path}`, {
      method,
      credentials: 'include',
      headers: body ? { 'Content-Type': 'application/json' } : undefined,
      body: body ? JSON.stringify(body) : undefined,
    });
  } catch {
    throw new ApiError('Tidak bisa menghubungi server', 0);
  }

  let json = null;
  try {
    json = await res.json();
  } catch {
    // respons tanpa body JSON
  }
  if (!res.ok) {
    throw new ApiError(json?.message || `Permintaan gagal (${res.status})`, res.status);
  }
  return json?.data;
}

export const api = {
  me: () => request('/me'),
  login: (username, password) => request('/login', { method: 'POST', body: { username, password } }),
  register: (username, password) => request('/register', { method: 'POST', body: { username, password } }),
  logout: () => request('/logout', { method: 'POST' }),

  ringkasan: () => request('/ringkasan'),

  daftarBarang: () => request('/barang'),
  tambahBarang: (b) => request('/barang', { method: 'POST', body: b }),
  ubahBarang: (id, b) => request(`/barang?id=${id}`, { method: 'PUT', body: b }),
  hapusBarang: (id) => request(`/barang?id=${id}`, { method: 'DELETE' }),

  daftarTransaksi: () => request('/transaksi'),
  catatTransaksi: (t) => request('/transaksi', { method: 'POST', body: t }),
  hapusTransaksi: (id) => request(`/transaksi?id=${id}`, { method: 'DELETE' }),

  googleLoginUrl: `${API_URL}/auth/google/login`,
};
