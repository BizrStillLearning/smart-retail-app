package models

type Promo struct {
	IDPromo        int     `gorm:"primaryKey;column:id_promo;autoIncrement" json:"id_promo"`
	KodePromo      string  `gorm:"unique;column:kode_promo" json:"kode_promo"`
	NamaPromo      string  `gorm:"column:nama_promo" json:"nama_promo"`
	TipeDiskon     string  `gorm:"column:tipe_diskon" json:"tipe_diskon"`
	NilaiDiskon    float64 `gorm:"column:nilai_diskon" json:"nilai_diskon"`
	MinimalBelanja float64 `gorm:"column:minimal_belanja" json:"minimal_belanja"`
	KhususMember   bool    `gorm:"column:khusus_member" json:"khusus_member"`
	Status         string  `gorm:"column:status" json:"status"`
}

func (Promo) TableName() string {
	return "promo"
}
