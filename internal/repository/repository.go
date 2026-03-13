package repository

import (
	"context"
	"fmt"
	"time"

	"nunu/internal/model"
	"nunu/pkg/log"
	"nunu/pkg/zapgorm2"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const ctxTxKey = "TxKey"

type Repository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewRepository(
	logger *log.Logger,
	db *gorm.DB,
) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

type Transaction interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewTransaction(r *Repository) Transaction {
	return r
}

func (r *Repository) DB(ctx context.Context) *gorm.DB {
	v := ctx.Value(ctxTxKey)
	if v != nil {
		if tx, ok := v.(*gorm.DB); ok {
			return tx
		}
	}
	return r.db.WithContext(ctx)
}

func (r *Repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, ctxTxKey, tx)
		return fn(ctx)
	})
}

// NewDB creates the main application database connection.
func NewDB(conf *viper.Viper, l *log.Logger) *gorm.DB {
	logger := zapgorm2.New(l.Logger)
	dsn := conf.GetString("data.db.app.dsn")

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect app database: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Enable pgvector extension
	db.Exec("CREATE EXTENSION IF NOT EXISTS vector")

	// Auto-migrate models
	if err := db.AutoMigrate(&model.Memory{}, &model.Knowledge{}, &model.Conversation{}); err != nil {
		panic(fmt.Sprintf("auto migrate error: %v", err))
	}

	return db
}

// QueryDB is a named type for the read-only query database.
type QueryDB struct {
	*gorm.DB
}

// NewQueryDB creates the read-only query database connection.
func NewQueryDB(conf *viper.Viper, l *log.Logger) *QueryDB {
	logger := zapgorm2.New(l.Logger)
	dsn := conf.GetString("data.db.query.dsn")
	if dsn == "" {
		// Fall back to app DB if query DB not configured
		dsn = conf.GetString("data.db.app.dsn")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect query database: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &QueryDB{DB: db}
}
