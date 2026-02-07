package Document

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title    string `json:"title" gorm:"size:255;not null"`
	Username string `json:"username" gorm:"size:255;not null"`
	URL      string `json:"url" gorm:"size:255;unique"`
	LikeNum  int    `json:"like_num" gorm:"default:0"`
}

func (a *Article) UpdLike(op bool) {
	if op {
		a.LikeNum++
	} else {
		a.LikeNum--
	}
}

func (a *Article) Print() gin.H {
	return gin.H{
		"id":        a.ID,
		"title":     a.Title,
		"url":       a.URL,
		"like_num":  a.LikeNum,
		"createdAt": a.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
