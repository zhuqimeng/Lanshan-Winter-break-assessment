package QA

import "gorm.io/gorm"

type Answer struct {
	gorm.Model
	Link        int    `json:"link"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
