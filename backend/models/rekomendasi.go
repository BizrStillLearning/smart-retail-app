package models

type RekomendasiProduk struct {
	ID              int `gorm:"primaryKey;autoIncrement" json:"id"`
	IDProdukUtama   int `json:"id_produk_utama"`
	IDProdukTerkait int `json:"id_produk_terkait"`
	SkorKorelasi    int `json:"skor_korelasi"`
}

func (RekomendasiProduk) TableName() string {
	return "rekomendasi_produk"
}
