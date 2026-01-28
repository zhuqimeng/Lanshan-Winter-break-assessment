package Document

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
}
