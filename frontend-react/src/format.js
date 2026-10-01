const rupiah = new Intl.NumberFormat('id-ID', {
  style: 'currency',
  currency: 'IDR',
  maximumFractionDigits: 0,
});

export const formatRupiah = (n) => rupiah.format(n);

// Dengan tanda eksplisit: +Rp5.000 / -Rp2.000
export const formatBertanda = (tipe, n) => `${tipe === 'pemasukan' ? '+' : '-'}${rupiah.format(n)}`;

export const formatWaktu = (iso) =>
  new Date(iso).toLocaleString('id-ID', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' });
