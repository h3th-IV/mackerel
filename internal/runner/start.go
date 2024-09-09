package runner

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	database "github.com/h3th-IV/mackerel/internal/database"
	"github.com/h3th-IV/mackerel/internal/handlers"
	"github.com/h3th-IV/mackerel/internal/server"
	"github.com/h3th-IV/mackerel/internal/utils"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

type StartRunner struct {
	ListenAddr string

	LoggingProduction      bool
	LoggingOutputPath      string
	LoggingLevel           string
	ErrorLoggingOutputPath string

	MySQLDatabaseHost     string
	MySQLDatabasePort     string
	MySQLDatabaseUser     string
	MySQLDatabasePassword string
	MySQLDatabaseName     string
}

func (runner *StartRunner) Run(c *cli.Context) error {
	var (
		loggerConfig        = zap.NewDevelopmentConfig()
		logger              *zap.Logger
		err                 error
		mysqlDbInstance     *sql.DB
		mysqlDatabaseClient database.Database
	)
	if runner.LoggingProduction {
		loggerConfig = zap.NewProductionConfig()
		loggerConfig.OutputPaths = []string{runner.LoggingOutputPath}
		loggerConfig.ErrorOutputPaths = []string{runner.ErrorLoggingOutputPath}
	}

	if err = loggerConfig.Level.UnmarshalText([]byte(runner.LoggingLevel)); err != nil {
		return err
	}

	if logger, err = loggerConfig.Build(); err != nil {
		return err
	}

	logger.Sync()
	databaseConfig := &mysql.Config{
		User:                 runner.MySQLDatabaseUser,
		Passwd:               runner.MySQLDatabasePassword,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", runner.MySQLDatabaseHost, runner.MySQLDatabasePort),
		DBName:               runner.MySQLDatabaseName,
		AllowNativePasswords: true,
	}

	const maxRetries = 3
	const retryDelay = 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		if mysqlDbInstance, err = sql.Open("mysql", databaseConfig.FormatDSN()); err == nil {
			if err = mysqlDbInstance.Ping(); err == nil {
				// Successfully connected
				break
			}
		}
		log.Printf("Failed to connect to MySQL database, attempt %d: %v", i+1, err)
		time.Sleep(retryDelay)
	}

	if err != nil {
		utils.Logger.Error("Error connecting to the database after multiple attempts", zap.Error(err))
		return fmt.Errorf("unable to open connection to MySQL Server: %s", err.Error())
	}

	if mysqlDatabaseClient, err = database.NewMySQLDatabase(mysqlDbInstance); err != nil {
		utils.Logger.Log(zap.ErrorLevel, "unable to create mysql client")
		return fmt.Errorf("unable to create MySQL database client: %s", err.Error())
	}
	mailerConfig := utils.LoadMailerConfig()
	mailer, err := utils.NewMailer(mailerConfig)
	if err != nil {
		utils.Logger.Log(zap.ErrorLevel, "unable to create mailer")
		return fmt.Errorf("unable to create mailer client: %s", err.Error())
	}
	utils.Logger.Info("connected to database successfully")
	server := &server.GracefulShutdownServer{
		HTTPListenAddr: runner.ListenAddr,
		//contollers here
		MCSAttackHandler:   handlers.NewMCSAttackHandler(logger, mailer),
		CaptureDataHandler: handlers.NewCaptureHandler(logger, mysqlDatabaseClient),
		Logger:             logger,
		Mysqlclient:        mysqlDatabaseClient,
	}
	server.Start()
	return nil
}
