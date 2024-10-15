package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidasiTabunganame(t *testing.T) {
	tests := []struct {
		nama          string
		input         int
		expectedLayak bool
		expectedPesan string
	}{
		{
			nama:          "Layak menerima bantuan UKT",
			input:         20000,
			expectedLayak: true,
			expectedPesan: "Anda layak menerima bantuan UKT.",
		},
		{
			nama:          "Tidak layak menerima bantuan UKT",
			input:         40000,
			expectedLayak: false,
			expectedPesan: "Anda tidak layak menerima bantuan UKT karena total tabungan harian melebihi batas Rp5 juta.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.nama, func(t *testing.T) {
			// Buat permintaan JSON dari input integer
			reqBody, _ := json.Marshal(map[string]int{
				"tabungan_harian": tt.input,
			})

			req := httptest.NewRequest(http.MethodPost, "/validasi-tabungan", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(ValidasiTabungan)

			handler.ServeHTTP(rr, req)

			// Check response code
			if status := rr.Code; status != http.StatusOK {
				t.Errorf("status code salah: dapat %v, ingin %v", status, http.StatusOK)
			}

			// Check the response body
			responseBody := rr.Body.String()
			if tt.expectedLayak && !strings.Contains(responseBody, tt.expectedPesan) {
				t.Errorf("Pesan salah: dapat %v, ingin %v", responseBody, tt.expectedPesan)
			}
			if !tt.expectedLayak && !strings.Contains(responseBody, tt.expectedPesan) {
				t.Errorf("Pesan salah: dapat %v, ingin %v", responseBody, tt.expectedPesan)
			}
		})
	}
}
