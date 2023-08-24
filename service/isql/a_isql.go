package isql

import (
	"errors"

	"github.com/haobinfei/ginner/model"
	"github.com/haobinfei/ginner/public/common"
	"github.com/haobinfei/ginner/public/tools"
)

type UserService struct{}

// Login 登录
func (s UserService) Login(user *model.User) (*model.User, error) {
	var firstUser model.User

	err := s.Find(tools.H{"username": user.UserName}, &firstUser)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	// 判断用户的状态
	userStatus := firstUser.Status
	if userStatus != 1 {
		return nil, errors.New("用户被禁用")
	}

	if tools.NewParPassword(firstUser.Password) != user.Password {
		return nil, errors.New("密码错误")
	}

	return &firstUser, nil
}

func (s UserService) Find(filter interface{}, data *model.User) error {
	return common.DB.Where(filter).Preload("Roles").First(data).Error
}
