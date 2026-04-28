# Go PM

Aplikasi web project management berbasis Golang untuk mengelola project, ticket, board, roadmap, referential master data, role, dan permission.

## Tech Stack

- Go
- Gin Web Framework
- MySQL
- HTML template
- Tailwind CSS

## Fitur Utama

- Login dan session user
- Manajemen user, role, dan permission
- Master data referential
  - Activities
  - Project statuses
  - Ticket statuses
  - Ticket types
  - Ticket priorities
- Management module
  - Projects
  - Tickets
  - Board
  - Road Map
- Roadmap gantt view per project
- CRUD Epic dan Ticket dari halaman roadmap

## Struktur Proyek

- `main.go` untuk bootstrap aplikasi
- `routes/` untuk definisi routing
- `controllers/` untuk handler HTTP
- `services/` untuk business logic
- `repositories/` untuk akses database
- `models/` untuk struktur data
- `middleware/` untuk auth dan permission
- `templates/` untuk halaman HTML
- `assets/` untuk CSS, JS, font, dan aset statis
- `config/` untuk koneksi dan konfigurasi aplikasi

## Kebutuhan

- Go 1.21+
- MySQL

## Konfigurasi

Buat file `.env` di root project.

Contoh:

```env
APP_PORT=8080

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASS=password
DB_NAME=gobase_app
BASE_URL=http://localhost:8080
```

Import database dari file:

- `gobase_app.sql`

## Menjalankan Aplikasi

```bash
go mod tidy
go run main.go
```

Lalu buka:

```text
http://localhost:8080
```

## Catatan

Project ini merupakan rewrite aplikasi project management ke Golang dengan arsitektur:

- controller
- service
- repository

## Lisensi

Digunakan untuk kebutuhan internal dan pengembangan lanjutan.
