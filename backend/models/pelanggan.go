package models

type Pelanggan struct {
	IDPelanggan   int    `gorm:"primaryKey;autoIncrement" json:"id_pelanggan"`
	NamaPelanggan string `json:"nama_pelanggan"`
	Alamat        string `json:"alamat"`
	TipeMember    string `json:"tipe_member"`
	PoinLoyalitas int    `json:"poin_loyalitas"`

	Email         string `gorm:"unique;default:null" json:"email"`
	PasswordHash  string `json:"-"`
}

type LoginPelangganRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterPelangganRequest struct {
	NamaPelanggan string `json:"nama_pelanggan" binding:"required"`
	Email         string `json:"email" binding:"required"`
	Password      string `json:"password" binding:"required"`
}