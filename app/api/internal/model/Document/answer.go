package Document

import "gorm.io/gorm"

type Answer struct {
	gorm.Model
	Username string `json:"username" gorm:"not null"`
	Link     string `json:"link" gorm:"not null"`
	URL      string `json:"url" gorm:"not null;unique"`
}
