# ADR-0001: Batas API vs Repository — Strategi Transaksi / Unit of Work

| Field | Nilai |
|---|---|
| Status | Accepted |
| Tanggal | 2026-11-25 |
| Batas implementasi | April 2027 |
| Menggantikan | — |
| Digantikan oleh | — |


## 1. Konteks

### 1.1 Ringkasan sistem

VecSys adalah REST API Go untuk registrasi event: peserta mendaftar, memilih
kategori/region, membayar (dengan referal), diverifikasi admin, mengisi biodata,
mengunci data, lalu mengunduh kartu peserta PDF.

Stack: Go 1.26, Fiber v2, GORM v1.31, PostgreSQL, Redis (session).
Tidak ada framework DI; semua dependency dirakit manual di `route.SetupRoute`.

Saat ini **tidak ada service layer**. Controller melakukan semuanya: parsing
request, validasi DTO, orkestrasi beberapa repository, dan — di empat handler —
membuka transaksi database sendiri.

### 1.2 Kondisi kode saat ini

- `*gorm.DB` di-inject ke controller:
  - `internal/rest/rest.go:29` — `New(..., db *gorm.DB, ...)`
  - `internal/rest/route/routes.go:72,77` — `db` diteruskan ke
    `ParticipantAdministrationController` dan `ParticipantDataController`.
- Controller memanggil `db.Begin()`, `Commit()`, `Rollback()` langsung:
  - `internal/rest/controller/public.participant.data.go:80` (`AddInstitution`)
  - `internal/rest/controller/public.participant.data.go:181` (`AddMembers`)
  - `internal/rest/controller/public.participant.administration.controller.go:154` (`ValidateReferal`)
  - `internal/rest/controller/public.participant.administration.controller.go:286` (`Payment`)
- Hanya 4 dari 8 repository punya `WithDB(db *gorm.DB)`:
  - `user` (`user.repository.go:32`), `payment` (`payment.repository.go:25`),
    `institution` (`institution.repository.go:16`), `referal` (`referal.repository.go:18`).
  - `admin`, `event`, `category`, `region` tidak punya.
- Tiga method repository membuka transaksi sendiri:
  - `internal/rest/repository/user/user.repository.impl.go:43` (`BulkAddParticipant`)
  - `internal/rest/repository/user/user.repository.impl.go:56` (`BulkUpdateBiodata`)
  - `internal/rest/repository/user/user.repository.impl.go:69` (`BulkAddBiodata`)
- Kontrak error repository hanya `repository.ErrNotFound` +
  `repository.MapNotFound` (`internal/rest/repository/errors.go`), dipetakan ke
  HTTP oleh `respondError` (`internal/rest/controller/error.go`).
- Pola rollback saat ini memakai `defer func() { if r := recover(); ... }()`.
  Ini hanya menangkap **panic**, bukan `return err` biasa; rollback per-cabang
  ditulis manual.
- Error dari `Begin()` dan `Commit()` tidak diperiksa.

### 1.3 Peta transaksi per use case

| Use case / handler | Agregat yang terlibat | Transaksi dibuka di | Catatan faktual |
|---|---|---|---|
| `AddInstitution` | `Institution`, `Participant` | Controller | Simpan file tidak ada di sini |
| `AddMembers` | `Biodata`, `Participant` | Controller | Upload file `id_card` terjadi **di dalam** transaksi; rollback DB tidak menghapus file |
| `ValidateReferal` | `Referal`, `Payment` | Controller | Pengurangan `seat_available` tanpa row lock; ada satu cabang error yang tidak `Rollback` (`...administration.controller.go:197`) |
| `Payment` | `Payment`, `Participant` | Controller | Setelah `CreatePayment`, alur lanjut memanggil `UpdatePayment` dengan ID nol (`...administration.controller.go:314`) |
| `BulkAddParticipantFromCSV` | `Participant` (banyak) | Repository | Transaksi internal repo; tidak lewat `WithDB` |
| `UpdateParticipantBiodata` (admin) | `Biodata` (banyak) | Repository | Transaksi internal repo |

### 1.4 Masalah yang memicu ADR

1. **Batas layer bocor.** Layer transport (controller) mengimpor GORM dan tahu
   cara membuka/menutup transaksi. Repository dan controller sama-sama bisa
   memulai transaksi.
