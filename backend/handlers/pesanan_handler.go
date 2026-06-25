package handlers

import (
	"net/http"
	"smart-retail-backend/config"
	"smart-retail-backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
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

	var subtotal float64 = 0

	for _, item := range input.Items {
		var produk models.Produk

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id_produk = ?", item.IDProduk).First(&produk).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Produk dengan ID " + strconv.Itoa(item.IDProduk) + " tidak ditemukan"})
			return
		}

		if produk.Stok < item.Kuantitas {
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{"error": "Stok gagal diproses. Sisa " + produk.NamaProduk + " hanya " + strconv.Itoa(produk.Stok)})
			return
		}

		if err := tx.Model(&produk).Update("stok", produk.Stok-item.Kuantitas).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate sistem inventaris"})
			return
		}

		subtotal += produk.Harga * float64(item.Kuantitas)
	}

	var totalDiskon float64 = 0
	var idPromo *int

	if input.KodePromo != "" {
		var promo models.Promo
		if err := tx.Where("kode_promo = ? AND status = 'aktif'", input.KodePromo).First(&promo).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kode promo tidak valid atau sedang tidak aktif"})
			return
		}

		if subtotal < promo.MinimalBelanja {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Total belanja belum mencapai syarat minimal promo ini"})
			return
		}

		if promo.KhususMember {
			var pelanggan models.Pelanggan
			tx.Where("id_pelanggan = ?", input.IDPelanggan).First(&pelanggan)
			if pelanggan.TipeMember == "Guest" {
				tx.Rollback()
				c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak! Promo ini khusus untuk Member terdaftar."})
				return
			}
		}

		if promo.TipeDiskon == "persentase" {
			totalDiskon = subtotal * (promo.NilaiDiskon / 100)
		} else {
			totalDiskon = promo.NilaiDiskon
		}

		if totalDiskon > subtotal {
			totalDiskon = subtotal
		}

		idPromo = &promo.IDPromo
	}

	totalAkhir := subtotal - totalDiskon

	pesananBaru := models.Pesanan{
		IDPelanggan:  input.IDPelanggan,
		TanggalPesan: time.Now(),
		Status:       "Selesai",
		IDPromo:      idPromo,
		TotalDiskon:  totalDiskon,
		TotalAkhir:   totalAkhir,
	}

	for _, item := range input.Items {
		pesananBaru.Detail = append(pesananBaru.Detail, models.DetailPesanan{
			IDProduk:  item.IDProduk,
			Kuantitas: item.Kuantitas,
		})
	}

	if err := tx.Create(&pesananBaru).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan transaksi: " + err.Error()})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":      "Sukses! Pesanan berhasil dibuat.",
		"id_pesanan":   pesananBaru.IDPesanan,
		"subtotal":     subtotal,
		"total_diskon": totalDiskon,
		"total_akhir":  totalAkhir,
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

	var pesanan models.Pesanan
	if err := tx.Preload("Detail").Where("id_pesanan = ?", idPesanan).First(&pesanan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Data pesanan tidak ditemukan"})
		return
	}

	if pesanan.Status == "Dibatalkan" {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Pesanan ini sudah dibatalkan sebelumnya"})
		return
	}

	for _, item := range pesanan.Detail {
		var produk models.Produk
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id_produk = ?", item.IDProduk).First(&produk).Error; err == nil {
			tx.Model(&produk).Update("stok", produk.Stok+item.Kuantitas)
		}
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat log audit pembatalan"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "sukses",
		"message": "Pesanan berhasil dibatalkan dan stok telah dikembalikan.",
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

func GetTotalTagihanPesanan(c *gin.Context) {
	idPesanan := c.Param("id")
	var total float64

	if err := config.DB.Raw("SELECT HitungTotalPesanan(?)", idPesanan).Scan(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung total tagihan melalui basis data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "sukses",
		"id_pesanan":    idPesanan,
		"total_tagihan": total,
	})
}
