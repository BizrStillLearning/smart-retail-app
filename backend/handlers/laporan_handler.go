package handlers

import (
	"net/http"
	"smart-retail-backend/config"
	"smart-retail-backend/models"

	"github.com/gin-gonic/gin"
)

func GetTopKategori(c *gin.Context) {
	var hasil []models.TopKategori
	batasPendapatan := 0

	if err := config.DB.Raw("CALL LihatTopKategori(?)", batasPendapatan).Scan(&hasil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal memuat laporan kategori",
		})
		return
	}

	var totalPesananSelesai int64
	config.DB.Model(&models.Pesanan{}).Where("status = ?", "Selesai").Count(&totalPesananSelesai)

	var totalPesananBatal int64
	config.DB.Model(&models.LogPembatalan{}).Count(&totalPesananBatal)

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data":   hasil,
		"stats": gin.H{
			"pesanan_selesai":    totalPesananSelesai,
			"pesanan_dibatalkan": totalPesananBatal,
		},
	})
}
