package repository

import (
	"context"

	"gorm.io/gorm"
)

type txCtxKey struct{}

// WithTx attaches a GORM transaction to the context so GetDBByCtx can retrieve it.
func WithTx(ctx context.Context, tx *gorm.DB) context.Context {
	if tx == nil {
		return ctx
	}
	return context.WithValue(ctx, txCtxKey{}, tx)
}

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
	if tx, ok := ctx.Value(txCtxKey{}).(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
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
