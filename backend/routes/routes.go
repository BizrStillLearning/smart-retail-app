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
		api.GET("/produk/:id/rekomendasi", handlers.GetRecommendations)
		api.GET("/pelanggan/:id", handlers.GetPelangganByID)
		adminGroup := api.Group("/admin")
		adminGroup.Use(middlewares.RequireAdmin())
		{
			adminGroup.GET("/dashboard", handlers.GetDashboardAdmin)
			adminGroup.POST("/produk", handlers.TambahProduk)
			adminGroup.PUT("/produk/:id", handlers.UpdateProduk)
			adminGroup.DELETE("/produk/:id", handlers.HapusProduk)
			adminGroup.GET("/pesanan", handlers.GetAllPesanan)
			adminGroup.PUT("/pesanan/:id/status", handlers.UpdateStatusPesanan)
			adminGroup.GET("/pesanan/:id/total", handlers.GetTotalTagihanPesanan)
			adminGroup.POST("/ai/sync", handlers.SyncAndTrainAI)
			adminGroup.GET("/promo", handlers.GetAllPromo)
			adminGroup.POST("/promo", handlers.TambahPromo)
			adminGroup.PUT("/promo/:id/status", handlers.ToggleStatusPromo)
			adminGroup.DELETE("/promo/:id", handlers.HapusPromo)
			adminGroup.GET("/pelanggan", handlers.GetAllPelanggan)
			adminGroup.POST("/pelanggan", handlers.TambahPelanggan)
			adminGroup.PUT("/pelanggan/:id", handlers.UpdatePelanggan)
			adminGroup.GET("/pesanan-online", handlers.GetPesananOnlineAdmin)
			adminGroup.PUT("/pesanan/:id/status", handlers.UpdateStatusPesanan)
			adminGroup.POST("/produk/:id/gambar", handlers.UnggahGambarProduk)
			adminGroup.GET("/statistik", handlers.GetStatistikDashboard)
		}
		ecommerce := api.Group("/pelanggan")
		ecommerce.GET("/produk", handlers.GetKatalogProdukPublik)
		{
			ecommerce.POST("/register", handlers.RegisterPelangganPublik)
			ecommerce.POST("/login", handlers.LoginPelangganPublik)
			ecommerce.GET("/produk/:id", handlers.GetProdukByIDPublik)
			ecommerce.GET("/produk/:id/rekomendasi", handlers.GetRecommendations)
		}
		api.Static("/uploads", "./uploads")
		shopGroup := api.Group("/shop")
		shopGroup.Use(middlewares.RequireCustomerToken())
		{
			shopGroup.POST("/checkout", handlers.ProsesCheckoutECommerce)
			shopGroup.GET("/riwayat", handlers.GetRiwayatPesananPelanggan)
			shopGroup.POST("/upload-bukti", handlers.UnggahBuktiTransfer)
			shopGroup.PUT("/profil", handlers.UpdateProfilPelanggan)
			shopGroup.GET("/pesanan/:id/invoice", handlers.GenerateInvoicePDF)
		}
	}
}
