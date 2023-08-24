package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"type:varchar(50);not null;unique;comment:'用户名'" json:"user_name"`
	Password string `gorm:"size:255;not null;comment:'密码'" json:"password"`
	Status   uint   `gorm:"type:tinyint(1);default:1;comment:'状态:1在职, 2离职'" json:"status"` // 状态
}
