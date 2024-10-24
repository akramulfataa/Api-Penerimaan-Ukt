package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
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
	ID       int    `json:"id"`
	Layak    bool   `json:"layak"`
	Pesan    string `json:"pesan"`
	Tabungan string `json:"tabungan"`
}

var (
	dataValidasi = make(map[int]Respon)
	mu           sync.Mutex
	idCounter    int
)

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

	mu.Lock()
	idCounter++
	res.ID = idCounter
	dataValidasi[res.ID] = res
	defer mu.Unlock()

	if totalTabungan >= UKTMax {
		res.Layak = false
		res.Pesan = fmt.Sprintf("anda tidak layak menerima bantuan ukt karena total tabungan harian melebihi batas rp5 juta.")
	} else {
		res.Layak = true
		res.Pesan = fmt.Sprintf("Anda layak menerima bantuan UKT.")
	}

	// Kirimkan response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func getTabungan(w http.ResponseWriter, r *http.Request) {

	mu.Lock()
	var semuaHasil []Respon
	for _, res := range dataValidasi {
		semuaHasil = append(semuaHasil, res)
	}
	defer mu.Unlock()
	// Kirimkan seluruh hasil validasi
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(semuaHasil)

}

func deleteTabungan(w http.ResponseWriter, r *http.Request) {
	// Ambil ID pengguna dari query string
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Parameter 'id' diperlukan", http.StatusBadRequest)
		return
	}

	// Konversi ID dari string ke int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Nilai 'id' tidak valid", http.StatusBadRequest)
		return
	}

	// Hapus hasil validasi dari map
	mu.Lock()
	_, ok := dataValidasi[id]
	if ok {
		delete(dataValidasi, id)
	}
	defer mu.Unlock()

	if !ok {
		http.Error(w, "Hasil validasi tidak ditemukan untuk ID tersebut", http.StatusNotFound)
		return
	}

	// Kirimkan respons sukses
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hasil validasi dengan ID %d berhasil dihapus.", id)

}

func main() {

	http.HandleFunc("/tabungan", ValidasiTabungan)
	http.HandleFunc("/tabungans", getTabungan)
	http.HandleFunc("/delete-tabungan", deleteTabungan)
	fmt.Println("Server berjalan di port 9393")

	http.ListenAndServe(":9393", nil)
}
