package handlers

import (
	"net/http"
	"smart-retail-backend/config"
	"smart-retail-backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func BuatPesanan(c *gin.Context) {
	var input models.CheckoutInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah"})
		return
	}

	if input.IDPelanggan == 0 {
		input.IDPelanggan = 1
	}

	tx := config.DB.Begin()

	pesananBaru := models.Pesanan{
		IDPelanggan:  input.IDPelanggan,
		TanggalPesan: time.Now(),
		Status:       "Selesai",
	}

	for _, item := range input.Items {
		pesananBaru.Detail = append(pesananBaru.Detail, models.DetailPesanan{
			IDProduk:  item.IDProduk,
			Kuantitas: item.Kuantitas,
		})
	}

	if err := tx.Create(&pesananBaru).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan: " + err.Error()})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":    "Sukses! Pesanan berhasil dibuat.",
		"id_pesanan": pesananBaru.IDPesanan,
	})
}

func BatalkanPesanan(c *gin.Context) {
	idPesananStr := c.Param("id")
	idPesanan, err := strconv.Atoi(idPesananStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Pesanan tidak valid"})
		return
	}

	tx := config.DB.Begin()

	if err := tx.Model(&models.Pesanan{}).Where("id_pesanan = ?", idPesanan).Update("status", "Dibatalkan").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membatalkan pesanan"})
		return
	}

	logBatal := models.LogPembatalan{
		IDPesanan:  idPesanan,
		WaktuBatal: time.Now(),
	}
	if err := tx.Create(&logBatal).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat log pembatalan"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "sukses",
		"message": "Pesanan berhasil dibatalkan",
	})
}

func GetAllPesanan(c *gin.Context) {
	var pesanans []models.Pesanan

	if err := config.DB.Preload("Detail.Produk").Order("tanggal_pesan desc").Find(&pesanans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data":   pesanans,
	})
}

func UpdateStatusPesanan(c *gin.Context) {
	idPesanan := c.Param("id")
	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data status tidak valid"})
		return
	}

	if err := config.DB.Model(&models.Pesanan{}).Where("id_pesanan = ?", idPesanan).Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status pesanan berhasil diperbarui menjadi: " + input.Status})
}
