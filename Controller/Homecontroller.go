package Homecontroller

import (
	"html/template"
	"net/http"
)

// Dataini adalah struktur data untuk kelembaban dan suhu
type Dataini struct {
	Humidity    float64 `json:"Humidity"`
	Temperature float64 `json:"Temperature"`
	WaktuData   string  `json:"Time"`
}

// Welcome adalah handler untuk route "/"
func Welcome(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("views/Home/display.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Your existing code to serve HTML
	if err := temp.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DisplayData adalah handler untuk menampilkan data dari database ke halaman HTML
func DisplayData(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan data dari database
	dataFromDB := GetDataFromDB()

	// Membuat template HTML
	temp, err := template.ParseFiles("views/Home/display.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Menyisipkan data ke template HTML
	if err := temp.Execute(w, map[string]interface{}{
		"Humidity":    dataFromDB.Humidity,
		"Temperature": dataFromDB.Temperature,
		"Time":        dataFromDB.WaktuData,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetDataFromMQTT adalah fungsi placeholder untuk mendapatkan data dari database
func GetDataFromMQTT() Dataini {
	// Gantilah dengan logika untuk mendapatkan data MQTT yang terbaru
	simulatedData := Dataini{
		Humidity:    60,
		Temperature: 28.3,
	}
	return simulatedData
}

// GetDataFromDB adalah fungsi placeholder untuk mendapatkan data dari database
func GetDataFromDB() Dataini {
	// Gantilah dengan logika untuk mendapatkan data dari database
	// Misalnya, ambil data dari variabel global atau panggil fungsi database
	simulatedData := Dataini{
		Humidity:    60,
		Temperature: 28.3,
	}

	return simulatedData
}
