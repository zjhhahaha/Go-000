package account

import (
	"context"
	pb "demo/api/account"
	"demo/internal/biz"
)

var _ pb.AccountServiceServer = &Service{}

type Service struct {
	biz biz.AccountBiz
	pb.UnimplementedAccountServiceServer
}

func New(biz biz.AccountBiz) *Service {
	return &Service{
		biz: biz,
	}
}

func (s *Service) GetAccountByName(ctx context.Context, req *pb.GetAccountByIDRequest) (*pb.GetAccountByIDReply, error) {
	account, err := s.biz.GetAccountByID(ctx, req.Id)
	if err != nil {

	}
	return &pb.GetAccountByIDReply{
		Name: account.Name,
		Sex:  account.Sex,
		Age:  account.Age,
	}, nil
}
