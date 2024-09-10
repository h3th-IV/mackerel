package database

import (
	"context"
	"database/sql"
	"io"
	"log"

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
		captureDataStmt = `insert into data (email, username, password, location, ip_address) values (?, ?, ?, ?, ?);`
		database        = &mysqlDatabase{DB: db}
		err             error
	)
	if database.captureData, err = db.Prepare(captureDataStmt); err != nil {
		return nil, err
	}
	return database, nil
}

func (db *mysqlDatabase) CaptureData(ctx context.Context, user *models.User) (bool, error) {
	result, err := db.captureData.ExecContext(ctx, user.Email, user.UserName, user.Password, user.Location, user.IpAddress)
	if err != nil {
		log.Println(err)
		return false, err
	}
	res_lid, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return false, err
	}
	res_ra, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return false, err
	}
	if res_lid <= 0 && res_ra <= 0 {
		log.Println("err inserting new data")
		return false, err
	}
	return true, nil
}
