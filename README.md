# Warehouse Management & Order System (BE)

Backend untuk **Warehouse Management & Order System** menggunakan **Golang**, **Gin**, **GORM (PostgreSQL)**, dan Docker untuk kemudahan development & deployment.

---

## 📦 Teknologi & Dependency Utama

- **Golang 1.24**
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **Database Driver**: PostgreSQL (`gorm.io/driver/postgres`)
- **UUID**: `github.com/google/uuid`
- **Environment Loader**: `github.com/joho/godotenv`
- **JWT**: `github.com/dgrijalva/jwt-go`
- **Docker & Docker Compose** untuk menjalankan DB dan service

---

## Mengapa Memilih PostgreSQL

Database **PostgreSQL** dipilih untuk proyek ini karena beberapa alasan:

- **Stabil dan handal**: Cocok untuk aplikasi produksi.
- **Fitur lengkap**: Mendukung SQL standar, JSON, indexing canggih, dan transaksi ACID.
- **Skalabilitas tinggi**: Bisa menangani data besar dan kompleks dengan performa baik.
- **Open source**: Gratis digunakan, dengan komunitas besar dan dokumentasi lengkap.
- **Integrasi mudah**: Didukung oleh banyak framework dan bahasa pemrograman.

PostgreSQL memastikan aplikasi dapat berjalan dengan aman, efisien, dan mudah dikembangkan ke depannya.

---

## 🗂 Struktur Project

```
wms-be/
├─ cmd/ # Main entry point
├─ config/ # Config Golang
├─ database/ # Migrations SQL
├─ docker/ # Dockerfile & docker-compose.yml
├─ domain/ # Models, Repositories, Services
├─ infrastructure/# Database, JWT, Middleware
├─ interfaces/ # HTTP Handlers, Router
├─ go.mod
├─ go.sum
```

---

## ⚡ Persiapan

1. Clone repository:

```bash
git clone https://github.com/febryana0219/wms_be.git
cd wms-be
```

2. Salin file environment:

```
cp .env.example .env
```

4. Sesuaikan .env dengan konfigurasi lokal:

```
DB_HOST=localhost
DB_PORT=55432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=wms_db
JWT_SECRET=your_secret_key
```

## 🐳 Menjalankan dengan Docker

Project sudah termasuk Docker Compose untuk PostgreSQL.

1. Jalankan Docker Compose:

```
docker compose up -d
```

2. Cek container:

```
docker ps
```
- wms_postgres → PostgreSQL

3. Migrasi database

```
psql -U postgres -d wms_db -f /path/to/migrations/001_init_warehouses.up.sql
```

# ulangi untuk file migrations lainnya


## 🏃 Menjalankan Backend Golang

1. Pastikan Golang sudah terinstal (versi >= 1.24):

```
go version
```

2. Jalankan backend:

```
go run cmd/main.go
```
- Server akan berjalan di `ocalhost:8000`.

## 👨‍💻 Author

- Febry
- Email / GitHub: `febryana0219@gmail.com`