package storage

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/wen-git-acc/orderbook/config"
	"github.com/wen-git-acc/orderbook/pkg/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// DBInstance holds the singleton database instance
var (
	once sync.Once
)

type StorageImpl struct {
	db     *sqlx.DB
	logger logger.LoggerClientInterface
}

type IStorageImpl interface {
	InitCfg(mode config.Mode, cfg *config.PgSqlConfig, logger logger.LoggerClientInterface)
	ExecuteInTransaction(ctx context.Context, fn TransactionFunc) error
	ExecuteInTransactionWithRead(ctx context.Context, fn TransactionFuncWithResult[interface{}]) (interface{}, error)
	ExecuteRead(ctx context.Context, query string, dest interface{}, args ...interface{}) error
	ExecuteReadMany(ctx context.Context, query string, dest interface{}, args ...interface{}) error
}

var PgStorage IStorageImpl = &StorageImpl{}

// GetDB returns a singleton instance of the database connection
func (c *StorageImpl) InitCfg(mode config.Mode, cfg *config.PgSqlConfig, logger logger.LoggerClientInterface) {
	c.logger = logger.GetLoggerWithProfile("storage")
	var pgDb *sqlx.DB
	var err error
	once.Do(func() {
		c.logger.Info("Initializing database connection")
		// Create the connection string
		connStr := "host=" + cfg.PostgresDbHost + " port=" + strconv.Itoa(cfg.PostgresDbPort) + " dbname=" + cfg.PostgresDbName + " sslmode=disable" + " connect_timeout=" + strconv.Itoa(cfg.PostgresPoolTimeout)

		if mode != config.Development {
			c.logger.Info("Registering user and password for non-development mode")
			connStr = connStr + " user=" + cfg.PostgresDbUser + " password=" + cfg.PostgresDbPassword
		}

		pgDb, err = sqlx.Connect("postgres", connStr)
		if err != nil {
			log.Fatalln(err)
		}

		// Configure the connection pool
		pgDb.SetMaxOpenConns(cfg.PostgresMaxOverflow)                   // Maximum number of open connections
		pgDb.SetMaxIdleConns(cfg.PostgresPoolSize)                      // Maximum number of idle connections
		pgDb.SetConnMaxLifetime(time.Duration(cfg.PostgresPoolRecycle)) // Connection max lifetime (0 = no limit)
	})
	c.db = pgDb

}

// TransactionFunc defines the type of function that can be executed within a transaction
type TransactionFunc func(tx *sql.Tx) error

// ExecuteInTransaction manages the transaction lifecycle
func (storage *StorageImpl) ExecuteInTransaction(ctx context.Context, fn TransactionFunc) error {
	// Begin a new transaction
	tx, err := storage.db.BeginTx(ctx, nil)
	if err != nil {
		storage.logger.Error("error", err, "failed to begin transaction")
		return err
	}

	// Defer rollback in case of error
	defer func() {
		if err != nil {
			storage.logger.Error("error", err, "rolling back transaction")
			tx.Rollback() // Rollback the transaction
			return
		}
		err = tx.Commit() // Commit the transaction
	}()

	// Call the provided function, passing in the transaction
	if err = fn(tx); err != nil {
		storage.logger.Error("error", err, "failed to execute transaction")
		return err // Return the error to trigger rollback
	}

	return nil // Return nil to indicate success
}

// TransactionFuncWithResult defines the type of function that can be executed within a transaction
type TransactionFuncWithResult[T any] func(tx *sql.Tx) (T, error)

// ExecuteInTransactionWithRead manages the transaction lifecycle and allows reading results
func (storage *StorageImpl) ExecuteInTransactionWithRead(ctx context.Context, fn TransactionFuncWithResult[interface{}]) (interface{}, error) {
	tx, err := storage.db.BeginTx(ctx, nil)
	if err != nil {
		storage.logger.Error("error", err, "failed to begin read transaction")
		return nil, err
	}

	defer func() {
		if err != nil {
			storage.logger.Error("error", err, "rolling back read transaction")
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	result, err := fn(tx)
	if err != nil {
		storage.logger.Error("error", err, "failed to execute read transaction")
		return nil, err
	}

	return result, nil
}

// ExecuteRead executes a read query and scans the result into dest.
func (storage *StorageImpl) ExecuteRead(ctx context.Context, query string, dest interface{}, args ...interface{}) error {
	err := storage.db.GetContext(ctx, dest, query, args...)
	if err != nil {
		storage.logger.Error("error", err, "failed to execute read query")
		return err
	}
	return nil
}

// ExecuteReadMany executes a read query that returns multiple rows and scans the result into dest.
func (storage *StorageImpl) ExecuteReadMany(ctx context.Context, query string, dest interface{}, args ...interface{}) error {
	err := storage.db.SelectContext(ctx, dest, query, args...)
	if err != nil {
		storage.logger.Error("error", err, "failed to execute read many query")
		return err
	}
	return nil
}
