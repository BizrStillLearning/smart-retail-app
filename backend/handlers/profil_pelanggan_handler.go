package handlers

import (
	"net/http"
	"smart-retail-backend/config"
	"smart-retail-backend/models"

	"github.com/gin-gonic/gin"
)

type UpdateProfilInput struct {
	NamaPelanggan string `json:"nama_pelanggan" binding:"required"`
	Alamat        string `json:"alamat" binding:"required"`
}

func UpdateProfilPelanggan(c *gin.Context) {
	idPelangganData, _ := c.Get("id_pelanggan")
	idPelanggan := int(idPelangganData.(float64))

	var input UpdateProfilInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data profil tidak valid"})
		return
	}

	if err := config.DB.Model(&models.Pelanggan{}).
		Where("id_pelanggan = ?", idPelanggan).
		Updates(map[string]interface{}{
			"nama_pelanggan": input.NamaPelanggan,
			"alamat":         input.Alamat,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal melakukan mutasi data profil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "sukses", "message": "Profil berhasil diperbarui"})
}
