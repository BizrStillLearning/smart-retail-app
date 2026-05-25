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

var jwtSecretKey = []byte("rahasia_smart_retail_super_aman")

func Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format input salah"})
		return
	}

	var user models.Pengguna
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id_pengguna": user.IDPengguna,
		"role":        user.Role,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login sukses",
		"token":   tokenString,
		"user": gin.H{
			"nama_lengkap": user.NamaLengkap,
			"role":         user.Role,
		},
	})
}

func SetupAdminAwal(c *gin.Context) {
	config.DB.AutoMigrate(&models.Pengguna{})

	var count int64
	config.DB.Model(&models.Pengguna{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Akun admin sudah pernah dibuat!"})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	adminBaru := models.Pengguna{
		Username:    "admin",
		Password:    string(hashPassword),
		NamaLengkap: "Administrator Toko",
		Role:        "admin",
	}

	config.DB.Create(&adminBaru)

	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses! Akun admin berhasil dibuat.",
		"akun": gin.H{
			"username": "admin",
			"password": "admin123",
		},
	})
}
