// +build wireinject
package pkg

import (
	"Go-000/Week04/internal/biz"
	"Go-000/Week04/internal/data"

	"github.com/google/wire"
)

func InitializeUser1() (*biz.UserBiz, error) {
	wire.Build(data.NewDB, data.NewData,
		wire.Bind(new(biz.UserRepo), new(*data.Data)), biz.NewBiz)
	return &biz.UserBiz{}, nil
}
