package pg

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Repo struct {
	C *sqlx.DB
}

func New(conn *sql.DB) *Repo {
	return &Repo{C: sqlx.NewDb(conn, "postgres")}
}

func (r *Repo) FindTableByName(ctx context.Context, tableName string) error {
	var (
		unused string
	)
	return r.C.GetContext(ctx, &unused, queryExistsTableWithName, tableName)
}

func (r *Repo) Exec(ctx context.Context, q string) error {
	_, err := r.C.ExecContext(ctx, q)
	return err
}
