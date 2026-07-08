package handlers

import (
	"net/http"
	"smart-retail-backend/config"
	"smart-retail-backend/models"

	"github.com/gin-gonic/gin"
)

func GetPelangganByID(c *gin.Context) {
	id := c.Param("id")
	var pelanggan models.Pelanggan

	if err := config.DB.Where("id_pelanggan = ?", id).First(&pelanggan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ID Member tidak terdaftar di sistem"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data":   pelanggan,
	})
}

func GetAllPelanggan(c *gin.Context) {
	var pelanggans []models.Pelanggan

	if err := config.DB.Order("id_pelanggan desc").Find(&pelanggans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pelanggan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data":   pelanggans,
	})
}

func TambahPelanggan(c *gin.Context) {
	var input models.Pelanggan

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data pelanggan tidak valid"})
		return
	}

	if input.TipeMember == "" {
		input.TipeMember = "Guest"
	}
	input.PoinLoyalitas = 0

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data pelanggan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pelanggan berhasil didaftarkan",
		"data":    input,
	})
}

func UpdatePelanggan(c *gin.Context) {
	id := c.Param("id")
	var input models.Pelanggan

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data pembaruan tidak valid"})
		return
	}

	if err := config.DB.Model(&models.Pelanggan{}).Where("id_pelanggan = ?", id).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui profil pelanggan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profil pelanggan berhasil diperbarui"})
}
