package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/pablor21/goms/pkg/database/config"
	"github.com/pablor21/goms/pkg/database/migrator"

	// "gorm.io/driver/mysql"
	// "gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormConnection struct {
	_conn     *gorm.DB
	_sqlDb    *sql.DB
	_config   config.DatabaseConnectionConfig
	_migrator migrator.Migrator
}

func NewGormConnection(config config.DatabaseConnectionConfig) DbConnection {
	conn, err := sql.Open(config.Driver, config.URI)
	if err != nil {
		panic(err)
	}

	// conn.SetMaxIdleConns(64)
	// conn.SetMaxOpenConns(64)
	// conn.SetConnMaxLifetime(time.Minute)

	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 64
	}

	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 64
	}

	if config.ConnMaxLifetime == 0 {
		config.ConnMaxLifetime = 60
	}

	conn.SetMaxIdleConns(config.MaxIdleConns)
	conn.SetMaxOpenConns(config.MaxOpenConns)
	conn.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	var gormDb *gorm.DB
	switch config.Driver {
	case "sqlite", "sqlite3":
		gormDb, err = gorm.Open(sqlite.New(sqlite.Config{
			Conn: conn,
		}), gormConfig)
	// case "mysql":
	// 	gormDb, err = gorm.Open(mysql.New(mysql.Config{
	// 		Conn: conn,
	// 	}), gormConfig)
	// case "postgres":
	// 	gormDb, err = gorm.Open(postgres.New(postgres.Config{
	// 		Conn: conn,
	// 	}), gormConfig)
	default:
		panic("Unsupported database driver" + config.Driver)
	}

	if err != nil {
		panic(err)
	}

	return &GormConnection{
		_conn:     gormDb,
		_sqlDb:    conn,
		_config:   config,
		_migrator: NewDefaultMigrator(config, conn),
	}
}

func (g *GormConnection) Close() error {
	sqlDb, err := g._conn.DB()
	if err != nil {
		return err
	}
	return sqlDb.Close()
}

func (g *GormConnection) Conn() *gorm.DB {
	return g._conn
}

func (g *GormConnection) Migrator() migrator.Migrator {
	return g._migrator
}

func (g *GormConnection) Paginate(query *gorm.DB, take, skip int) *gorm.DB {
	return query.Limit(take).Offset(skip)
}

// Get dbConn from context
// If there is a transaction in the context, returns the transaction
func (g *GormConnection) GetContextDb(ctx context.Context) (*gorm.DB, context.Context, error) {
	var txKey txKeyType = txKeyType("db_tx_" + g._config.Name)
	db, ok := ctx.Value(txKey).(*gormTx)
	if ok {
		return db.DB, ctx, nil
	} else {
		return g._conn, ctx, nil
	}
}

// Get a transaction from the context or create a new one
// Do not forget to commit or rollback the transaction and use the same context
func (g *GormConnection) GetContextTx(ctx context.Context) (*gormTx, context.Context, error) {
	var txKey txKeyType = txKeyType("db_tx_" + g._config.Name)
	tx, ok := ctx.Value(txKey).(*gormTx)
	if !ok {
		tx = &gormTx{
			DB:     g._conn.Begin(),
			isMain: true,
		}

	} else {
		tx = &gormTx{
			DB:     tx.DB,
			isMain: false,
		}
	}
	// tx = &gormTx{
	// 	DB:     g._conn.Begin(),
	// 	isMain: true,
	// }

	ctx = context.WithValue(ctx, txKey, tx)
	return tx, ctx, nil
}

type gormTx struct {
	*gorm.DB
	isMain bool
}

func (tx *gormTx) Commit() error {
	if tx.isMain {
		return tx.DB.Commit().Error
	}
	return nil
}

func (tx *gormTx) Rollback() error {
	if tx.isMain {
		return tx.DB.Rollback().Error
	}
	return nil
}
