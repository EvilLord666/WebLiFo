package model

import (
	"WebLiFo/logging"
	"errors"
	g "github.com/wissance/gwuu/gorm"
	"github.com/wissance/stringFormatter"
	"gorm.io/gorm"
)

type AppModelContext struct {
	dbContext *gorm.DB
	logger    *logging.AppLogger
}

func CreateAppModelContext(logger *logging.AppLogger) *AppModelContext {
	return &AppModelContext{logger: logger}
}

func (context *AppModelContext) Init(host string, port int, dbName string, user string, password string,
	maxConns int, maxIdleConns int, useSsl string, createIfNotExists bool, config *gorm.Config) *gorm.DB {
	connStr := g.BuildConnectionString(g.Postgres, host, port, dbName, user, password, useSsl)
	db := g.OpenDb2(g.Postgres, connStr, createIfNotExists, true, config)
	context.dbContext = db

	// db.DB() returns *sql.DB and error. err != nil signifies one of two events:
	// either connection attempt was unsuccessful, or the underlying db connection is not of the *sql.DB type
	// (it should be under normal circumstances).
	// Either way we return the *gorm.DB, which should be nil if something's wrong,
	// then it'll be caulght inside application.go's initDatabase().
	sqlConfig, err := db.DB()
	if err != nil {
		return db
	}

	// If no max_idle_connections is passed, the value will be kept at default = 10
	if maxIdleConns != 0 {
		sqlConfig.SetMaxIdleConns(maxIdleConns)
	}

	// if no max_open_connections is passed, the value will be kept at default = 100
	if maxConns != 0 {
		sqlConfig.SetMaxOpenConns(maxConns)
	}
	// Please note that the max_open_connections value should probably be kept @ default
	// unless postgres deployment is set up (either through postgres config or pooling software)
	// to accept more than default number of connections, which is 100.
	// Also consider that stateMachineService and stateMachineExecutor share the database (as of 2.7),
	// hence total max_connections between them should be kept under postgres deployment's limits.

	return db
}

func (context *AppModelContext) GetContext() *gorm.DB {
	return context.dbContext
}

func (context *AppModelContext) SetContext(dbContext *gorm.DB) {
	context.dbContext = dbContext
}

func (context AppModelContext) Close() bool {
	return g.CloseDb(context.dbContext)
}

// Prepare function for Auto migrate & set up relations & data init
func (context AppModelContext) Prepare() bool {
	err := context.migrateDatabaseTable(context.dbContext)
	if err != nil {
		return false
	}

	return true
}

func (context AppModelContext) migrateDatabaseTable(db *gorm.DB) error {
	errorTemplate := "An error occurred during \"{0}\" table migration: {1}"

	err := db.AutoMigrate(&Lifo{})
	if err != nil {
		errorMsg := stringFormatter.Format(errorTemplate, "lifo", err.Error())
		context.logger.Error(errorMsg)
		return errors.New(errorMsg)
	}

	err = db.AutoMigrate(&LifoItem{})
	if err != nil {
		errorMsg := stringFormatter.Format(errorTemplate, "lifo_item", err.Error())
		context.logger.Error(errorMsg)
		return errors.New(errorMsg)
	}

	return nil
}
