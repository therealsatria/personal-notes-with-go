# Personal Notes with Go

Aplikasi Personal Notes adalah aplikasi manajemen catatan pribadi yang aman dengan enkripsi end-to-end. Aplikasi ini memungkinkan pengguna untuk membuat, membaca, memperbarui, dan menghapus catatan pribadi serta mengorganisirnya dalam kategori.

## Fitur Utama

- **Enkripsi End-to-End**: Semua catatan dienkripsi menggunakan AES-256 dengan encoding Base64
- **Manajemen Kunci Enkripsi**: Kunci enkripsi disimpan dalam file `settings.json` dan dapat diperbarui
- **Validasi Enkripsi**: Sistem validasi untuk memastikan kunci enkripsi berfungsi dengan benar
- **Manajemen Kategori**: Pengelompokan catatan berdasarkan kategori
- **Pencarian**: Pencarian catatan berdasarkan subjek dan konten
- **Prioritas**: Penandaan prioritas catatan (rendah, sedang, tinggi)
- **Tag**: Penambahan tag pada catatan untuk pengorganisasian lebih lanjut
- **Pembatasan Data**: Opsi untuk membatasi jumlah catatan yang ditampilkan
- **Antarmuka Pengguna Responsif**: Tampilan yang responsif untuk berbagai ukuran layar

## Teknologi

- **Backend**: Go (Golang) dengan Gin framework
- **Database**: SQLite
- **Frontend**: HTML, CSS, JavaScript (Vanilla)
- **Enkripsi**: AES-256 dengan GCM mode
- **Encoding**: Base64

## Struktur Proyek

```
personal-notes-with-go/
├── database/         # Inisialisasi dan konfigurasi database
├── frontend/         # Antarmuka pengguna (HTML, CSS, JS)
├── handlers/         # Handler HTTP untuk endpoint API
├── models/           # Model data
├── repositories/     # Akses dan manipulasi data
├── settings/         # Konfigurasi aplikasi dan manajemen kunci
├── utils/            # Utilitas umum (enkripsi, error handling)
├── main.go           # Entry point aplikasi
├── settings.json     # File konfigurasi
└── db.sqlite3        # Database SQLite
```

## Endpoint API

### Status Enkripsi

- **`GET /encryption/status`**
  - **Deskripsi**: Memeriksa status sistem enkripsi
  - **Respons**:
    - `200 OK`: Status enkripsi
    ```json
    {
      "encryption_valid": true,
      "message": "Encryption system is properly initialized and working correctly."
    }
    ```

### Catatan (Notes)

- **`POST /notes`**
  - **Deskripsi**: Membuat catatan baru
  - **Memerlukan Enkripsi Valid**: Ya
  - **Request Body**:
    ```json
    {
      "subject": "Judul Catatan",
      "content": "Isi catatan",
      "priority": "medium",
      "tags": "tag1, tag2",
      "category_id": "id-kategori"
    }
    ```
  - **Respons**:
    - `201 Created`: Catatan berhasil dibuat
    - `400 Bad Request`: Request tidak valid
    - `403 Forbidden`: Enkripsi tidak valid
    - `500 Internal Server Error`: Kesalahan server

- **`GET /notes`**
  - **Deskripsi**: Mengambil semua catatan
  - **Parameter Query**:
    - `q`: Kata kunci pencarian (opsional)
    - `all`: Jika "true", menampilkan semua catatan; jika tidak, terbatas sesuai pengaturan (opsional)
  - **Respons**:
    - `200 OK`: Daftar catatan
    - `500 Internal Server Error`: Kesalahan server

- **`PUT /notes/{id}`**
  - **Deskripsi**: Memperbarui catatan berdasarkan ID
  - **Memerlukan Enkripsi Valid**: Ya
  - **Request Body**:
    ```json
    {
      "subject": "Judul Catatan Diperbarui",
      "content": "Isi catatan diperbarui",
      "priority": "high",
      "tags": "tag1, tag2, tag3",
      "category_id": "id-kategori"
    }
    ```
  - **Respons**:
    - `200 OK`: Catatan berhasil diperbarui
    - `400 Bad Request`: Request tidak valid
    - `403 Forbidden`: Enkripsi tidak valid
    - `404 Not Found`: Catatan tidak ditemukan
    - `500 Internal Server Error`: Kesalahan server

