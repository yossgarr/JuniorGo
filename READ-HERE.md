library
- jwt               = Menyimpan data klaim identitas pengguna (ID, username, role) tanpa membebani database untuk verifikasi sesi.
- OAuth
- Redis ( Upstah )  = Menyimpan salinan data sementara (caching) di sisi server untuk mempercepat waktu respons dan meringankan beban database utama (Neon/PostgreSQL).
- Cookie based Auth = Menyimpan token autentikasi di dalam browser secara aman melalui atribut HttpOnly untuk menangkal pencurian data via serangan XSS.
- API Token         = Memberikan akses khusus mesin atau otomasi (seperti sensor IoT, skrip cron, atau aplikasi pihak ketiga) tanpa melalui siklus antarmuka login.
- HTTP Cache
- REST
- APIs JSON