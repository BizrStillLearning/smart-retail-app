package models

type Pelanggan struct {
	IDPelanggan   int    `gorm:"primaryKey;column:id_pelanggan;autoIncrement" json:"id_pelanggan"`
	NamaPelanggan string `gorm:"column:nama_pelanggan" json:"nama_pelanggan"`
	Alamat        string `gorm:"column:alamat" json:"alamat"`
	TipeMember    string `gorm:"column:tipe_member" json:"tipe_member"`
	PoinLoyalitas int    `gorm:"column:poin_loyalitas" json:"poin_loyalitas"`
}

func (Pelanggan) TableName() string {
	return "pelanggan"
}