2. **Boilerplate dan inkonsistensi.** Rollback ditulis manual di tiap cabang;
   minimal satu cabang error tidak melakukan rollback. `Begin`/`Commit` error
   diabaikan.
3. **Risiko nested transaction.** Jika `Bulk*` (yang sudah punya transaksi
   internal) nanti dipanggil dengan `WithDB(tx)`, GORM akan membuat savepoint
   atau transaksi bersarang tanpa aturan yang jelas.
4. **Sulit diuji.** Interface repository mengekspos `*gorm.DB` lewat
   `WithDB`, sehingga unit test transaksi butuh database asli atau sqlmock.
   Saat ini hanya ada satu file test (`repository/errors_test.go`).
5. **Tidak ada error domain.** Selain "not found", kegagalan bisnis (mis.
   referal habis, step peserta tidak valid) tidak punya tipe error, sehingga
   service/controller sulit membedakan 400 vs 409 vs 500 tanpa mengecek string
   atau kondisi di handler.
6. **Efek samping non-DB di dalam transaksi.** `AddMembers` menyimpan file di
   antara operasi DB; rollback tidak mengembalikan file. Perlu aturan eksplisit
   apakah I/O boleh berada di dalam unit of work.
7. **Alur transaksi manual rawan bug.** Bug `Payment` di atas menunjukkan
   orkestrasi multi-agregat yang ditulis tangan sulit dijaga konsisten.

### 1.5 Batasan (constraints)

- Stack tidak berubah: Go, Fiber, GORM, PostgreSQL.
- Kontrak HTTP dan format `dto.APIResponse` tidak boleh berubah.
- Pemetaan `ErrNotFound` → 404 dan logging error tetap dipertahankan.
- Tidak menambah framework DI; konstruksi manual di `SetupRoute` tetap dipakai.
- Repository tetap berupa interface agar bisa dipakai lintas handler.
- Ukuran tim: 2 orang; tenggat refactor: April 2027.
- Batasan lain:
  - Tidak ada breaking ke kontrak API ke frontend

---

## 2. Kriteria keputusan

Skor 1–5; semakin tinggi semakin baik. Skor terbobot: **A 72,0%**, B 68,7%,
C 70,7%, D 56,0%. Skor manual pada tiap opsi bersifat kualitatif dan tidak
dikonversi langsung ke persentase; keputusan akhir menimbang beban kognitif tim.

| Kriteria | Mengapa relevan di VecSys | Bobot | A | B | C | D |
|---|---|---|---|---|---|---|
| Konsistensi lintas agregat | `Payment` + `Referal` + `Participant` harus berubah bersama | 5 | 5 | 5 | 5 | 5 |
| Batas tanggung jawab layer | Controller tetap memegang GORM dan `*gorm.DB` | 3 | 2 | 4 | 4 | 3 |
| Testability tanpa DB | Kontrak repo memakai `*gorm.DB` konkret, sulit di-fake | 2 | 2 | 4 | 3 | 2 |
| Biaya migrasi | 4 handler + seluruh signature repo + call site | 4 | 3 | 2 | 2 | 2 |
| Risiko nested transaction | 3 method repo sudah punya tx internal | 4 | 3 | 3 | 3 | 3 |
| Kejelasan alur (audit/debug) | Alur tx eksplisit di signature vs implisit | 3 | 4 | 4 | 4 | 2 |
| Skalabilitas fitur (retry, lock, idempotency) | Butuh klaim referal aman dan retry | 4 | 3 | 4 | 4 | 2 |
| Beban kognitif tim | Tim kecil, ada yang pemula, codebase baru | 5 | 5 | 2 | 3 | 2 |

---

## 3. Opsi

### Opsi A — Controller-owned transaction dengan `*gorm.DB` eksplisit

**Mekanisme:** controller memegang `*gorm.DB` dan membuka transaksi. Repository
tidak lagi menyimpan/mengganti DB lewat `WithDB`; setiap method menerima
`db *gorm.DB` sebagai parameter, sehingga controller bisa mengoper `tx` yang sama.

```go
tx := c.DB.Begin()
if tx.Error != nil { /* 500 */ }
defer func() { if r := recover(); r != nil { tx.Rollback() } }()

if err := c.CategoryRepo.Create(tx.WithContext(c.Context()), &category); err != nil {
    tx.Rollback()
    return respondError(c, err, ...)
}
if err := tx.Commit().Error; err != nil { /* 500 */ }
```

