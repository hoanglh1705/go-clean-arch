package sqlxhelper

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type sqlServerDB struct {
	db *sqlx.DB
}

func NewsqlServerDB(host string, port int, username, password, database string) SqlDatabase {
	db, err := initSqlServerDB(host, port, username, password, database)
	if err != nil {
		zap.S().Panic("Failed to init sqlserver", zap.Error(err))
	}
	return &sqlServerDB{
		db: db,
	}
}

func (h *sqlServerDB) QueryRowsPaging(statement string, offset, limit uint32,
	agruments []interface{}) (rows *sqlx.Rows, errQuery error) {
	db := h.Open()
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(statement)
	queryBuilder.WriteString(fmt.Sprintf(" offset :%d rows fetch next :%d rows only", len(agruments)+1, len(agruments)+2))
	agruments = append(agruments, offset, limit)
	rows, errQuery = db.Queryx(queryBuilder.String(), agruments...)
	if errQuery != nil {
		return nil, errQuery
	}
	return rows, nil
}

func (h *sqlServerDB) QueryRowPaging(statement string, offset, limit uint32,
	agruments []interface{}) (row *sqlx.Row) {
	db := h.Open()
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(statement)
	queryBuilder.WriteString(fmt.Sprintf(" offset :%d rows fetch next :%d rows only", len(agruments)+1, len(agruments)+2))
	agruments = append(agruments, offset, limit)
	row = db.QueryRowx(queryBuilder.String(), agruments...)
	return row
}

func (h sqlServerDB) Open() *sqlx.DB {
	return h.db
}

func (h *sqlServerDB) Close() error {
	return h.db.Close()
}

func (h sqlServerDB) Begin() (*sqlx.Tx, error) {
	return h.db.Beginx()
}

func (h sqlServerDB) Commit(tx *sqlx.Tx) error {
	return tx.Commit()
}

func (h sqlServerDB) Rollback(tx *sqlx.Tx) error {
	return tx.Rollback()
}

func (h sqlServerDB) GetConn() (*sqlx.Conn, error) {
	return h.db.Connx(context.Background())
}

func initSqlServerDB(host string, port int, username, password, database string) (*sqlx.DB, error) {
	databaseURL := fmt.Sprintf(
		"server=%s;port=%d;user id=%s;password=%s;database=%s;",
		host,
		port,
		username,
		password,
		database,
	)
	db, err := sqlx.Connect("sqlserver", databaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	zap.S().Infof("Successfully connected!")
	return db, nil
}
