package handlers

import (
	"encoding/json"
	"net/http"
	"reflect"
)

func apiResponse(w http.ResponseWriter, responseByte []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(responseByte)
}

type response struct {
	Data    interface{} `json:"data,omitempty"`
	Err     string      `json:"err,omitempty"`
	Success bool        `json:"success"`
	TTL     int         `json:"ttl"`
}

func GetSuccessResponse(data interface{}, ttl int) []byte {
	resp := &response{
		Success: true,
		TTL:     ttl,
	}

	if data == nil || (reflect.ValueOf(data).Kind() == reflect.Ptr && reflect.ValueOf(data).IsNil()) {
		resp.Data = nil
	} else {
		resp.Data = data
	}
	responseBytes, _ := json.Marshal(resp)
	return responseBytes
}

func GetErrorResponseBytes(data interface{}, ttl int, err error) []byte {
	resp := &response{
		Success: false,
		Err:     "",
		TTL:     ttl,
	}
	if data == nil || (reflect.ValueOf(data).Kind() == reflect.Ptr && reflect.ValueOf(data).IsNil()) {
		resp.Data = nil
	} else {
		resp.Data = data
	}
	if err != nil {
		resp.Err = err.Error()
	}

	responseBytes, _ := json.Marshal(resp)

	return responseBytes
}
