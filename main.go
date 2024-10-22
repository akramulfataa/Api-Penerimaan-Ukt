package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
)

const (
	HariPerBulan = 30
	Bulan        = 6
	UKTMax       = 5000000
)

type Permintaan struct {
	TabunganHarian int `json:"tabungan_harian"`
}

type Respon struct {
	Layak    bool   `json:"layak"`
	Pesan    string `json:"pesan"`
	Tabungan string `json:"tabungan"`
}

var res []Respon

func HitungTotalTabunganHarian(tabunganHarian int) int {
	return tabunganHarian * (Bulan * HariPerBulan)
}

func FormatRupiah(amount int) string {
	return fmt.Sprintf("Rp%d", amount)
}

func ValidasiTabungan(w http.ResponseWriter, r *http.Request) {
	var req Permintaan
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Gagal memproses permintaan", http.StatusBadRequest)
		return
	}

	// Hitung total tabungan selama 6 bulan
	totalTabungan := HitungTotalTabunganHarian(req.TabunganHarian)

	// Buat response
	res := Respon{
		Tabungan: FormatRupiah(totalTabungan),
	}

	// Validasi apakah pengguna layak atau tidak berdasarkan total tabungan
	if totalTabungan >= UKTMax {
		res.Layak = false
		res.Pesan = "Anda tidak layak menerima bantuan UKT karena total tabungan harian melebihi batas Rp5 juta."
	} else {
		res.Layak = true
		res.Pesan = "Anda layak menerima bantuan UKT."
	}

	// Kirimkan response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func getTabungan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		slog.Error("load env error %v\n", err)
		return
	}

	http.HandleFunc("/validasi-tabungan", ValidasiTabungan)
	http.HandleFunc("/get-tabugan", getTabungan)
	fmt.Println("Server berjalan di port 9393")

	http.ListenAndServe(":9393", nil)
}
