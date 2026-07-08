package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"smart-retail-backend/config"
	"smart-retail-backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

func UnggahBuktiTransfer(c *gin.Context) {
	idPelangganData, exists := c.Get("id_pelanggan")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses ditolak"})
		return
	}
	idPelanggan := int(idPelangganData.(float64))

	idPesanan := c.PostForm("id_pesanan")
	if idPesanan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Pesanan wajib disertakan"})
		return
	}

	var pesanan models.Pesanan
	if err := config.DB.Where("id_pesanan = ? AND id_pelanggan = ?", idPesanan, idPelanggan).First(&pesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan atau bukan milik Anda"})
		return
	}

	file, err := c.FormFile("bukti_transfer")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca berkas gambar"})
		return
	}

	ekstensi := filepath.Ext(file.Filename)
	namaFileBaru := fmt.Sprintf("transfer_%s_%d%s", idPesanan, time.Now().Unix(), ekstensi)
	pathSimpan := filepath.Join("uploads", "bukti_transfer", namaFileBaru)

	if err := c.SaveUploadedFile(file, pathSimpan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan berkas ke server"})
		return
	}

	pesanan.BuktiTransfer = "/" + filepath.ToSlash(pathSimpan)
	pesanan.Status = "verifikasi pembayaran"

	if err := config.DB.Save(&pesanan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "sukses",
		"message": "Bukti transfer berhasil diunggah. Menunggu verifikasi admin.",
		"path":    pesanan.BuktiTransfer,
	})
}
