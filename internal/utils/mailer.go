package utils

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"go.uber.org/zap"
)

// var mailTemplates embed.FS

type MailerConfig struct {
	TimeOut         time.Duration
	SenderAddr      string
	TemplatePath    string
	SMTPRegion      string
	AccessIDKey     string
	SecretAccessKey string
}

func LoadMailerConfig() *MailerConfig {
	return &MailerConfig{
		TimeOut:         time.Hour,
		SenderAddr:      MailerSender,
		TemplatePath:    "./../../mail_templates",
		SMTPRegion:      MailerRegion,
		AccessIDKey:     MailerAccessIDKey,
		SecretAccessKey: MailerSecretAccessKey,
	}
}

var (
	MailerSender          string
	MailerRegion          string
	MailerAccessIDKey     string
	MailerSecretAccessKey string
)

type Mailer struct {
	config *MailerConfig
	client *ses.Client
}

func NewMailer(config *MailerConfig) (*Mailer, error) {
	//set up the configuration for Mailer
	cfg, err := awscfg.LoadDefaultConfig(context.TODO(), awscfg.WithRegion(config.SMTPRegion), awscfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.AccessIDKey, config.SecretAccessKey, "")))
	if err != nil {
		return nil, err
	}

	//create the ses client
	ses_client := ses.NewFromConfig(cfg)
	if ses_client == nil {
		return nil, fmt.Errorf("unable to create ses client")
	}
	return &Mailer{
		config: config,
		client: ses_client,
	}, nil
}

func (m *Mailer) SendEmail(ctx context.Context, recipient, subject, templateName string, data interface{}) error {
	workingDir, err := os.Getwd()
	if err != nil {
		Logger.Error("err getting working directory", zap.Error(err))
		return err
	}
	// fmt.Println("Current working directory:", workingDir)
	t, err := template.ParseFiles(fmt.Sprintf("%s/mail_templates/%s", workingDir, templateName))
	// t, err := template.ParseFiles(fmt.Sprintf(mail_templates.ReadFile(), templateName))
	if err != nil {
		Logger.Error("err parsing email templates", zap.Error(err))
		return err
	}

	// render templates + data
	var emailBody bytes.Buffer
	err = t.Execute(&emailBody, data)
	if err != nil {
		Logger.Error("err executing templates", zap.Error(err))
		return err
	}

	input := ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{recipient},
		},

		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(emailBody.String()),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(m.config.SenderAddr),
	}
	_, err = m.client.SendEmail(ctx, &input)
	if err != nil {
		return err
	}
	return nil
}

// this will send a microsoft alert security email to the victim
func (m *Mailer) MSCAttack(ctx context.Context, recipient string, data interface{}) error {
	return m.SendEmail(ctx, recipient, "Microsoft account security info", "mcs_attack.html", data)
}
