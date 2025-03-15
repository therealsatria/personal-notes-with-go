# Personal Notes Frontend

Ini adalah frontend untuk aplikasi Personal Notes. Aplikasi ini merupakan Single Page Application (SPA) yang dibangun dengan JavaScript vanilla, HTML, dan CSS, dengan fokus pada keamanan data melalui enkripsi end-to-end.

## Fitur

- **Enkripsi End-to-End**: Semua data sensitif dienkripsi menggunakan AES-256 dengan encoding Base64
- **Manajemen Kunci**: Antarmuka untuk membangkitkan dan mengelola kunci enkripsi
- **Validasi Keamanan**: Indikator status enkripsi dan pembatasan akses saat enkripsi tidak valid
- **Manajemen Catatan**: Operasi CRUD lengkap untuk catatan dengan prioritas dan tag
- **Manajemen Kategori**: Pengelompokan catatan dalam kategori yang dapat disesuaikan
- **Pencarian**: Pencarian catatan berdasarkan subjek dan konten dengan highlight hasil
- **Pembatasan Data**: Opsi untuk membatasi jumlah catatan yang ditampilkan
- **UI Modern**: Antarmuka pengguna yang bersih dan responsif
- **Navigasi SPA**: Navigasi tanpa reload halaman
- **Notifikasi Toast**: Umpan balik pengguna melalui notifikasi toast
- **Dialog Modal**: Form dan konfirmasi dalam dialog modal
- **Penanganan Error**: Penanganan error yang baik dengan pesan yang informatif

## Struktur

Frontend mengikuti arsitektur berbasis komponen untuk organisasi kode yang bersih:

```
frontend/
├── css/
│   └── styles.css          # Semua style untuk aplikasi
├── js/
│   ├── components/         # Komponen UI
│   │   ├── notes.js        # Komponen catatan
│   │   ├── categories.js   # Komponen kategori
│   │   └── key-generator.js # Komponen pembangkit kunci
│   ├── services/           # Layanan bersama
│   │   ├── api.js          # Layanan komunikasi API
│   │   ├── toast.js        # Layanan notifikasi toast
│   │   └── encryption-status.js # Layanan status enkripsi
│   └── app.js              # Logika aplikasi utama
└── index.html              # File HTML utama
```

## Cara Menggunakan

1. Jalankan server backend:
   ```
   go run main.go
   ```

2. Akses aplikasi:
   - Buka browser dan navigasikan ke `http://localhost:8080/frontend`
   - Server akan otomatis mengarahkan Anda ke frontend

3. Menggunakan aplikasi:
   - Beralih antara Notes dan Categories menggunakan menu navigasi
   - Tambahkan catatan/kategori baru menggunakan tombol "Add"
   - Edit atau hapus item yang ada menggunakan tombol pada setiap kartu
   - Gunakan fitur pencarian untuk menemukan catatan berdasarkan subjek atau konten
   - Gunakan tombol "Show All" untuk menampilkan semua catatan (tanpa batasan)
   - Gunakan tombol floating key generator untuk membuat kunci enkripsi baru
   - Perhatikan banner status enkripsi jika ada masalah dengan sistem enkripsi

## Fitur Keamanan

1. **Validasi Status Enkripsi**:
   - Aplikasi memeriksa status enkripsi saat startup
   - Banner peringatan ditampilkan jika enkripsi tidak valid
   - Tombol "Add Note" dan "Add Category" dinonaktifkan jika enkripsi tidak valid

2. **Pembangkit Kunci**:
   - Tombol floating untuk mengakses pembangkit kunci
   - Kunci baru dapat dibangkitkan berdasarkan teks input
   - Kunci dapat disalin dan diterapkan ke file `settings.json`

3. **Penanganan Error Enkripsi**:
   - Pesan error yang jelas saat operasi enkripsi/dekripsi gagal
   - Opsi untuk mencoba kembali operasi yang gagal

## Pengembangan

Jika Anda ingin memodifikasi frontend:

1. Edit file HTML, CSS, atau JavaScript sesuai kebutuhan
2. Refresh browser Anda untuk melihat perubahan
3. Tidak diperlukan langkah build karena ini adalah aplikasi JavaScript vanilla

## Prinsip Clean Code

Implementasi frontend ini mengikuti prinsip clean code berikut:

- **Single Responsibility Principle**: Setiap komponen dan layanan memiliki tanggung jawab tunggal
- **DRY (Don't Repeat Yourself)**: Fungsionalitas umum diekstrak ke dalam layanan yang dapat digunakan kembali
- **Separation of Concerns**: UI, logika bisnis, dan akses data dipisahkan
- **Consistent Error Handling**: Semua error ditangkap dan ditampilkan kepada pengguna
- **Defensive Programming**: Validasi input dan pemeriksaan error di seluruh aplikasi
- **Meaningful Names**: Variabel, fungsi, dan kelas memiliki nama yang jelas dan deskriptif
- **Comments and Documentation**: Kode didokumentasikan dengan baik dengan komentar gaya JSDoc

## Integrasi dengan Backend

Frontend berkomunikasi dengan backend Go melalui API RESTful. Semua operasi data dilakukan melalui endpoint API yang didokumentasikan dalam README utama proyek. Komunikasi menggunakan format JSON dan menangani enkripsi/dekripsi data sensitif. 