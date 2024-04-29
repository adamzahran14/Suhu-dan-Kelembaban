package Config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Dataini adalah struktur data untuk kelembaban dan suhu
type Dataini struct {
	Humidity    float64 `json:"Humidity"`
	Temperature float64 `json:"Temperature"`
	WaktuData   string  `json:"Time"`
}

var DB *sql.DB

// ConnectDB digunakan untuk menghubungkan ke database MySQL
func ConnectDB() error {
	var err error
	DB, err = sql.Open("mysql", "root@/suhu dan kelembaban")
	if err != nil {
		return err
	}

	// Coba ping database untuk memeriksa apakah koneksi berhasil
	err = DB.Ping()
	if err != nil {
		return err
	}

	log.Println("Terhubung ke database MySQL")

	return nil
}

// InsertData digunakan untuk memasukkan data ke tabel database
func InsertData(data Dataini) error {
	// Pastikan sudah terhubung ke database sebelum memasukkan data
	if DB == nil {
		return fmt.Errorf("Tidak terhubung ke database")
	}

	// Query SQL untuk memasukkan data ke tabel
	insertQuery := "INSERT INTO sensor_data (humidity, temperature, waktu_data) VALUES (?, ?, ?)"

	// Eksekusi query dengan parameter
	_, err := DB.Exec(insertQuery, data.Humidity, data.Temperature, data.WaktuData)
	if err != nil {
		return err
	}

	log.Println("Data berhasil dimasukkan ke database.")
	return nil
}

// InsertDataToWeb digunakan untuk mengambil data dari database secara periodik
func InsertDataToWeb() {
	for {
		// Setup database connection
		if DB == nil {
			log.Println("Tidak terhubung ke database. Menunggu koneksi...")
			time.Sleep(5 * time.Second)
			continue
		}

		// Query to fetch temperature, humidity, and time from database
		rows, err := DB.Query("SELECT temperature, humidity, waktu_data FROM sensor_data ORDER BY waktu_data DESC LIMIT 1")
		if err != nil {
			log.Println("Gagal mengambil data dari database:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Iterate through the result set
		for rows.Next() {
			var temperature float64
			var humidity float64
			var waktuDataStr string
			if err := rows.Scan(&temperature, &humidity, &waktuDataStr); err != nil {
				log.Println("Gagal membaca baris data:", err)
				continue
			}

			// Parse waktuDataStr menjadi tipe time.Time
			waktuData, err := time.Parse("2006-01-02 15:04:05", waktuDataStr)
			if err != nil {
				log.Println("Gagal mem-parse waktu:", err)
				continue
			}

			// Konversi ke zona waktu lokal
			waktuDataLocal := waktuData.Local()

			// Gunakan data yang telah dipindai
			log.Printf("Temperature: %.2f, Humidity: %.2f, Time: %s\n", temperature, humidity, waktuDataLocal.Format("2006-01-02 15:04:05"))
		}

		// Tutup baris hasil
		if err := rows.Close(); err != nil {
			log.Println("Gagal menutup baris hasil:", err)
		}

		// Delay to fetch data periodically
		time.Sleep(5 * time.Second)
	}
}
