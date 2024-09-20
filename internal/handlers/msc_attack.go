package handlers

import (
	"context"
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
	//extract victim email from the request (for testing or real CLI integration)
	victimEmail := r.URL.Query().Get("email")
	if victimEmail == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	//prepare the phishing email data
	data := struct {
		Email         string
		MaliciousLink string
	}{
		Email:         victimEmail,
		MaliciousLink: "http://fake-malicious-site.com/login", // replace with your fake phishing link
	}

	//send the email
	err := handler.mailer.MSCAttack(context.TODO(), victimEmail, data)
	if err != nil {
		handler.logger.Error("Failed to send phishing email", zap.Error(err))
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Phishing email sent successfully"))
}

func (handler *MCSAttackHandler) MCSAttack(w http.ResponseWriter, r *http.Request) {
	//extract victim email from the request (for testing or real CLI integration)
	victimEmail := r.URL.Query().Get("email")
	if victimEmail == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	//prepare the phishing email data
	data := struct {
		Email         string
		MaliciousLink string
	}{
		Email:         victimEmail,
		MaliciousLink: "http://fake-malicious-site.com/login", // replace with your fake phishing link
	}

	//send the email
	err := handler.mailer.MSCAttack(context.TODO(), victimEmail, data)
	if err != nil {
		handler.logger.Error("Failed to send phishing email", zap.Error(err))
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Phishing email sent successfully"))
}
