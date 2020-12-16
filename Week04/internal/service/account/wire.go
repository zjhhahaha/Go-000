package account

import (
	"demo/internal/biz"
	"demo/internal/data"
	"demo/internal/pkg"

	"github.com/google/wire"
)

func initializeService(dsn string) *Service {
	wire.Build(New, data.Provider, pkg.NewRDBConnection, biz.Provider)
	return &Service{}
}
