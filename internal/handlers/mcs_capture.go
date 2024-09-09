package handlers

import (
	"net/http"

	"github.com/h3th-IV/mackerel/internal/database"
	"go.uber.org/zap"
)

/*
THis Handler/controller is used to simulate capturing victim data
*/

var (
	_ http.Handler = &CaptureDataHandler{}
)

type CaptureDataHandler struct {
	logger      *zap.Logger
	mysqlclient database.Database
}

func NewCaptureHandler(logger *zap.Logger, mysqlclient database.Database) *CaptureDataHandler {
	return &CaptureDataHandler{
		logger:      logger,
		mysqlclient: mysqlclient,
	}
}

func (handler *CaptureDataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
