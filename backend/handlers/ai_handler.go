package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"smart-retail-backend/config"
	"smart-retail-backend/models"

	"github.com/gin-gonic/gin"
)

type AIItemPayload struct {
	IDPesanan int `json:"id_pesanan"`
	IDProduk  int `json:"id_produk"`
}

type AIResponse struct {
	IDProdukReferensi int   `json:"id_produk_referensi"`
	RekomendasiID     []int `json:"rekomendasi_id"`
}

func SyncAndTrainAI(c *gin.Context) {
	var details []models.DetailPesanan

	if err := config.DB.Find(&details).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data transaksi dari database"})
		return
	}

	var payload []AIItemPayload
	for _, d := range details {
		payload = append(payload, AIItemPayload{
			IDPesanan: d.IDPesanan,
			IDProduk:  d.IDProduk,
		})
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal enkripsi data ke JSON"})
		return
	}

	resp, err := http.Post("http://localhost:8000/api/ai/train", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghubungi AI Service. Pastikan server Python aktif di port 8000"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AI Service menolak data atau volume transaksi belum mencukupi batas support minimum"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "sukses",
		"message": "Sinkronisasi berhasil! Model AI FP-Growth telah diperbarui dengan data transaksi terbaru.",
	})
}

func GetRecommendations(c *gin.Context) {
	idProduk := c.Param("id")

	resp, err := http.Get("http://localhost:8000/api/ai/recommend/" + idProduk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal terhubung dengan layanan analitik AI"})
		return
	}
	defer resp.Body.Close()

	var aiResult AIResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca dokumen output AI Service"})
		return
	}

	if len(aiResult.RekomendasiID) == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "sukses", "data": []models.Produk{}})
		return
	}

	var produkRekomendasi []models.Produk
	if err := config.DB.Where("id_produk IN ?", aiResult.RekomendasiID).Find(&produkRekomendasi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menarik detail entitas produk rekomendasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data":   produkRekomendasi,
	})
}
