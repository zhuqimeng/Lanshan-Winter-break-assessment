package User

import (
	"gorm.io/gorm"
)

type FeedItem struct {
	gorm.Model
	Username string `json:"username" gorm:"size:64;not null"`
	Content  string `json:"content" gorm:"not null"`
	Author   string `json:"author" gorm:"size:64;not null"`
}
