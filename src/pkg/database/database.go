package database

import (
	"github.com/pablor21/goms/pkg/database/config"
	"github.com/pablor21/goms/pkg/database/migrator"
)

var db *DbManager

type txKeyType string

// const txKey txKeyType = "tx"

type DbTx interface {
	Commit() error
	Rollback() error
}

type DbConnection interface {
	// Conn() T
	Migrator() migrator.Migrator
	Close() error
	// GetContextTx(context.Context) (DbTx, context.Context, error)
}

type DbManager struct {
	connections map[string]DbConnection
	config      config.DatabaseConfig
}

func NewDbManager(config config.DatabaseConfig) *DbManager {
	mng := &DbManager{
		connections: make(map[string]DbConnection),
		config:      config,
	}
	db = mng
	return mng
}

func GetConnection(name string) DbConnection {

	if db == nil {
		panic("Database manager not initialized")
	}
	conn, ok := db.connections[name]
	if !ok {
		config, ok := db.config[name]
		if !ok {
			panic("Database connection not found")
		}

		switch config.Type {
		case "gorm":
			conn = NewGormConnection(config)
			db.connections[name] = conn
		default:
			panic("Unsupported database type")
		}
	}

	return conn
}

func Close() {
	if db == nil {
		panic("Database manager not initialized")
	}
	for _, conn := range db.connections {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}
}

// func (m *DbManager) AddConnection(name string, conn DbConnection[any]) {
// 	m.connections[name] = conn
// }

func GetManager() *DbManager {
	if db == nil {
		panic("Database manager not initialized")
	}
	return db
}
