package models

type Pengguna struct {
	IDPengguna  int    `gorm:"primaryKey;column:id_pengguna;autoIncrement" json:"id_pengguna"`
	Username    string `gorm:"unique;column:username" json:"username"`
	Password    string `gorm:"column:password" json:"-"`
	NamaLengkap string `gorm:"column:nama_lengkap" json:"nama_lengkap"`
	Role        string `gorm:"column:role;default:'pelanggan'" json:"role"`
}

func (Pengguna) TableName() string {
	return "pengguna"
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
