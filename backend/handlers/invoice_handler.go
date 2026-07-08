package handlers

import (
	"fmt"
	"net/http"
	"smart-retail-backend/config"
	"smart-retail-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func GenerateInvoicePDF(c *gin.Context) {
	idPesanan := c.Param("id")
	idPelangganData, _ := c.Get("id_pelanggan")
	idPelanggan := int(idPelangganData.(float64))

	var pesanan models.Pesanan

	if err := config.DB.Preload("Detail.Produk").
		Where("id_pesanan = ? AND id_pelanggan = ?", idPesanan, idPelanggan).
		First(&pesanan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(40, 10, "INVOICE SMART RETAIL")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Order ID: #%d", pesanan.IDPesanan))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Tanggal: %s", pesanan.TanggalPesan.Format("02 Jan 2006 15:04")))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Status: %s", pesanan.Status))
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(80, 10, "Item", "1", 0, "L", false, 0, "")
	pdf.CellFormat(30, 10, "Qty", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Subtotal", "1", 0, "R", false, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 12)
	for _, item := range pesanan.Detail {
		pdf.CellFormat(80, 10, item.Produk.NamaProduk, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", item.Kuantitas), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("Rp %.0f", item.Subtotal), "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}

	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(110, 10, "TOTAL PEMBAYARAN:")
	pdf.CellFormat(40, 10, fmt.Sprintf("Rp %.0f", pesanan.TotalAkhir), "0", 0, "R", false, 0, "")

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=invoice_%d.pdf", pesanan.IDPesanan))
	if err := pdf.Output(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal merender PDF"})
	}
}
