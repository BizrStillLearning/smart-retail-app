package handlers

import (
	"net/http"
	"smart-retail-backend/config"
	"smart-retail-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var customerJwtSecret = []byte("retail_customer_secret_key_2026")

func RegisterPelangganPublik(c *gin.Context) {
	var input models.RegisterPelangganRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format pendaftaran tidak valid. Pastikan nama, email, dan sandi terisi."})
		return
	}

	var cekEmail int64
	config.DB.Model(&models.Pelanggan{}).Where("email = ?", input.Email).Count(&cekEmail)
	if cekEmail > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Alamat email ini sudah terdaftar di sistem."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan saat memproses keamanan sandi."})
		return
	}

	pelangganBaru := models.Pelanggan{
		NamaPelanggan: input.NamaPelanggan,
		Email:         input.Email,
		PasswordHash:  string(hashedPassword),
		TipeMember:    "Guest",
		PoinLoyalitas: 0,
		Alamat:        "-",
	}

	if err := config.DB.Create(&pelangganBaru).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data akun pelanggan."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "sukses",
		"message": "Akun berhasil dibuat. Silakan masuk.",
	})
}

func LoginPelangganPublik(c *gin.Context) {
	var input models.LoginPelangganRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format login tidak valid."})
		return
	}

	var pelanggan models.Pelanggan
	if err := config.DB.Where("email = ?", input.Email).First(&pelanggan).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau kata sandi tidak cocok."})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pelanggan.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau kata sandi tidak cocok."})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id_pelanggan":   pelanggan.IDPelanggan,
		"nama_pelanggan": pelanggan.NamaPelanggan,
		"tipe_member":    pelanggan.TipeMember,
		"exp":            time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(customerJwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menerbitkan token otorisasi."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data": gin.H{
			"token":          tokenString,
			"id_pelanggan":   pelanggan.IDPelanggan,
			"nama_pelanggan": pelanggan.NamaPelanggan,
			"tipe_member":    pelanggan.TipeMember,
		},
	})
}