**Perubahan minimum:** helper transaksi di controller/util; hapus `WithDB` dari
repo yang memilikinya; tambah parameter `db *gorm.DB` pada seluruh method repo;
sesuaikan seluruh call site. Tidak ada package service/UoW baru.

**Dampak yang sudah pasti:**
- Controller tetap mengimpor `gorm.io/gorm`; separation of concern tetap longgar.
- `WithDB(*gorm.DB)` dihapus; kontrak repo berubah di semua domain.
- Semua call site repo berubah, termasuk handler yang tidak memakai transaksi.
- Unit test transaksi tetap butuh DB/sqlmock; integration test lebih realistis.

**Pro:** Tim baru yang masih programmer pemula akan lebih mudah karena minim kompleksitas

**Kontra:** Separation of layer masih terbilang kurang

**Penilaian:** 7/10

---

### Opsi B — Service layer + Unit of Work (callback)

**Mekanisme:** controller hanya memanggil service. Service meminta UoW untuk
menjalankan satu blok kerja; UoW membuka transaksi GORM, menyediakan kumpulan
repository yang terikat ke transaksi, lalu commit/rollback otomatis.

```go
// internal/uow
type UnitOfWork interface {
    Do(ctx context.Context, fn func(ctx context.Context, r Repositories) error) error
}

// internal/service
func (s *paymentService) Pay(ctx context.Context, in PayInput) error {
    return s.uow.Do(ctx, func(ctx context.Context, r uow.Repositories) error {
        p, err := r.Payment.GetByParticipant(ctx, in.ParticipantID)
        if err != nil && !errors.Is(err, repository.ErrNotFound) { return err }
        // ...
        return r.Participant.AdvanceStep(ctx, in.ParticipantID, entity.StepPaidParticipant)
    })
}
```

**Perubahan minimum:** package `internal/service` + `internal/uow`; konstruksi
repository per-transaksi; controller kehilangan `*gorm.DB`; `WithDB` dihapus
atau dijadikan internal UoW.

**Dampak yang sudah pasti:**
- 4 handler pindah orkestrasi ke service.
- `route.SetupRoute` berubah: tidak lagi mengoper `db` ke controller.
- Semua repo yang dipakai lintas agregat harus bisa dibuat dari `*gorm.DB`
  (bukan hanya 4 yang sekarang punya `WithDB`).
- Butuh keputusan desain untuk error domain di dalam callback.

**Pro:** Separation of layer menjadi bagus

**Kontra:** Programmer pemula yang baru masuk akan kebingungan soal alur app

**Penilaian:** 5/10

---

### Opsi C — Repository coarse-grained (satu method = satu use case atomik)

**Mekanisme:** tidak ada transaksi di controller/service. Repository
menyediakan method khusus use case yang menjalankan transaksinya sendiri.

```go
// kontrak repo
SavePaymentAndAdvance(ctx context.Context, in SavePaymentInput) (entity.Payment, error)

// implementasi
func (r *paymentRepositoryImpl) SavePaymentAndAdvance(ctx context.Context, in SavePaymentInput) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // upsert payment + update progress_step peserta
        return nil
    })
}
```

**Perubahan minimum:** method baru di repo `payment`, `referal`,
`institution`, `participant`; controller memanggil satu method per use case;
`WithDB` dan `db` di controller dihapus.

**Dampak yang sudah pasti:**
- Logika lintas agregat pindah ke repository (repo jadi lebih gemuk dan
  use-case-specific).
- Service tetap mungkin diperlukan untuk efek non-DB (upload file, email).
- Duplikasi pola transaksi di beberapa repo, kecuali dibuat helper bersama.

**Pro:** Perubahan Minimum yang mendukung separation of concern yang masih bisa ditolerir programmer pemula

**Kontra:** Semakin banyak atomic usecase, layer ini akan membengkak dan bingung membacanya

**Penilaian:** 6/10

---

### Opsi D — Transaksi dibawa lewat `context.Context`

**Mekanisme:** pemanggil menaruh transaksi di context; repository mengambilnya
secara implisit dan memakai `WithContext`.

```go
ctx, tx, err := uow.Begin(ctx, db)
defer tx.RollbackUnlessCommitted()

if err := paymentRepo.Save(ctx, p); err != nil { return err } // repo: db := uow.From(ctx, r.db)
return tx.Commit().Error
```

