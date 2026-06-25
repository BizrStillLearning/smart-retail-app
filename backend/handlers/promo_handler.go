package handlers

import (
	"net/http"
	"smart-retail-backend/config"
	"smart-retail-backend/models"

	"github.com/gin-gonic/gin"
)

func GetAllPromo(c *gin.Context) {
	var promos []models.Promo

	if err := config.DB.Find(&promos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data promo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data":   promos,
	})
}

func TambahPromo(c *gin.Context) {
	var input models.Promo

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input data promo tidak valid"})
		return
	}

	input.Status = "aktif"

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan promo. Pastikan Kode Promo unik dan belum pernah digunakan."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Promo berhasil ditambahkan",
		"data":    input,
	})
}

func ToggleStatusPromo(c *gin.Context) {
	idPromo := c.Param("id")

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter status tidak valid"})
		return
	}

	if input.Status != "aktif" && input.Status != "nonaktif" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status hanya menerima nilai 'aktif' atau 'nonaktif'"})
		return
	}

	if err := config.DB.Model(&models.Promo{}).Where("id_promo = ?", idPromo).Update("status", input.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status operasional promo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status promo berhasil diperbarui menjadi " + input.Status})
}

func HapusPromo(c *gin.Context) {
	idPromo := c.Param("id")

	if err := config.DB.Where("id_promo = ?", idPromo).Delete(&models.Promo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus promo dari basis data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data promo berhasil dihapus secara permanen"})
}
