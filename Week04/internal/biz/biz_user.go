package biz

import (
	"Go-000/Week04/internal/data"
)

type UserRepo interface {
	GetUserById(id int) (*data.UserInfor, error)
}
type UserBiz struct {
	userRepo UserRepo
}

func NewBiz(userRepo UserRepo) *UserBiz {
	return &UserBiz{userRepo: userRepo}
}

func (u *UserBiz) GetUserById(id int) (*data.UserInfor, error) {
	return u.userRepo.GetUserById(id)
}
