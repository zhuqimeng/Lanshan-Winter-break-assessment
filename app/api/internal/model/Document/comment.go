package Document

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Username string `json:"username" gorm:"size:255;not null"`
	Link     string `json:"link" gorm:"not null"`
	Content  string `json:"content" gorm:"type:longtext;not null"`
}

func (c Comment) Print() gin.H {
	return gin.H{
		"id":        c.ID,
		"link":      c.Link,
		"content":   c.Content,
		"createdAt": c.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
