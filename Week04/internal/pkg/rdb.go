package pkg

import "gorm.io/gorm"

func NewRDBConnection(dsn string) *gorm.DB {
	return &gorm.DB{}
}
