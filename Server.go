package main

import (
	"SuhuKelembaban/Config"
	Homecontroller "SuhuKelembaban/controller"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
)

// Struktur data sesuaikan dengan struktur data yang diharapkan dari payload JSON
type Dataini struct {
	Humidity    float64 `json:"Humidity"`
	Temperature float64 `json:"Temperature"`
	WaktuData   string  `json:"Time"`
}

// SaveDataToMySQL menyimpan data ke database MySQL
func SaveDataToMySQL(db *sql.DB, data Dataini) error {
	// Query SQL untuk memasukkan data ke tabel
	insertQuery := "INSERT INTO sensor_data (temperature, humidity, Waktu_data) VALUES (?, ?, ?)"

	// Eksekusi query dengan parameter
	_, err := db.Exec(insertQuery, data.Temperature, data.Humidity, data.WaktuData)
	if err != nil {
		return err
	}

	fmt.Println("Data berhasil dimasukkan ke database.")
	return nil
}

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("Menerima pesan pada topik %s: %s\n", message.Topic(), message.Payload())

	// Parsing data JSON
	var data Dataini
	err := json.Unmarshal(message.Payload(), &data)
	if err != nil {
		fmt.Println("Salah mem-parsing JSON:", err)
		return
	}

	// Mendapatkan waktu saat ini dalam format yang diinginkan
	waktuSekarang := time.Now().Format("15:04:05")
	data.WaktuData = waktuSekarang

	// Simpan data ke MySQL
	err = SaveDataToMySQL(Config.DB, data)
	if err != nil {
		fmt.Println("Error menyimpan data ke MySQL:", err)
		return
	}

}

func main() {
	Config.ConnectDB()
	log.Println("Server running di port 8000")
	http.HandleFunc("/", Homecontroller.Welcome)

	// Konfigurasi broker MQTT dan topik
	brokerURL := "tcp://mqtt.telkomiot.id:1883"
	topic := "v2.0/subs/APP65a6366c14d5733573/DEV65a6369e56fd892462"

	// Set up opsi klien MQTT
	opts := MQTT.NewClientOptions().AddBroker(brokerURL).
		SetClientID("tcp:/mqtt.telkomiot.id:1883/DEV65a6369e56fd892462").
		SetUsername("18d1562d722b11cc").
		SetPassword("18d1562d72362b27")

	// Membuat klien MQTT
	client := MQTT.NewClient(opts)

	// Terhubung ke broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Error saat terhubung ke broker MQTT:", token.Error())
		return
	}

	// Berlangganan ke topik yang ditentukan
	if token := client.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
		fmt.Println("Error saat melakukan langganan ke topik:", token.Error())
		return
	}

	fmt.Printf("Berlangganan ke topik: %s\n", topic)

	// Menjalankan server HTTP di goroutine terpisah
	go func() {
		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			log.Fatal("Error server HTTP:", err)
		}
	}()

	// Menjalankan fungsi insertdatatoweb dalam goroutine terpisah
	go Config.InsertDataToWeb()

	// Menunggu sinyal terminasi untuk memutus koneksi dari broker dengan baik
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals // Menunggu sinyal terminasi

	// Memutus koneksi dari broker
	disconnectTimeout := uint(250 * time.Millisecond / time.Nanosecond)
	client.Disconnect(disconnectTimeout)
}
