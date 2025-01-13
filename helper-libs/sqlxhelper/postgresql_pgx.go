package sqlxhelper

import (
	"context"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type postgresqlPgxDB struct {
	db *sqlx.DB
}

func NewPostgresqlPgxDB(host string, port int, username, password, database, schema string) SqlDatabase {
	db, err := initPostgresqlPgxDB(host, port, username, password, database, schema)
	if err != nil {
		zap.S().Panic("Failed to init postgresql", zap.Error(err))
	}
	return &postgresqlPgxDB{
		db: db,
	}
}

func (h *postgresqlPgxDB) QueryRowsPaging(statement string, offset, limit uint32,
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

func (h *postgresqlPgxDB) QueryRowPaging(statement string, offset, limit uint32,
	agruments []interface{}) (row *sqlx.Row) {
	db := h.Open()
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(statement)
	queryBuilder.WriteString(fmt.Sprintf(" offset $%d rows fetch next $%d rows only", len(agruments)+1, len(agruments)+2))
	agruments = append(agruments, offset, limit)
	row = db.QueryRowx(queryBuilder.String(), agruments...)
	return row
}

func (h postgresqlPgxDB) Open() *sqlx.DB {
	return h.db
}

func (h *postgresqlPgxDB) Close() error {
	return h.db.Close()
}

func (h postgresqlPgxDB) Begin() (*sqlx.Tx, error) {
	return h.db.Beginx()
}

func (h postgresqlPgxDB) Commit(tx *sqlx.Tx) error {
	return tx.Commit()
}

func (h postgresqlPgxDB) Rollback(tx *sqlx.Tx) error {
	return tx.Rollback()
}

func (h postgresqlPgxDB) GetConn() (*sqlx.Conn, error) {
	return h.db.Connx(context.Background())
}

func initPostgresqlPgxDB(host string, port int, username, password, database, schema string) (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable search_path=%s",
		host,
		port,
		username,
		password,
		database,
		schema,
	)
	db, err := sqlx.Connect("pgx", psqlInfo)
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
