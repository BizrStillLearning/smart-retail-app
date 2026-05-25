package routes

import (
	"smart-retail-backend/handlers"
	"smart-retail-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Static("/uploads", "./uploads")
	api := r.Group("/api")
	{
		api.POST("/login", handlers.Login)
		api.POST("/setup-admin", handlers.SetupAdminAwal)
		api.GET("/produk", handlers.GetSemuaProduk)
		api.POST("/pesanan", handlers.BuatPesanan)
		api.PUT("/pesanan/:id/batal", handlers.BatalkanPesanan)
		api.GET("/laporan/top-kategori", handlers.GetTopKategori)

		adminGroup := api.Group("/admin")
		adminGroup.Use(middlewares.RequireAdmin())
		{
			adminGroup.POST("/produk", handlers.TambahProduk)
			adminGroup.PUT("/produk/:id", handlers.UpdateProduk)
			adminGroup.DELETE("/produk/:id", handlers.HapusProduk)
			adminGroup.GET("/pesanan", handlers.GetAllPesanan)
			adminGroup.PUT("/pesanan/:id/status", handlers.UpdateStatusPesanan)
		}
	}
}
