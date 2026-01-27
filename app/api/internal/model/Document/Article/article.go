package Article

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title string `json:"title"`
	Body  string `json:"body"`
}
