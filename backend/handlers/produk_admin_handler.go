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

func UnggahGambarProduk(c *gin.Context) {
	idProduk := c.Param("id")

	file, err := c.FormFile("gambar_produk")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca berkas gambar"})
		return
	}

	ekstensi := filepath.Ext(file.Filename)
	namaFileBaru := fmt.Sprintf("produk_%s_%d%s", idProduk, time.Now().Unix(), ekstensi)
	pathSimpan := filepath.Join("uploads", "produk", namaFileBaru)

	if err := c.SaveUploadedFile(file, pathSimpan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan berkas ke server"})
		return
	}

	pathURL := "/" + filepath.ToSlash(pathSimpan)

	if err := config.DB.Model(&models.Produk{}).Where("id_produk = ?", idProduk).Update("gambar", pathURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengikat gambar ke database produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "sukses",
		"message": "Visual produk berhasil diperbarui",
		"gambar":  pathURL,
	})
}
