package utils

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"
)

var (
	MYTH      string
	MYSTIC    string
	JWTISSUER string
	Client_ID string
	MYST      string
	ANU       string
)

var Logger, _ = zap.NewDevelopment()

// Middleware to log requests to the server ##
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// logger := NewLogger(os.Stdout, os.Stderr)
		Logger.Info((fmt.Sprintf("%v - %v %v %v", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())))
		next.ServeHTTP(w, r)
	})
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				//if panic close connection
				w.Header().Set("Connection", "Close")
				//write internal server error
				ServerError(w, "Connection Closed inabruptly", fmt.Errorf("%v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// used for all internal server Error
func ServerError(w http.ResponseWriter, errMsg string, err error) {
	fmt.Println("Reaxcher 1")
	errTrace := fmt.Sprintf("%v\n%v", err.Error(), debug.Stack())
	fmt.Println("Reaxcher 2")
	Logger.Error(errTrace)
	fmt.Println("Reaxcher 3")
	http.Error(w, errMsg, http.StatusInternalServerError)
	fmt.Println("Reaxcher 4")
}