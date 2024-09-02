# Clone

Clone repository ini terlebih dahulu ke komputer Anda

# Persiapan

Sebelum memulai menjalankan projek ini, Anda perlu beberapa tools seperti:

### 1. Docker

Docker diperlukan untuk membuat database [Postgresql](https://www.postgresql.org/) dan [Redis](https://redis.io/). Install Docker melalui link berikut:
https://www.docker.com/

### 2. Database Docker Image

#### - Postgresql

Postgresql akan digunakan sebagai database utama. Install Postgresql dengan command berikut di terminal:

```
docker pull postgres=16-alpine
```

#### - Redis

Redis digunakan untuk menampung session. Install Redis dengan command berikut di terminal:

```
docker pull redis:7.4-alpine
```

### 3. Env file

Pastikan Anda telah mendapatkan file `.env` (jika belum minta kepada admin) dan letakkan di root folder Anda.

### 4. Make

Untuk memudahkan inisialisasi database, kami telah membuatkan make file. Jika Anda bekerja dengan Windows, Anda perlu menginstall makefile di tutorial berikut (pengguna macbook dan linux tidak perlu):
https://www.technewstoday.com/install-and-use-make-in-windows/

# Inisialisasi Database

Untuk inisialisasi database, jalankan perintah berikut di terminal vscode:

```
make postgres
make redis
make createdb
```

# Mulai Bekerja

Jika Anda sudah melalui semua tahap di atas, seharusnya Anda sudah bisa mulai bekerja.
