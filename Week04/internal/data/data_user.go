package data

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("There is no data.")
	ErrDbServer = errors.New("db error")
)

type UserInfor struct {
	Name string
	Age  int
	Addr string
}

type Data struct {
	db *sql.DB
}

func NewData(db *sql.DB) *Data {
	return &Data{db: db}
}
func NewDB() (*sql.DB, error) {
	conn := "root:123@tcp(127.0.0.1:3306)/mysql?charset=utf8"
	return sql.Open("mysql", conn)
}

func (d *Data) GetUserById(id int) (*UserInfor, error) {
	var userInfor UserInfor
	query := fmt.Sprintf("Select name,age,addr from user where id=%d", id)
	err := d.db.QueryRow(query).
		Scan(&userInfor.Name, &userInfor.Age, &userInfor.Addr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("sql:%s error:[%v]", query, ErrNotFound)
		}
		return nil, errors.Wrapf(ErrDbServer, fmt.Sprintf("query: %s error(%v)", query, err))
	}
	return &userInfor, nil
}
