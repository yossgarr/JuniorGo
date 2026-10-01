import { formatRupiah } from '../format';

export default function Ringkasan({ data }) {
  const saldo = data?.saldo ?? 0;
  return (
    <div className="summary">
      <div className="card saldo">
        <span className="label">Saldo</span>
        <strong className={saldo < 0 ? 'neg' : ''}>{formatRupiah(saldo)}</strong>
      </div>
      <div className="card">
        <span className="label">Pemasukan</span>
        <strong className="pos">{formatRupiah(data?.total_pemasukan ?? 0)}</strong>
      </div>
      <div className="card">
        <span className="label">Pengeluaran</span>
        <strong className="neg">{formatRupiah(data?.total_pengeluaran ?? 0)}</strong>
      </div>
    </div>
  );
}
