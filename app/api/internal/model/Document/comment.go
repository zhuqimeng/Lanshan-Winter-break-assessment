package Document

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Username string `json:"username" gorm:"size:255;not null"`
	Link     string `json:"link" gorm:"not null"`
	Content  string `json:"content" gorm:"type:longtext;not null"`
}
