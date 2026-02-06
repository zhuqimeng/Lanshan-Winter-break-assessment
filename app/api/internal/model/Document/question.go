package Document

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	Title    string `json:"title" gorm:"size:255;not null"`
	Username string `json:"username" gorm:"size:255;not null"`
	URL      string `json:"url" gorm:"size:255;unique"`
	LikeNum  int    `json:"like_num" gorm:"default:0"`
}

func (q Question) Print() gin.H {
	return gin.H{
		"id":         q.ID,
		"title":      q.Title,
		"url":        q.URL,
		"like_num":   q.LikeNum,
		"created_at": q.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