**Perubahan minimum:** package `uow` berisi helper context; setiap method repo
yang ikut transaksi harus membaca tx dari context; controller masih memegang
`*gorm.DB` untuk memulai.

**Dampak yang sudah pasti:**
- Signature repo tidak berubah (tetap `ctx context.Context`).
- Dependency transaksi menjadi implisit: sulit terlihat dari signature.
- Risiko lupa membungkus context atau memakai `r.db` langsung alih-alih tx.

**Pro:** -

**Kontra:** Implementasi transactionnya terkesan jadi terpisah-pisah dan kalau terjadi banyak akan membengkakkan code dan juga tetap membuat refactor besar pada code

**Penilaian:** 5/10

---

## 4. Keputusan

**Opsi terpilih:** A

**Alasan utama:** Opsi A dipilih karena alur transaksi paling mudah dipahami
programmer pemula dan tidak menambah service/UoW. Separation of concern yang
kurang diterima sebagai trade-off; keputusan ditinjau ulang bila komposisi atau
skill tim berubah.

### 4.1 Aturan batas API vs Repository

Aturan berikut menjadi checklist review PR.

| # | Pertanyaan keputusan | Keputusan | Alasan |
|---|---|---|---|
| 1 | Layer mana yang boleh membuka `Begin`/`Commit`/`Rollback`? | Controller, kecuali method `Bulk*` yang sudah ada (lihat #6) | Agar programmer baru bisa memahami; pengecualian lama didokumentasikan |
| 2 | Bolehkah controller menerima `*gorm.DB`? | Ya | Agar mengontrol `Begin`/`Commit`/`Rollback` |
| 3 | Apakah `WithDB(*gorm.DB)` dipertahankan, diganti, atau dihapus? | Dihapus; diganti parameter `db *gorm.DB` eksplisit di method repo | DB dan transaksi dikontrol controller |
| 4 | Apakah service layer wajib ada? Di package mana? | Tidak wajib | Menambah kompleksitas |
| 5 | Bentuk kontrak transaksi: callback UoW, context, atau method repo? | Method repo menerima `db *gorm.DB` dari controller | Menyesuaikan Opsi A |
| 6 | Bagaimana memperlakukan `Bulk*` yang sudah punya transaksi internal? | Pertahankan transaksi internalnya | Alur ini tidak dipanggil di dalam transaksi controller, jadi tidak ada nested tx |
| 7 | Di mana error domain didefinisikan dan bagaimana dipetakan ke HTTP? | Diputuskan manual di controller; tanpa error domain baru | Lebih simple |
| 8 | Bolehkah I/O non-DB (upload file, email) berada di dalam transaksi? | Belum diputuskan (open question) | Perlu keputusan terpisah; risiko file yatim dicatat di 5.3 |
| 9 | Level isolasi/locking untuk klaim referal dan pembayaran? | Read Committed + `SELECT ... FOR UPDATE` pada baris referal | Mencegah dua peserta mengambil seat terakhir; pembayaran memakai tx yang sama |
| 10 | Bagaimana transaksi diuji (fake UoW, integration test, sqlmock)? | Integration test PostgreSQL | Lebih mudah dan realistis untuk perilaku rollback |

---

## 5. Konsekuensi

### 5.1 Dampak positif

Alur transaksi tetap sederhana dan mudah dipahami programmer pemula; tidak ada
service/UoW baru. Transaksi lintas agregat tetap bisa dijalankan karena
controller mengoper `tx` yang sama ke repository. Bug rollback dan klaim referal
bisa diperbaiki tanpa menambah layer baru.

### 5.2 Dampak negatif / biaya

Controller tetap mengimpor GORM dan memegang orkestrasi bisnis, sehingga
separation of concern masih kurang. Menghapus `WithDB` dan menambah parameter
`db *gorm.DB` menyentuh seluruh signature repo dan call site, jadi biaya
refactor lebih besar daripada status quo. Testing transaksi tetap bergantung
pada database.

### 5.3 Risiko dan mitigasi

| Risiko | Kemungkinan | Dampak | Mitigasi |
|---|---|---|---|
| Regresi pada 4 handler transaksi saat refactor | Sedang | Tinggi | Integration test sukses + rollback per use case; migrasi satu handler per PR |
| Nested transaction / savepoint tidak terkendali | Rendah | Sedang | `Bulk*` tidak boleh dipanggil di dalam transaksi controller; beri komentar/guard |
| File upload sudah tersimpan saat DB rollback | Sedang | Sedang | Belum diputuskan (#8); sementara simpan file setelah commit atau hapus saat rollback |
| Perubahan kontrak repo memecah call site lain | Tinggi | Tinggi | Ubah signature per domain, compile-driven; jangan campur dengan perubahan perilaku |
| Seat referal tetap race jika `FOR UPDATE` terlupa | Sedang | Tinggi | Wajib row lock saat baca + kurangi seat; test dua klaim seat terakhir bersamaan |
| Controller tetap gemuk dan tahu GORM | Tinggi | Sedang | Batasi blok transaksi; helper; review checklist; tinjau ulang ADR jika handler membesar |
| Mapping error manual tidak konsisten | Sedang | Rendah | Tabel mapping di controller; review PR; error domain bisa jadi ADR terpisah |

### 5.4 Area kode yang terdampak

| Area | File | Perubahan yang kemungkinan diperlukan |
|---|---|---|
| Wiring | `internal/rest/rest.go`, `internal/rest/route/routes.go` | Tetap mengoper `*gorm.DB` ke controller; tidak ada service/UoW |
| Kontrak repo | `internal/rest/repository/{admin,category,event,institution,payment,referal,region,user}/*.repository.go` | Hapus `WithDB`; semua method menerima `db *gorm.DB` |
| Implementasi repo | `.../*.repository.impl.go` | Pakai parameter `db` + `WithContext`, bukan `r.db` |
| Pemanggil repo | `internal/rest/controller/*.go` | Sesuaikan semua call site; non-tx pakai `c.DB.WithContext(ctx)`, tx pakai `tx` |
| Handler transaksi | `public.participant.data.go`, `public.participant.administration.controller.go` | Pakai helper tx; perbaiki rollback; bug `Payment`; `FOR UPDATE` referal |
| Repo transaksi internal | `internal/rest/repository/user/user.repository.impl.go:43,56,69` | Pertahankan; jangan panggil dari dalam transaksi controller |
| Error contract | `internal/rest/repository/errors.go`, `internal/rest/controller/error.go` | Tidak ada error domain baru; mapping tetap di controller |
| Test | `internal/rest/repository/errors_test.go` + test baru | Integration test 4 use case, rollback, dan konkurensi seat |
| Dokumentasi | `docs/adr/README.md` | Status `Accepted` |

### 5.5 Kompatibilitas

- Perubahan pada response HTTP: tidak berubah.
- Perubahan pada signature repository internal: `WithDB` dihapus; semua method
  menerima `db *gorm.DB`. Ini breaking secara compile untuk seluruh call site.
- Perubahan skema database: tidak ada.
- Perubahan perilaku runtime: hanya perbaikan bug (`Payment`, rollback
  `ValidateReferal`) dan penambahan row lock pada klaim referal.

### 5.6 Pertanyaan terbuka / utang teknis

- **I/O di dalam transaksi (#8):** `AddMembers` menyimpan file `id_card` di
  antara operasi DB. Pilih simpan-setelah-commit atau hapus-file-saat-rollback.
- **Error domain:** belum ada tipe error bisnis; mapping manual di controller
  bisa menyimpang. Pertimbangkan ADR terpisah bila handler bertambah.
- **`Bulk*`:** aturan "tidak dipanggil dalam tx" baru berupa konvensi; tambahkan
  komentar/guard agar tidak dilanggar tanpa sengaja.
- **Batas layer dan testability:** controller tetap memegang GORM dan unit test
  transaksi tetap butuh database. Diterima sebagai trade-off; integration test
  menjadi jalur verifikasi utama.
- **Ukuran controller:** Opsi A menerima controller gemuk. Tinjau ulang bila satu
  handler melewati batas yang disepakati tim.

---

## 6. Rencana implementasi

Centang setelah selesai. Urutan mengikuti Opsi A.

- [x] **Fase 0 — Terima keputusan.** Status ADR `Accepted`; index `docs/adr/README.md` diperbarui.
- [ ] **Fase 1 — Helper transaksi.** Tambah `withTransaction` di `internal/rest/controller` (atau `internal/util`) yang mengecek `Begin`, rollback otomatis saat error/panic, mengecek `Commit`, dan meneruskan `tx`.
- [ ] **Fase 2 — Perbaiki bug yang diketahui.** `Payment` tidak boleh lanjut `UpdatePayment` dengan ID nol setelah create; `ValidateReferal` wajib rollback di semua cabang; tambah `SELECT ... FOR UPDATE` untuk baris referal.
- [ ] **Fase 3 — Ubah kontrak repo.** Hapus `WithDB`; semua method menerima `db *gorm.DB`; implementasi memakai `db.WithContext(...)`.
- [ ] **Fase 4 — Migrasi call site.** Mulai dari 4 handler transaksi, lalu handler non-transaksi; gunakan `tx` atau `c.DB.WithContext(ctx)`.
- [ ] **Fase 5 — Aturan `Bulk*` dan I/O.** Tambah guard/komentar bahwa `Bulk*` tidak boleh dipanggil dalam tx controller; putuskan penanganan file (#8).
- [ ] **Fase 6 — Integration test.** Skenario sukses + rollback untuk 4 use case; test dua klaim seat terakhir bersamaan.
- [ ] **Fase 7 — Bersihkan.** Hapus kode mati, perbarui status/index, catat utang teknis.

**Definition of done:**

- [ ] Tidak ada `WithDB` di interface maupun implementasi repo.
- [ ] Semua method repo menerima `db *gorm.DB` dan memakai context dari pemanggil.
- [ ] Semua transaksi controller memakai helper; tidak ada `Rollback` manual di tiap cabang.
- [ ] Error `Begin`/`Commit` diperiksa.
- [ ] Bug `Payment` (update dengan ID nol) dan missing rollback `ValidateReferal` diperbaiki.
- [ ] Klaim referal memakai `SELECT ... FOR UPDATE`.
- [ ] `Bulk*` tidak dipanggil di dalam transaksi controller.
- [ ] Integration test 4 use case + rollback + konkurensi seat lulus.
- [ ] `go build ./...` dan `go test ./...` lulus.
- [ ] _[kriteria tambahan]_

---

## 7. Referensi

- Pola transaksi saat ini:
  [`public.participant.data.go#L80`](../../internal/rest/controller/public.participant.data.go#L80),
  [`public.participant.data.go#L181`](../../internal/rest/controller/public.participant.data.go#L181),
  [`public.participant.administration.controller.go#L154`](../../internal/rest/controller/public.participant.administration.controller.go#L154),
  [`public.participant.administration.controller.go#L197`](../../internal/rest/controller/public.participant.administration.controller.go#L197),
  [`public.participant.administration.controller.go#L286`](../../internal/rest/controller/public.participant.administration.controller.go#L286),
  [`public.participant.administration.controller.go#L314`](../../internal/rest/controller/public.participant.administration.controller.go#L314)
- Injeksi DB ke controller: [`route/routes.go#L72`](../../internal/rest/route/routes.go#L72), [`rest.go#L29`](../../internal/rest/rest.go#L29)
- `WithDB`: [`user.repository.go#L32`](../../internal/rest/repository/user/user.repository.go#L32), [`payment.repository.go#L25`](../../internal/rest/repository/payment/payment.repository.go#L25), [`institution.repository.go#L16`](../../internal/rest/repository/institution/institution.repository.go#L16), [`referal.repository.go#L18`](../../internal/rest/repository/referal/referal.repository.go#L18)
- Transaksi internal repo: [`user.repository.impl.go#L43`](../../internal/rest/repository/user/user.repository.impl.go#L43), [`#L56`](../../internal/rest/repository/user/user.repository.impl.go#L56), [`#L69`](../../internal/rest/repository/user/user.repository.impl.go#L69)
- Kontrak error: [`repository/errors.go`](../../internal/rest/repository/errors.go), [`controller/error.go`](../../internal/rest/controller/error.go)
- State machine peserta: [`entity/participants.entity.go`](../../entity/participants.entity.go) (`ValidateUserStep`)
- Catatan: referensi `WithDB` di atas menggambarkan kondisi **sebelum** ADR ini
  diimplementasikan. Setelah Fase 3, kontraknya menjadi parameter `db *gorm.DB`.
- Index ADR: [`docs/adr/README.md`](./README.md)
