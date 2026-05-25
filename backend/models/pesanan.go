package models

import "time"

type Pesanan struct {
	IDPesanan    int             `gorm:"primaryKey;column:id_pesanan;autoIncrement" json:"id_pesanan"`
	IDPelanggan  int             `gorm:"column:id_pelanggan" json:"id_pelanggan"`
	TanggalPesan time.Time       `gorm:"column:tanggal_pesan" json:"tanggal_pesan"`
	Status       string          `gorm:"column:status" json:"status"`
	Detail       []DetailPesanan `gorm:"foreignKey:IDPesanan" json:"detail"`
}

type DetailPesanan struct {
	IDDetail  int    `gorm:"primaryKey;column:id_detail;autoIncrement" json:"id_detail"`
	IDPesanan int    `gorm:"column:id_pesanan" json:"id_pesanan"`
	IDProduk  int    `gorm:"column:id_produk" json:"id_produk"`
	Kuantitas int    `gorm:"column:kuantitas" json:"kuantitas"`
	Produk    Produk `gorm:"foreignKey:IDProduk" json:"produk"`
}

func (Pesanan) TableName() string {
	return "pesanan"
}

func (DetailPesanan) TableName() string {
	return "detail_pesanan"
}

type CheckoutInput struct {
	IDPelanggan int `json:"id_pelanggan"`
	Items       []struct {
		IDProduk  int `json:"id_produk"`
		Kuantitas int `json:"kuantitas"`
	} `json:"items"`
}

type LogPembatalan struct {
	IDLog      int       `gorm:"primaryKey;column:id_log;autoIncrement" json:"id_log"`
	IDPesanan  int       `gorm:"column:id_pesanan" json:"id_pesanan"`
	WaktuBatal time.Time `gorm:"column:waktu_batal" json:"waktu_batal"`
}

func (LogPembatalan) TableName() string {
	return "log_pembatalan"
}
