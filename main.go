package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func getHelloWorld(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "hello wordld")
}

func main() {
	http.HandleFunc("/validasi-tabungan", ValidasiTabungan)
	http.HandleFunc("/get", getHelloWorld)
	fmt.Println("Server berjalan di http://localhost:9393")
	http.ListenAndServe(":9393", nil)
}
