package command

import (
	"context"
	"fmt"
	"time"

	"github.com/h3th-IV/mackerel/internal/models"
	"github.com/h3th-IV/mackerel/internal/runner"
	"github.com/h3th-IV/mackerel/internal/utils"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func StartCommand() *cli.Command {
	var startRunner = &runner.StartRunner{}

	cmd := &cli.Command{
		Name:  "start",
		Usage: "start the application and server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "listen-addr",
				EnvVars:     []string{"LISTEN_ADDR"},
				Usage:       "The address that the server will listen for request on",
				Destination: &startRunner.ListenAddr,
				Value:       ":9008",
			},
			&cli.StringFlag{
				Name:        "mysql-database-name",
				EnvVars:     []string{"MC_DBNAME"},
				Usage:       "Sample database name",
				Destination: &startRunner.MySQLDatabaseName,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mysql-database-password",
				EnvVars:     []string{"MC_PASSWORD"},
				Usage:       "Sample database password",
				Destination: &startRunner.MySQLDatabasePassword,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mysql-database-User",
				EnvVars:     []string{"MC_USER"},
				Usage:       "Sample database user",
				Destination: &startRunner.MySQLDatabaseUser,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mysql-database-Host",
				EnvVars:     []string{"MC_HOST"},
				Usage:       "Sample database host",
				Destination: &startRunner.MySQLDatabaseHost,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mysql-database-Port",
				EnvVars:     []string{"MC_PORT"},
				Usage:       "Sample database port",
				Destination: &startRunner.MySQLDatabasePort,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mailer-region",
				EnvVars:     []string{"MAILER_REGION"},
				Usage:       "Sample mailer host",
				Destination: &utils.MailerRegion,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mailer-access_id_key",
				EnvVars:     []string{"MAILER_ACCESS_ID_KEY"},
				Usage:       "Sample mailer port",
				Destination: &utils.MailerAccessIDKey,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mailer-secret_access_key",
				EnvVars:     []string{"MAILER_SECRET_ACCESS_KEY"},
				Usage:       "Sample mailer username",
				Destination: &utils.MailerSecretAccessKey,
				Value:       "",
			},
			&cli.StringFlag{
				Name:        "mailer-sender",
				EnvVars:     []string{"MAILER_SENDER"},
				Usage:       "Sample mailer username",
				Destination: &utils.MailerSender,
				Value:       "",
			},
			&cli.StringFlag{
				Name:     "mcs-alert",
				Aliases:  []string{"m"},
				Usage:    "Send microsoft security alert",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			fmt.Println(`
            ><(((º>  ><(((º> ><(((º>   ><(((º>  ><(((º> ><(((º>
       ><(((º>  ><(((º> ><(((º>           ><(((º>  ><(((º> ><(((º>
  ><(((º>  ><(((º> ><(((º> ><(((º>  ><(((º>  ><(((º> ><(((º> ><(((º>
       ><(((º>  ><(((º> ><(((º>   ><(((º>  ><(((º> ><(((º>
            `)
			fmt.Println("Mackerel is starting...")
			time.Sleep(5 * time.Second)
			if ctx.IsSet("mcs-alert") {
				email := ctx.String("mcs-alert")
				if err := MSCAttack(email); err != nil {
					return fmt.Errorf("failed to send Microsoft security alert: %v", err)
				}
			}
			return startRunner.Run(ctx)
		},
	}
	return cmd
}

// delcare new mailer and logger
// send mailer config here
func MSCAttack(email string) error {
	mailerConfig := utils.LoadMailerConfig()
	mailer, err := utils.NewMailer(mailerConfig)
	if err != nil {
		utils.Logger.Log(zap.ErrorLevel, "unable to create mailer")
		return fmt.Errorf("unable to create mailer client: %s", err.Error())
	}
	payload := models.AttackPayload{
		Email:         email,
		MaliciousLink: "link here",
	}
	err = mailer.MSCAttack(context.TODO(), email, payload)
	if err != nil {
		utils.Logger.Error("err sending phishing email:", zap.Error(err))
		return err
	}
	utils.Logger.Info("Phishing Email sent Successfully to", zap.Any("email", email))
	return nil
}
