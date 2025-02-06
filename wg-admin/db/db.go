package db

import (
	"context"
	"database/sql"
	"mp2720/wg-admin/wg-admin/db/sqlgen"
	"mp2720/wg-admin/wg-admin/transaction"
	"mp2720/wg-admin/wg-admin/utils"

	_ "embed"

	"github.com/mattn/go-sqlite3"
)

// Wrappers around sql.Tx and sql.DB that implement [transaction.Tx] and [transaction.Initiator]
// interfaces.

type DB struct {
	db *sql.DB
}

type Tx struct {
	sql *sql.Tx
}

func (tx Tx) Commit(ctx context.Context) error {
	return tx.sql.Commit()
}

func (tx Tx) Rollback(ctx context.Context) error {
	return tx.sql.Rollback()
}

func (db DB) Clock() error {
    return db.db.Close()
}

func (db DB) Begin(ctx context.Context) (transaction.Tx, error) {
	sqlTx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}
	return Tx{sqlTx}, nil
}

func (db DB) With(ctx context.Context) *sqlgen.Queries {
	tx := transaction.FromCtx(ctx)
	if tx == nil {
		return sqlgen.New(db.db)
	}
	return sqlgen.New(tx.(Tx).sql)
}

//go:embed sql/schema.sql
var dbSchemaSql string

func (db DB) CreateTablesIfNotExists(ctx context.Context) error {
	_, err := db.db.ExecContext(ctx, dbSchemaSql)
	return err
}

// Replaces errors from sql package and driver-specific errors with errors from utils.
//
// sql.ErrNoRows                                -> ErrNotFound
// (driver-specific) Unique constraint failure  -> ErrAlreadyExists
func HandleSQLError(what string, err error) error {
	if err == sql.ErrNoRows {
		return utils.ErrNotFound{What: what}
	}

	if sqliteErr, ok := err.(sqlite3.Error); ok {
		switch sqliteErr.ExtendedCode {
		case sqlite3.ErrConstraintUnique:
			return utils.ErrAlreadyExists{What: what}
		}
	}

	return err
}

func NewDB(ctx context.Context, filepath string) (DB, error) {
	sqldb, err := sql.Open("sqlite3", filepath+"?_txlock=immediate")
	if err != nil {
		return DB{}, err
	}

	db := DB{sqldb}
	if err = db.CreateTablesIfNotExists(ctx); err != nil {
		return DB{}, err
	}

	return db, nil
}
