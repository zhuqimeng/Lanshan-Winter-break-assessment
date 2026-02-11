package User

import "gorm.io/gorm"

type Relation struct {
	gorm.Model
	Username string `json:"username" gorm:"not null"`
	Follower string `json:"follower" gorm:"not null"`
}
