package Document

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	Title    string `json:"title" gorm:"size:255;not null"`
	Username string `json:"username" gorm:"size:255;not null"`
	URL      string `json:"url" gorm:"size:255;unique"`
}
