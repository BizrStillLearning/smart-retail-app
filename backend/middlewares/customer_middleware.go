package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireCustomerToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Akses publik ditolak. Harap masuk terlebih dahulu."})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("retail_customer_secret_key_2026"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sesi belanja tidak valid atau telah kedaluwarsa."})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("id_pelanggan", claims["id_pelanggan"])
			c.Set("nama_pelanggan", claims["nama_pelanggan"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Gagal memproses identitas pelanggan."})
			c.Abort()
			return
		}

		c.Next()
	}
}
