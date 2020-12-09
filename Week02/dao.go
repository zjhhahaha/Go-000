package demo

import (
	"database/sql"

	"github.com/pkg/errors"
)

type Dao interface {
	GetUser() (*User, error)
}

type dao struct {
}

type User struct{}

func mockQuery(sql string) error {
	return nil
}

func (d *dao) GetUser() (*User, error) {
	sql := "xxxxx"
	err := mockQuery(sql)
	if err != nil {
		return nil, errors.Wrapf(err, "Dao: get user error;sql:%s", sql)
	}
	return &User{}, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
