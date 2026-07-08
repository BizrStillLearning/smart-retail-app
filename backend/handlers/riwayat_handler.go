package handlers

import (
	"net/http"
	"smart-retail-backend/config"
	"smart-retail-backend/models"

	"github.com/gin-gonic/gin"
)

func GetRiwayatPesananPelanggan(c *gin.Context) {
	idPelangganData, exists := c.Get("id_pelanggan")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Otorisasi gagal"})
		return
	}
	idPelanggan := int(idPelangganData.(float64))

	var riwayatPesanan []models.Pesanan

	if err := config.DB.Preload("DetailPesanan").
		Where("id_pelanggan = ?", idPelanggan).
		Order("waktu desc").
		Find(&riwayatPesanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat riwayat pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data":   riwayatPesanan,
	})
}
