package handlers

import (
	"net/http"
	"smart-retail-backend/config"
	"smart-retail-backend/models"

	"github.com/gin-gonic/gin"
)

func GetStatistikDashboard(c *gin.Context) {
	var totalPesanan int64
	var totalPendapatan float64

	config.DB.Model(&models.Pesanan{}).Where("status != 'dibatalkan'").Count(&totalPesanan)

	config.DB.Model(&models.Pesanan{}).Where("status = 'selesai'").Select("COALESCE(SUM(total_akhir), 0)").Row().Scan(&totalPendapatan)

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data": gin.H{
			"total_pesanan":    totalPesanan,
			"total_pendapatan": totalPendapatan,
		},
	})
}
