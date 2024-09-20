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
	captureData    *sql.Stmt
	insertLocation *sql.Stmt
}

func NewMySQLDatabase(db *sql.DB) (*mysqlDatabase, error) {
	var (
		captureDataStmt    = `insert into users (email, user_name, password, ip_address) values (?, ?, ?, ?);`
		insertLocationStmt = `insert into geoLocation (user_id, city, country, ip_address, region, lat_long, organization, timezone) values (?, ?, ?, ?, ?, ?, ?, ?);`
		database           = &mysqlDatabase{DB: db}
		err                error
	)
	if database.captureData, err = db.Prepare(captureDataStmt); err != nil {
		return nil, err
	}
	if database.insertLocation, err = db.Prepare(insertLocationStmt); err != nil {
		return nil, err
	}
	return database, nil
}

func (db *mysqlDatabase) CaptureData(ctx context.Context, user *models.User) (bool, error) {
	result, err := db.captureData.ExecContext(ctx, user.Email, user.UserName, user.Password, user.IpAddress)
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
	locat, err := db.insertLocation.ExecContext(ctx, res_lid, user.Location.City, user.Location.Country, user.Location.IpAddress, user.Location.Region, user.Location.LatLong, user.Location.Organization, user.Location.TimeZone)
	if err != nil {
		log.Println(err)
		return false, err
	}
	locat_lid, err := locat.LastInsertId()
	if err != nil {
		log.Println(err)
		return false, err
	}
	locat_ra, err := locat.RowsAffected()
	if err != nil {
		log.Println(err)
		return false, err
	}
	if locat_lid <= 0 && locat_ra <= 0 {
		log.Println("err inserting location data")
		return false, err
	}
	return true, nil
}
