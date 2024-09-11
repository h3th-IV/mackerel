package handlers

import (
	"net/http"

	"github.com/h3th-IV/mackerel/internal/utils"
	"go.uber.org/zap"
)

/*
This handler/controller is used to simulate microsoft security alert email --intended to PHISH
*/
var (
	_ http.Handler = &MCSAttackHandler{}
)

type MCSAttackHandler struct {
	logger *zap.Logger
	mailer *utils.Mailer
}

func NewMCSAttackHandler(logger *zap.Logger, mailer *utils.Mailer) *MCSAttackHandler {
	return &MCSAttackHandler{
		logger: logger,
		mailer: mailer,
	}
}

func (handler *MCSAttackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//send the phishing email here
}
