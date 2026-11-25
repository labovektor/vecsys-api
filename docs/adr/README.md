# Architecture Decision Records (ADR)

Catatan keputusan arsitektur VecSys API. Satu ADR = satu keputusan. ADR yang sudah
diputuskan tidak dihapus; ubah statusnya dan tautkan penggantinya.

## Daftar ADR

| ID | Judul | Status | Tanggal | File |
|----|-------|--------|---------|------|
| ADR-0001 | Batas API vs Repository: Strategi Transaksi / Unit of Work | Accepted | 2026-11-25 | [0001-batas-api-vs-repository-strategi-transaksi.md](./0001-batas-api-vs-repository-strategi-transaksi.md) |

## Konvensi

- **ID:** 4 digit berurutan (`0001`, `0002`, ...).
- **Status:** `Proposed` → `Accepted` / `Rejected`; nanti bisa `Superseded by ADR-XXXX` atau `Deprecated`.
- **Isi wajib:** Konteks, Opsi, Keputusan, Konsekuensi.
- **Fakta kode:** sertakan path file dan nomor baris agar mudah diverifikasi.
- **Bahasa:** Indonesia; istilah teknis boleh tetap Inggris.
- **Membuat ADR baru:** salin struktur `0001-batas-api-vs-repository-strategi-transaksi.md`, ganti ID/judul, lalu tambahkan baris di tabel atas.
