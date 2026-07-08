package handlers

import (
	"net/http"
	"smart-retail-backend/config"
	"smart-retail-backend/models"

	"github.com/gin-gonic/gin"
)

func GetDashboardAdmin(c *gin.Context) {
	var stat struct {
		PendapatanKotor      float64 `json:"pendapatan_kotor"`
		PendapatanBersih     float64 `json:"pendapatan_bersih"`
		TotalDiskonDiberikan float64 `json:"total_diskon_diberikan"`
		PesananSelesai       int64   `json:"pesanan_selesai"`
	}

	config.DB.Table("pesanan").
		Select("COALESCE(SUM(subtotal), 0) as pendapatan_kotor, COALESCE(SUM(total_akhir), 0) as pendapatan_bersih, COALESCE(SUM(total_diskon), 0) as total_diskon_diberikan").
		Where("status = ?", "selesai").
		Scan(&stat)

	config.DB.Table("pesanan").Where("status = ?", "selesai").Count(&stat.PesananSelesai)

	type TopKategori struct {
		Kategori        string  `json:"kategori"`
		TotalPendapatan float64 `json:"total_pendapatan"`
	}
	var topKategori []TopKategori

	config.DB.Table("detail_pesanan dp").
		Select("p.kategori, SUM(dp.harga_satuan * dp.kuantitas) as total_pendapatan").
		Joins("JOIN produk p ON p.id_produk = dp.id_produk").
		Joins("JOIN pesanan ps ON ps.id_pesanan = dp.id_pesanan").
		Where("ps.status = ?", "selesai").
		Group("p.kategori").
		Order("total_pendapatan DESC").
		Limit(5).
		Scan(&topKategori)

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data": gin.H{
			"statistik":    stat,
			"top_kategori": topKategori,
		},
	})
}

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
