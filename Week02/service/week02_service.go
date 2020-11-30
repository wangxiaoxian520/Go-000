package service

import (
	"Go-000/Week02/dao"
	"errors"
	"strconv"

	xerrors "github.com/pkg/errors"
)

var (
	ErrorParam = errors.New("Param error")
	ErrorConv  = errors.New("conversion faile")
)

func GetUserInfor(id string) (string, error) {
	userId, err := strconv.Atoi(id)
	if err != nil {
		return "", xerrors.Wrap(ErrorConv, "userId faild")
	}
	if userId <= 0 {
		return "", xerrors.Wrap(ErrorParam, "userId invalid")
	}
	return dao.GetUserInfor(userId)
}
