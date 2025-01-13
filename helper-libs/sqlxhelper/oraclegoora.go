package sqlxhelper

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	go_ora "github.com/sijms/go-ora/v2"
	"go.uber.org/zap"
)

type oracleGoOraDB struct {
	db *sqlx.DB
}

func NewOracleGoOraDB(host string, port int, username, password, database string) SqlDatabase {
	db, err := initOracleGoOraDB(host, port, username, password, database)
	if err != nil {
		zap.S().Panic("Failed to init go_ora", zap.Error(err))
	}
	return &oracleGoOraDB{
		db: db,
	}
}

func (h *oracleGoOraDB) QueryRowsPaging(statement string, offset, limit uint32,
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

func (h *oracleGoOraDB) QueryRowPaging(statement string, offset, limit uint32,
	agruments []interface{}) (row *sqlx.Row) {
	db := h.Open()
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(statement)
	queryBuilder.WriteString(fmt.Sprintf(" offset :%d rows fetch next :%d rows only", len(agruments)+1, len(agruments)+2))
	agruments = append(agruments, offset, limit)
	row = db.QueryRowx(queryBuilder.String(), agruments...)
	return row
}

func (h oracleGoOraDB) Open() *sqlx.DB {
	return h.db
}

func (h *oracleGoOraDB) Close() error {
	return h.db.Close()
}

func (h oracleGoOraDB) Begin() (*sqlx.Tx, error) {
	return h.db.Beginx()
}

func (h oracleGoOraDB) Commit(tx *sqlx.Tx) error {
	return tx.Commit()
}

func (h oracleGoOraDB) Rollback(tx *sqlx.Tx) error {
	return tx.Rollback()
}

func (h oracleGoOraDB) GetConn() (*sqlx.Conn, error) {
	return h.db.Connx(context.Background())
}

func initOracleGoOraDB(host string, port int, username, password, database string) (*sqlx.DB, error) {
	databaseURL := go_ora.BuildUrl(host, port, database, username, password, nil)

	db, err := sqlx.Connect("oracle", databaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	zap.S().Infof("Successfully connected!")
	return db, nil
}
