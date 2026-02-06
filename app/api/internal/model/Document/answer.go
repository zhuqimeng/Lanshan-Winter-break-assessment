package Document

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Answer struct {
	gorm.Model
	Username string `json:"username" gorm:"not null"`
	Link     string `json:"link" gorm:"not null"`
	URL      string `json:"url" gorm:"not null;unique"`
	LikeNum  int    `json:"like_num" gorm:"default:0"`
}

func (a Answer) Print() gin.H {
	return gin.H{
		"id":           a.ID,
		"question_url": a.Link,
		"answer_url":   a.URL,
		"like_num":     a.LikeNum,
		"createdAt":    a.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
