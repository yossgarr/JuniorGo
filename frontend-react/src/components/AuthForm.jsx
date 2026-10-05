import { useState } from 'react';
import { api } from '../api';

export default function AuthForm({ onLogin }) {
  const [isRegister, setIsRegister] = useState(false);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [info, setInfo] = useState('');
  const [busy, setBusy] = useState(false);

  const submit = async (e) => {
    e.preventDefault();
    setError('');
    setInfo('');
    setBusy(true);
    try {
      if (isRegister) {
        await api.register(username, password);
        setIsRegister(false);
        setInfo('Registrasi berhasil, silakan masuk.');
      } else {
        const data = await api.login(username, password);
        onLogin(data.username);
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="auth-card">
      <h1>Manage Money By Yatashimura</h1>
      <p className="muted">{isRegister ? 'Buat akun baru' : 'Masuk untuk mengelola keuanganmu'}</p>

      <form onSubmit={submit}>
        <label>
          Username
          <input required value={username} onChange={(e) => setUsername(e.target.value)} autoComplete="username" />
        </label>
        <label>
          Password
          <input
            required
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            autoComplete={isRegister ? 'new-password' : 'current-password'}
          />
        </label>
        {error && <p className="msg error">{error}</p>}
        {info && <p className="msg ok">{info}</p>}
        <button className="btn primary block" disabled={busy}>
          {isRegister ? 'Daftar' : 'Masuk'}
        </button>
      </form>

      <button className="btn link block" onClick={() => { setIsRegister(!isRegister); setError(''); setInfo(''); }}>
        {isRegister ? 'Sudah punya akun? Masuk' : 'Belum punya akun? Daftar'}
      </button>

      <div className="divider"><span>atau</span></div>
      <a className="btn google block" href={api.googleLoginUrl}>Masuk dengan Google</a>
    </div>
  );
}
