package sqlxhelper

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type postgresqlDB struct {
	db *sqlx.DB
}

func NewPostgresqlDB(host string, port int, username, password, database, schema string) SqlDatabase {
	db, err := initPostgresqlDB(host, port, username, password, database, schema)
	if err != nil {
		zap.S().Panic("Failed to init postgresql", zap.Error(err))
	}
	return &postgresqlDB{
		db: db,
	}
}

func (h *postgresqlDB) QueryRowsPaging(statement string, offset, limit uint32,
	agruments []interface{}) (rows *sqlx.Rows, errQuery error) {
	db := h.Open()
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(statement)
	queryBuilder.WriteString(fmt.Sprintf(" offset $%d rows fetch next $%d rows only", len(agruments)+1, len(agruments)+2))
	agruments = append(agruments, offset, limit)
	rows, errQuery = db.Queryx(queryBuilder.String(), agruments...)
	if errQuery != nil {
		return nil, errQuery
	}
	return rows, nil
}

func (h *postgresqlDB) QueryRowPaging(statement string, offset, limit uint32,
	agruments []interface{}) (row *sqlx.Row) {
	db := h.Open()
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(statement)
	queryBuilder.WriteString(fmt.Sprintf(" offset $%d rows fetch next $%d rows only", len(agruments)+1, len(agruments)+2))
	agruments = append(agruments, offset, limit)
	row = db.QueryRowx(queryBuilder.String(), agruments...)
	return row
}

func (h postgresqlDB) Open() *sqlx.DB {
	return h.db
}

func (h *postgresqlDB) Close() error {
	return h.db.Close()
}

func (h postgresqlDB) Begin() (*sqlx.Tx, error) {
	return h.db.Beginx()
}

func (h postgresqlDB) Commit(tx *sqlx.Tx) error {
	return tx.Commit()
}

func (h postgresqlDB) Rollback(tx *sqlx.Tx) error {
	return tx.Rollback()
}

func (h postgresqlDB) GetConn() (*sqlx.Conn, error) {
	return h.db.Connx(context.Background())
}

func initPostgresqlDB(host string, port int, username, password, database, schema string) (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable search_path=%s",
		host,
		port,
		username,
		password,
		database,
		schema,
	)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	zap.S().Infof("Successfully connected!")
	return db, nil
}
