# Personal Notes with Go

Aplikasi Personal Notes adalah aplikasi manajemen catatan pribadi yang aman dengan enkripsi end-to-end. Aplikasi ini memungkinkan pengguna untuk membuat, membaca, memperbarui, dan menghapus catatan pribadi serta mengorganisirnya dalam kategori.

## Fitur Utama

- **End-to-End Encryption**: Semua data sensitif dienkripsi menggunakan AES-256 dan Base64 encoding.
- **Key Management**: Antarmuka untuk menghasilkan dan mengelola kunci enkripsi.
- **Validasi Keamanan**: Indikator status enkripsi dan pembatasan akses.
- **Manajemen Catatan**: Operasi CRUD lengkap untuk catatan dengan prioritas dan tag.
- **Manajemen Kategori**: Organisasi catatan berdasarkan kategori.
- **Pencarian**: Kemampuan mencari catatan berdasarkan subjek dan konten.
- **Pembatasan Data**: Opsi untuk membatasi jumlah catatan yang ditampilkan.
- **Activity Logging**: Pencatatan semua aktivitas sistem dengan timestamp dan informasi klien.
- **UI Responsif**: Antarmuka pengguna modern yang bekerja di berbagai perangkat.
- **Notifikasi Toast**: Umpan balik pengguna melalui notifikasi toast.
- **Dialog Modal**: Form dan konfirmasi menggunakan dialog modal.
- **Penanganan Error**: Pesan error yang informatif dan penanganan kesalahan yang baik.

## Teknologi

- **Backend**: Go (Golang) dengan Gin framework
- **Database**: SQLite
- **Frontend**: HTML, CSS, JavaScript (Vanilla)
- **Enkripsi**: AES-256 dengan GCM mode
- **Encoding**: Base64

## Struktur Proyek

```
personal-notes-with-go/
├── database/
│   └── db.go                  # Inisialisasi database dan pembuatan tabel
├── frontend/                  # Aplikasi frontend
│   ├── css/
│   │   └── styles.css         # Semua style untuk aplikasi
│   ├── js/
│   │   ├── components/        # Komponen UI
│   │   │   ├── notes.js       # Komponen catatan
│   │   │   ├── categories.js  # Komponen kategori
│   │   │   ├── activity-logs.js # Komponen log aktivitas
│   │   │   └── key-generator.js # Komponen pembangkit kunci
│   │   ├── services/          # Layanan bersama
│   │   │   ├── api.js         # Layanan komunikasi API
│   │   │   ├── toast.js       # Layanan notifikasi toast
│   │   │   └── encryption-status.js # Layanan status enkripsi
│   │   └── app.js             # Logika aplikasi utama
│   └── index.html             # File HTML utama
├── handlers/
│   ├── activity_log_handler.go # Handler untuk log aktivitas
│   ├── category_handler.go    # Handler untuk kategori
│   ├── encryption_handler.go  # Handler untuk status enkripsi
│   ├── key_handler.go         # Handler untuk generasi kunci
│   └── note_handler.go        # Handler untuk catatan
├── models/
│   ├── activity_log.go        # Model untuk log aktivitas
│   ├── category.go            # Model untuk kategori
│   └── note.go                # Model untuk catatan
├── repositories/
│   ├── activity_log_repository.go # Repository untuk log aktivitas
│   ├── category_repository.go # Repository untuk kategori
│   └── note_repository.go     # Repository untuk catatan
├── settings/
│   └── settings.go            # Pengaturan aplikasi
├── utils/
│   ├── encryption.go          # Utilitas enkripsi
│   └── errors.go              # Penanganan error
├── main.go                    # Entry point aplikasi
├── go.mod                     # Dependensi Go
├── go.sum                     # Checksum dependensi
└── settings.json              # File konfigurasi
```

## Endpoint API

### Encryption Status

- **GET /encryption/status**: Mendapatkan status enkripsi saat ini
  - Response: `{"encryption_valid": true|false, "message": "..."}`

### Notes

- **GET /notes**: Mendapatkan semua catatan
  - Query Parameters:
    - `category_id`: Filter berdasarkan kategori
    - `all`: Jika "true", tampilkan semua catatan tanpa batasan
    - `limit`: Jumlah maksimum catatan yang dikembalikan
    - `q`: Query pencarian untuk subjek dan konten
  - Response: Array dari objek Note

- **POST /notes**: Membuat catatan baru
  - Request Body: `{"subject": "...", "content": "...", "priority": "...", "tags": "...", "category_id": "..."}`
  - Response: Objek Note yang dibuat

- **PUT /notes/:id**: Memperbarui catatan yang ada
  - Request Body: `{"subject": "...", "content": "...", "priority": "...", "tags": "...", "category_id": "..."}`
  - Response: Objek Note yang diperbarui

