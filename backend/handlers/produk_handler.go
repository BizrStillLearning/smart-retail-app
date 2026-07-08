package handlers

import (
	"math"
	"net/http"
	"os"
	"path/filepath"
	"smart-retail-backend/config"
	"smart-retail-backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetSemuaProduk(c *gin.Context) {
	search := c.Query("search")
	kategori := c.Query("kategori")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	offset := (page - 1) * limit
	query := config.DB.Model(&models.Produk{})

	if kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}

	if search != "" {
		query = query.Where("nama_produk LIKE ?", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var produks []models.Produk
	query.Limit(limit).Offset(offset).Find(&produks)

	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"data":      produks,
		"total":     total,
		"page":      page,
		"last_page": lastPage,
	})
}

func TambahProduk(c *gin.Context) {
	nama := c.PostForm("nama_produk")
	kategori := c.PostForm("kategori")
	hargaStr := c.PostForm("harga")
	harga, _ := strconv.ParseFloat(hargaStr, 64)

	file, err := c.FormFile("gambar")
	var pathGambar string

	if err == nil {
		os.MkdirAll("./uploads", os.ModePerm)

		filename := strconv.FormatInt(time.Now().Unix(), 10) + "_" + filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, "./uploads/"+filename); err == nil {
			pathGambar = "/uploads/" + filename
		}
	}

	produk := models.Produk{
		NamaProduk: nama,
		Kategori:   kategori,
		Harga:      harga,
		Gambar:     pathGambar,
	}

	if err := config.DB.Create(&produk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk dan gambar berhasil ditambahkan"})
}

func UpdateProduk(c *gin.Context) {
	id := c.Param("id")
	var produk models.Produk

	if err := config.DB.First(&produk, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	if nama := c.PostForm("nama_produk"); nama != "" {
		produk.NamaProduk = nama
	}
	if kategori := c.PostForm("kategori"); kategori != "" {
		produk.Kategori = kategori
	}
	if hargaStr := c.PostForm("harga"); hargaStr != "" {
		harga, _ := strconv.ParseFloat(hargaStr, 64)
		produk.Harga = harga
	}

	file, err := c.FormFile("gambar")
	if err == nil {
		filename := strconv.FormatInt(time.Now().Unix(), 10) + "_" + filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, "./uploads/"+filename); err == nil {
			produk.Gambar = "/uploads/" + filename
		}
	}

	config.DB.Save(&produk)
	c.JSON(http.StatusOK, gin.H{"message": "Data produk berhasil diperbarui"})
}

func HapusProduk(c *gin.Context) {
	id := c.Param("id")
	var produk models.Produk

	if err := config.DB.First(&produk, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	config.DB.Delete(&produk)
	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
}

func GetKatalogProdukPublik(c *gin.Context) {
	var daftarProduk []models.Produk

	if err := config.DB.Where("stok > 0").Find(&daftarProduk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil katalog produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data":   daftarProduk,
	})
}

func GetProdukByIDPublik(c *gin.Context) {
	idProduk := c.Param("id")
	var produk models.Produk

	if err := config.DB.Where("id_produk = ? AND stok > 0", idProduk).First(&produk).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan atau stok habis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "sukses",
		"data":   produk,
	})
}
