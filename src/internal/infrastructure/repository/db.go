package repository

import (
	"context"

	"gorm.io/gorm"
)

type DBInterface interface {
	GetDB() *gorm.DB
	GetDBByCtx(ctx context.Context) *gorm.DB
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	Close() error
}

type SchemaDB struct {
	db *gorm.DB
}

func NewSchemaDB(db *gorm.DB) DBInterface {
	return &SchemaDB{db: db}
}

func (d *SchemaDB) GetDB() *gorm.DB {
	return d.db
}

func (d *SchemaDB) GetDBByCtx(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx)
}

func (d *SchemaDB) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return d.GetDBByCtx(ctx).Transaction(fn)
}

func (d *SchemaDB) Close() error {

	return closeGormDB(d.db)
}

func closeGormDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