- **DELETE /notes/:id**: Menghapus catatan
  - Response: `{"message": "Note deleted successfully"}`

### Categories

- **GET /categories**: Mendapatkan semua kategori
  - Response: Array dari objek Category

- **POST /categories**: Membuat kategori baru
  - Request Body: `{"name": "..."}`
  - Response: Objek Category yang dibuat

- **PUT /categories/:id**: Memperbarui kategori yang ada
  - Request Body: `{"name": "..."}`
  - Response: Objek Category yang diperbarui

- **DELETE /categories/:id**: Menghapus kategori
  - Response: `{"message": "Category deleted successfully"}`

### Key Generation

- **POST /generate-key**: Menghasilkan kunci enkripsi dari teks
  - Request Body: `{"text": "..."}`
  - Response: `{"key": "..."}`

### Activity Logs

- **GET /activity-logs**: Mendapatkan semua log aktivitas dengan pagination
  - Query Parameters:
    - `limit`: Jumlah maksimum log yang dikembalikan (default: 20)
    - `offset`: Offset untuk pagination (default: 0)
  - Response: Array dari objek ActivityLog

- **GET /activity-logs/count**: Mendapatkan jumlah total log aktivitas
  - Response: `{"count": 123}`

- **GET /activity-logs/entity-type/:entityType**: Mendapatkan log aktivitas berdasarkan tipe entitas
  - Path Parameters:
    - `entityType`: Tipe entitas (misalnya "note", "category", "encryption", "key")
  - Query Parameters:
    - `limit`: Jumlah maksimum log yang dikembalikan (default: 20)
    - `offset`: Offset untuk pagination (default: 0)
  - Response: Array dari objek ActivityLog

- **GET /activity-logs/entity-type/:entityType/count**: Mendapatkan jumlah log aktivitas berdasarkan tipe entitas
  - Path Parameters:
    - `entityType`: Tipe entitas (misalnya "note", "category", "encryption", "key")
  - Response: `{"count": 123}`

- **GET /activity-logs/action/:action**: Mendapatkan log aktivitas berdasarkan aksi
  - Path Parameters:
    - `action`: Tipe aksi (misalnya "create", "update", "delete", "read", "check", "generate")
  - Query Parameters:
    - `limit`: Jumlah maksimum log yang dikembalikan (default: 20)
    - `offset`: Offset untuk pagination (default: 0)
  - Response: Array dari objek ActivityLog

- **GET /activity-logs/action/:action/count**: Mendapatkan jumlah log aktivitas berdasarkan aksi
  - Path Parameters:
    - `action`: Tipe aksi (misalnya "create", "update", "delete", "read", "check", "generate")
  - Response: `{"count": 123}`

- **DELETE /activity-logs/older-than/:days**: Menghapus log aktivitas yang lebih lama dari jumlah hari tertentu
  - Path Parameters:
    - `days`: Jumlah hari
  - Response: `{"message": "Old activity logs deleted successfully", "rowsAffected": 123}`

## Konfigurasi

Aplikasi menggunakan file `settings.json` untuk menyimpan konfigurasi:

```json
{
  "encryption_key": "Base64EncodedKey==",
  "notes_limit": 10
}
```

- **encryption_key**: Kunci enkripsi dalam format Base64
- **notes_limit**: Jumlah maksimum catatan yang ditampilkan secara default

## Fitur Keamanan

### Enkripsi Data Sensitif

Aplikasi menggunakan enkripsi AES-256 untuk mengamankan data sensitif seperti subjek dan konten catatan. Kunci enkripsi disimpan dalam file konfigurasi dan dapat dihasilkan menggunakan endpoint `/generate-key`.

### Validasi Kunci

Sistem melakukan validasi kunci enkripsi saat startup. Jika kunci tidak valid atau tidak ada, modifikasi data akan dinonaktifkan untuk alasan keamanan.

### Pembatasan Akses

Endpoint yang memodifikasi data (POST, PUT, DELETE) memerlukan kunci enkripsi yang valid. Jika kunci tidak valid, permintaan akan ditolak dengan kode status 403 Forbidden.

### Generasi Kunci

Aplikasi menyediakan endpoint untuk menghasilkan kunci enkripsi yang konsisten dari teks input. Kunci dihasilkan menggunakan SHA-256 dan dikodekan dengan Base64.

### Pencatatan Aktivitas

Semua aktivitas sistem dicatat dengan detail seperti jenis aksi, entitas yang terpengaruh, deskripsi, timestamp, alamat IP, dan user agent. Log aktivitas dapat diakses melalui endpoint API dan dapat difilter berdasarkan berbagai kriteria.

