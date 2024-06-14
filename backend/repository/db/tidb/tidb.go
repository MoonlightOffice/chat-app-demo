package tidb

import (
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"chat-app-demo/config"
	"chat-app-demo/repository/db/types"

	"github.com/go-sql-driver/mysql"
)

var (
	pool       *sql.DB
	maxConns   = int(max(config.AppConfig.DBMaxConns, 1))
	dsn        = config.AppConfig.MySQLDSN
	serverName = config.AppConfig.MySQLServerName
)

func setupPool() error {
	if config.AppConfig.Mode != config.ModeLocal {
		mysql.RegisterTLSConfig("tidb", &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: serverName,
		})
	}

	var err error
	pool, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	pool.SetConnMaxLifetime(3600)
	pool.SetMaxIdleConns(maxConns)
	pool.SetMaxOpenConns(maxConns)

	return nil
}

func NewTiDBRepository() (types.DBRepository, error) {
	// Check if pool is initialized
	if pool == nil {
		err := setupPool()
		if err != nil {
			return nil, fmt.Errorf("crdb.TiDBRepository.NewTiDBRepository(): %w", err)
		}
	}

	return TiDBRepository{crud: connExtended{pool}}, nil
}

type TiDBRepository struct {
	crud CRUD
}

func (repo TiDBRepository) BeginTx() (types.DBTxRepository, error) {
	tx, err := repo.crud.Begin()
	if err != nil {
		return nil, err
	}

	return TiDBRepository{crud: tx}, nil
}

func (repo TiDBRepository) Close() {}

func (repo TiDBRepository) CommitTx() error {
	return repo.crud.Commit()
}

func (repo TiDBRepository) RollbackTx() error {
	return repo.crud.Rollback()
}

type CRUD interface {
	Begin() (*txExtended, error)
	Rollback() error
	Commit() error
	Exec(string, ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(string, ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Implementation of DBRepository

type connExtended struct {
	*sql.DB
}

func (conn connExtended) Begin() (*txExtended, error) {
	_tx, err := conn.DB.Begin()
	if err != nil {
		return nil, err
	}

	return &txExtended{Tx: _tx}, nil
}

func (conn connExtended) Rollback() error {
	return nil
}

func (conn connExtended) Commit() error {
	return nil
}

// Implementation of DBTxRepository

type txExtended struct {
	*sql.Tx

	next      *txExtended
	savePoint int
	resolved  bool
}

func (tx *txExtended) Begin() (*txExtended, error) {
	if tx.next != nil && !tx.next.resolved {
		return nil, errors.New("child transaction already exists")
	}

	tx.next = &txExtended{
		Tx:        tx.Tx,
		savePoint: tx.savePoint + 1,
	}

	_, err := tx.Exec("SAVEPOINT SP" + strconv.Itoa(tx.next.savePoint))
	if err != nil {
		return nil, err
	}

	return tx.next, nil
}

func (tx *txExtended) Rollback() error {
	if tx.resolved {
		return errors.New("transaction has already been committed or rolled back")
	}

	// Rollback uncommited nested transaction
	if tx.next != nil && !tx.next.resolved {
		err := tx.next.Rollback()
		if err != nil {
			return err
		}

		tx.next.resolved = true
	}

	if tx.savePoint > 0 {
		// Rollback to the savepoint
		_, err := tx.Exec("ROLLBACK TO SAVEPOINT SP" + strconv.Itoa(tx.savePoint))
		if err != nil {
			return err
		}

		// Remove the savepoint to abort the nested transaction
		_, err = tx.Exec("RELEASE SAVEPOINT SP" + strconv.Itoa(tx.savePoint))
		if err != nil {
			return err
		}

		tx.resolved = true

		return err
	}

	err := tx.Tx.Rollback()
	if err != nil {
		return err
	}

	tx.resolved = true

	return nil
}

func (tx *txExtended) Commit() error {
	if tx.resolved {
		return errors.New("transaction has already been committed or rolled back")
	}

	// Rollback uncommited nested transaction
	if tx.next != nil && !tx.next.resolved {
		err := tx.next.Rollback()
		if err != nil {
			return err
		}

		tx.next.resolved = true
	}

	if tx.savePoint > 0 {
		_, err := tx.Exec("RELEASE SAVEPOINT SP" + strconv.Itoa(tx.savePoint))
		if err != nil {
			return err
		}

		tx.resolved = true

		return nil
	} else {
		err := tx.Tx.Commit()
		if err != nil {
			return err
		}

		tx.resolved = true

		return nil
	}
}
