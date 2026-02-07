package Document

import "gorm.io/gorm"

type LikeUrlUser struct {
	gorm.Model
	Username string `json:"username" gorm:"not null"`
	Url      string `json:"url" gorm:"not null"`
	Status   bool   `json:"status" gorm:"default:true"`
}
