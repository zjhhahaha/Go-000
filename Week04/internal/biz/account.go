package biz

import (
	"context"
	"demo/internal/data"
)

type AccountBiz interface {
	GetAccountByID(ctx context.Context, id int64) (*data.Account, error)
}

type accountBiz struct {
	repo data.Repository
}

func New(repo data.Repository) *accountBiz {
	return &accountBiz{
		repo: repo,
	}
}

func (biz *accountBiz) GetAccountByID(ctx context.Context, id int64) (*data.Account, error) {
	account, err := biz.repo.GetAccountById(ctx, id)
	if err != nil {
		return nil, err
	}
	return account, nil
}