## Cara Menjalankan

1. **Clone repository**:
   ```bash
   git clone <repository-url>
   cd personal-notes-with-go
   ```

2. **Jalankan aplikasi**:
   ```bash
   go run main.go
   ```

3. **Akses aplikasi**:
   - Browser akan terbuka otomatis di `http://localhost:8080/frontend`
   - Jika browser tidak terbuka otomatis, buka browser dan akses URL tersebut

4. **Menonaktifkan pembukaan browser otomatis**:
   ```bash
   NO_BROWSER=1 go run main.go
   ```

## Pengujian dengan Curl

### Status Enkripsi

```bash
curl http://localhost:8080/encryption/status
```

### Catatan

```bash
# Membuat catatan baru
curl -X POST -H "Content-Type: application/json" -d '{"subject":"Catatan Baru","content":"Isi catatan","priority":"medium","tags":"tag1, tag2"}' http://localhost:8080/notes

# Mengambil semua catatan
curl http://localhost:8080/notes

# Mengambil catatan berdasarkan kategori
curl http://localhost:8080/notes?category_id=123

# Mengambil semua catatan tanpa batasan
curl http://localhost:8080/notes?all=true

# Memperbarui catatan
curl -X PUT -H "Content-Type: application/json" -d '{"subject":"Catatan Diperbarui","content":"Isi diperbarui","priority":"high","tags":"tag1, tag2, tag3"}' http://localhost:8080/notes/{id}

# Menghapus catatan
curl -X DELETE http://localhost:8080/notes/{id}
```

### Kategori

```bash
# Membuat kategori baru
curl -X POST -H "Content-Type: application/json" -d '{"name":"Kategori Baru"}' http://localhost:8080/categories

# Mengambil semua kategori
curl http://localhost:8080/categories

# Memperbarui kategori
curl -X PUT -H "Content-Type: application/json" -d '{"name":"Kategori Diperbarui"}' http://localhost:8080/categories/{id}

# Menghapus kategori
curl -X DELETE http://localhost:8080/categories/{id}
```

### Pembangkit Kunci

```bash
curl -X POST -H "Content-Type: application/json" -d '{"text":"Teks untuk membangkitkan kunci"}' http://localhost:8080/generate-key
```

### Log Aktivitas

```bash
# Mendapatkan semua log aktivitas
curl http://localhost:8080/activity-logs

# Mendapatkan jumlah total log aktivitas
curl http://localhost:8080/activity-logs/count

# Mendapatkan log aktivitas berdasarkan tipe entitas
curl http://localhost:8080/activity-logs/entity-type/note

# Mendapatkan jumlah log aktivitas berdasarkan tipe entitas
curl http://localhost:8080/activity-logs/entity-type/note/count

# Mendapatkan log aktivitas berdasarkan aksi
curl http://localhost:8080/activity-logs/action/create

# Mendapatkan jumlah log aktivitas berdasarkan aksi
curl http://localhost:8080/activity-logs/action/create/count

# Menghapus log aktivitas yang lebih lama dari 30 hari
curl -X DELETE http://localhost:8080/activity-logs/older-than/30
```

## Fitur Frontend

### Halaman Utama
- Navigasi SPA antara Notes, Categories, dan Activity Logs
- Indikator status enkripsi dengan pesan informatif
- Tombol floating untuk pembangkit kunci

### Manajemen Catatan
- Daftar catatan dengan tampilan kartu yang informatif
- Form untuk menambah dan mengedit catatan
- Pencarian catatan berdasarkan subjek dan konten
- Filter catatan berdasarkan kategori
- Opsi untuk menampilkan semua catatan tanpa batasan

### Manajemen Kategori
- Daftar kategori dengan opsi edit dan hapus
- Form untuk menambah dan mengedit kategori

### Log Aktivitas
- Daftar log aktivitas dengan informasi lengkap
- Filter berdasarkan tipe entitas dan aksi
- Paginasi untuk navigasi mudah
- Opsi untuk menghapus log lama
- Kontrol filter yang tetap terlihat saat scroll

### Pembangkit Kunci
- Form untuk menghasilkan kunci enkripsi dari teks input
- Opsi untuk menyalin kunci ke clipboard

### Notifikasi
- Notifikasi toast untuk umpan balik operasi
- Pesan error yang informatif
- Konfirmasi untuk operasi penghapusan

## Lisensi

Proyek ini dilisensikan di bawah **No Commercial Use License**.

Ini berarti bahwa meskipun Anda bebas menggunakan, memodifikasi, dan mendistribusikan proyek ini untuk tujuan pribadi atau pendidikan, Anda dilarang keras menjualnya atau menggunakannya untuk keuntungan komersial apa pun.