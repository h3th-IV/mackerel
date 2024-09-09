package database

import (
	"context"
	"database/sql"
	"io"

	"github.com/h3th-IV/mackerel/internal/models"
)

var (
	_ Database  = &mysqlDatabase{}
	_ io.Closer = &mysqlDatabase{}
)

type mysqlDatabase struct {
	*sql.DB
	captureData *sql.Stmt
}

func NewMySQLDatabase(db *sql.DB) (*mysqlDatabase, error) {
	var (
		captureDataStmt = ``
		database      = &mysqlDatabase{DB: db}
		err           error
	)
	if database.captureData, err = db.Prepare(captureDataStmt); err != nil {
		return nil, err
	}
	return database, nil
}

func (db *mysqlDatabase) CaptureData(ctx context.Context, user models.User) (bool, error) {
	return true, nil
}
