package handlers

import (
	"net/http"
	"smart-retail-backend/config"
	"smart-retail-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

type ItemKeranjang struct {
	ID        int     `json:"id" binding:"required"`
	Kuantitas int     `json:"kuantitas" binding:"required"`
	Harga     float64 `json:"harga" binding:"required"`
}

type CheckoutRequest struct {
	Items            []ItemKeranjang `json:"items" binding:"required"`
	MetodePembayaran string          `json:"metode_pembayaran" binding:"required"`
}

func ProsesCheckoutECommerce(c *gin.Context) {
	idPelangganData, exists := c.Get("id_pelanggan")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Identitas pelanggan tidak ditemukan"})
		return
	}
	idPelanggan := int(idPelangganData.(float64))

	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data keranjang tidak valid"})
		return
	}

	var totalHarga float64
	for _, item := range req.Items {
		totalHarga += (item.Harga * float64(item.Kuantitas))
	}

	statusAwal := "menunggu pembayaran"
	if req.MetodePembayaran == "COD" {
		statusAwal = "diproses"
	}

	tx := config.DB.Begin()

	pesananBaru := models.Pesanan{
		IDPelanggan:      idPelanggan,
		TanggalPesan:     time.Now(),
		TotalAkhir:       totalHarga,
		TotalDiskon:      0,
		Status:           statusAwal,
		MetodePembayaran: req.MetodePembayaran,
	}

	if err := tx.Create(&pesananBaru).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pesanan induk"})
		return
	}

	for _, item := range req.Items {
		detail := models.DetailPesanan{
			IDPesanan: pesananBaru.IDPesanan,
			IDProduk:  item.ID,
			Kuantitas: item.Kuantitas,
		}
		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan rincian produk"})
			return
		}

		if err := tx.Exec("UPDATE produk SET stok = stok - ? WHERE id_produk = ?", item.Kuantitas, item.ID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memotong stok persediaan"})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":     "sukses",
		"message":    "Pesanan berhasil dibuat",
		"id_pesanan": pesananBaru.IDPesanan,
	})
}