- **`DELETE /notes/{id}`**
  - **Deskripsi**: Menghapus catatan berdasarkan ID
  - **Memerlukan Enkripsi Valid**: Ya
  - **Respons**:
    - `200 OK`: Catatan berhasil dihapus
    - `403 Forbidden`: Enkripsi tidak valid
    - `404 Not Found`: Catatan tidak ditemukan
    - `500 Internal Server Error`: Kesalahan server

### Kategori (Categories)

- **`POST /categories`**
  - **Deskripsi**: Membuat kategori baru
  - **Memerlukan Enkripsi Valid**: Ya
  - **Request Body**:
    ```json
    {
      "name": "Nama Kategori"
    }
    ```
  - **Respons**:
    - `201 Created`: Kategori berhasil dibuat
    - `400 Bad Request`: Request tidak valid
    - `403 Forbidden`: Enkripsi tidak valid
    - `500 Internal Server Error`: Kesalahan server

- **`GET /categories`**
  - **Deskripsi**: Mengambil semua kategori
  - **Respons**:
    - `200 OK`: Daftar kategori
    - `500 Internal Server Error`: Kesalahan server

- **`PUT /categories/{id}`**
  - **Deskripsi**: Memperbarui kategori berdasarkan ID
  - **Memerlukan Enkripsi Valid**: Ya
  - **Request Body**:
    ```json
    {
      "name": "Nama Kategori Diperbarui"
    }
    ```
  - **Respons**:
    - `200 OK`: Kategori berhasil diperbarui
    - `400 Bad Request`: Request tidak valid
    - `403 Forbidden`: Enkripsi tidak valid
    - `404 Not Found`: Kategori tidak ditemukan
    - `500 Internal Server Error`: Kesalahan server

- **`DELETE /categories/{id}`**
  - **Deskripsi**: Menghapus kategori berdasarkan ID
  - **Memerlukan Enkripsi Valid**: Ya
  - **Respons**:
    - `200 OK`: Kategori berhasil dihapus
    - `403 Forbidden`: Enkripsi tidak valid
    - `404 Not Found`: Kategori tidak ditemukan
    - `500 Internal Server Error`: Kesalahan server

### Pembangkit Kunci (Key Generator)

- **`POST /generate-key`**
  - **Deskripsi**: Membangkitkan kunci enkripsi baru berdasarkan teks input
  - **Request Body**:
    ```json
    {
      "text": "Teks untuk membangkitkan kunci"
    }
    ```
  - **Respons**:
    - `200 OK`: Kunci berhasil dibangkitkan
    ```json
    {
      "key": "Base64EncodedKey=="
    }
    ```
    - `400 Bad Request`: Request tidak valid
    - `500 Internal Server Error`: Kesalahan server

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

1. **Enkripsi Data Sensitif**:
   - Konten catatan dan tag dienkripsi menggunakan AES-256 dengan GCM mode
   - Kunci enkripsi disimpan dalam file `settings.json`

2. **Validasi Kunci Enkripsi**:
   - Sistem melakukan validasi kunci enkripsi saat aplikasi dimulai
   - Jika kunci tidak valid, operasi modifikasi data dinonaktifkan

3. **Pembatasan Akses**:
   - Endpoint modifikasi data (POST, PUT, DELETE) dinonaktifkan jika enkripsi tidak valid
   - Tombol "Add Note" dan "Add Category" dinonaktifkan jika enkripsi tidak valid

4. **Pembangkit Kunci**:
   - Fitur untuk membangkitkan kunci enkripsi baru berdasarkan teks input
   - Kunci baru dapat disalin dan diterapkan ke file `settings.json`

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
   - Buka browser dan akses `http://localhost:8080/frontend`

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

# Mengambil semua catatan dengan pencarian
curl http://localhost:8080/notes?q=catatan

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

## Lisensi

Proyek ini dilisensikan di bawah **No Commercial Use License**.

Ini berarti bahwa meskipun Anda bebas menggunakan, memodifikasi, dan mendistribusikan proyek ini untuk tujuan pribadi atau pendidikan, Anda dilarang keras menjualnya atau menggunakannya untuk keuntungan komersial apa pun.

