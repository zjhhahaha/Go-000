package data

import (
	"context"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Repository interface {
	GetAccountById(ctx context.Context, id int64) (*Account, error)
}

type Account struct {
	ID        int64
	Name      string
	Sex       bool
	Age       int64
	Address   string
	CreatedAt time.Time
	UpdateAt  time.Time
}

type DB struct {
	rdb   *gorm.DB
	cache *redis.Client
}

func New(db *gorm.DB) *DB {
	return &DB{rdb: db}
}

func (db *DB) GetAccountById(ctx context.Context, id int64) (*Account, error) {
	account := Account{}
	data, err := db.cache.Get(fmt.Sprintf("xxx:%d", id)).Result()
	if err != nil {
		if err != redis.Nil {
			// handle
		}
	} else {
		err := json.UnmarshalFromString(data, &account)
		if err != nil {
			// handle
		}
	}
	err = db.rdb.WithContext(ctx).Where("id = ?", id).First(&account).Error
	if err != nil {
		return nil, errors.Wrap(err, "Dao: get account by id error")
	}
	return &account, nil
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
