package sqlxhelper

import "github.com/jmoiron/sqlx"

type (
	SqlDatabase interface {
		Open() *sqlx.DB
		Close() error
		Begin() (*sqlx.Tx, error)
		Commit(tx *sqlx.Tx) error
		Rollback(tx *sqlx.Tx) error
		QueryRowsPaging(stmt string, offset, limit uint32, args []any) (*sqlx.Rows, error)
		QueryRowPaging(stmt string, offset, limit uint32, args []any) *sqlx.Row
		GetConn() (*sqlx.Conn, error)
	}
)
