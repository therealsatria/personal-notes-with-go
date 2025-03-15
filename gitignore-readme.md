# Penggunaan File .gitignore

File `.gitignore` telah dikonfigurasi untuk mengecualikan file-file tertentu dari repositori Git. Ini penting untuk:

1. Melindungi informasi sensitif
2. Menghindari penambahan file-file sementara atau yang dihasilkan secara otomatis
3. Menjaga repositori tetap bersih dan fokus pada kode sumber

## Pengecualian Penting

### File Konfigurasi Sensitif
```
settings.json
```
File `settings.json` berisi kunci enkripsi dan pengaturan aplikasi lainnya. File ini **TIDAK** disertakan dalam repositori untuk alasan keamanan. Setiap pengguna harus membuat file ini sendiri berdasarkan template yang disediakan.

### Template settings.json
Gunakan template berikut untuk membuat file `settings.json` Anda sendiri:
```json
{
  "encryption_key": "YOUR_BASE64_ENCODED_KEY",
  "notes_limit": 10
}
```

Anda dapat menggunakan endpoint `/generate-key` untuk membuat kunci enkripsi baru.

### File Database
```
*.db
*.sqlite
*.sqlite3
*.db-journal
*.sqlite-journal
*.sqlite3-journal
```

File database SQLite dan file jurnal terkait tidak disertakan dalam repositori karena:
1. Berisi data pengguna yang mungkin sensitif
2. Dapat berukuran besar dan sering berubah
3. Harus dibuat secara lokal untuk setiap instalasi

## Pengaturan Awal

Saat pertama kali mengkloning repositori ini, Anda perlu:

1. Membuat file `settings.json` dengan kunci enkripsi yang valid
2. Menjalankan aplikasi untuk membuat database secara otomatis

```bash
# Setelah mengkloning repositori
cd personal-notes-with-go

# Buat file settings.json
echo '{"encryption_key":"YOUR_BASE64_ENCODED_KEY","notes_limit":10}' > settings.json

# Jalankan aplikasi untuk membuat database
go run main.go
```

## Catatan Penting

- **JANGAN** menambahkan file `settings.json` atau file database ke repositori Git
- Jika Anda perlu berbagi pengaturan dengan tim, gunakan file template atau variabel lingkungan
- Selalu buat backup file database Anda secara terpisah dari repositori Git

## File Lain yang Dikecualikan

Selain file konfigurasi dan database, `.gitignore` juga mengecualikan:
- File biner dan objek yang dikompilasi
- File sementara dan cache
- File khusus IDE
- File log
- File backup
- Direktori dependensi seperti `node_modules` 