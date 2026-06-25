package models

type Produk struct {
	IDProduk   int     `gorm:"primaryKey;column:id_produk" json:"id_produk"`
	NamaProduk string  `gorm:"column:nama_produk" json:"nama_produk"`
	Kategori   string  `gorm:"column:kategori" json:"kategori"`
	Harga      float64 `gorm:"column:harga" json:"harga"`
	Gambar     string  `gorm:"column:gambar" json:"gambar"`
	Stok       int     `gorm:"column:stok" json:"stok"`
}

type TopKategori struct {
	Kategori        string  `json:"kategori"`
	TotalPendapatan float64 `json:"total_pendapatan"`
}

func (Produk) TableName() string {
	return "produk"
}
