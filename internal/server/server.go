package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/h3th-IV/mackerel/internal/database"
	"github.com/h3th-IV/mackerel/internal/utils"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

type GracefulShutdownServer struct {
	HTTPListenAddr     string
	MCSAttackHandler   http.Handler
	CaptureDataHandler http.Handler
	StealDataHandler   http.Handler
	httpServer         *http.Server
	WriteTimeout       time.Duration
	ReadTimeout        time.Duration
	IdleTimeout        time.Duration
	HandlerTimeout     time.Duration
	Logger             *zap.Logger
	Mysqlclient        database.Database
}

func NewGracefulShutdownServer(listenaddr string, logger *zap.Logger, mysqlclient database.Database) *GracefulShutdownServer {
	return &GracefulShutdownServer{
		Logger:      logger,
		Mysqlclient: mysqlclient,
	}
}

func (server *GracefulShutdownServer) getRouter() *mux.Router {
	router := mux.NewRouter()

	mux.CORSMethodMiddleware(router)
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	middleWareChain := alice.New(utils.RequestLogger, utils.RecoverPanic, cors.Handler)
	//api routes here
	router.Handle("/mcs-attack", server.MCSAttackHandler).Methods(http.MethodPost)
	router.Handle("/capture-data", server.CaptureDataHandler).Methods(http.MethodPost)
	router.Use(middleWareChain.Then) //request logging will be handled here
	mux.CORSMethodMiddleware(router)
	router.SkipClean(true)
	return router
}

func (server *GracefulShutdownServer) Start() {
	router := server.getRouter()
	server.httpServer = &http.Server{
		Addr:         server.HTTPListenAddr,
		WriteTimeout: server.WriteTimeout,
		ReadTimeout:  server.ReadTimeout,
		IdleTimeout:  server.IdleTimeout,
		Handler:      router,
	}
	utils.Logger.Info(fmt.Sprintf("listening and serving on %s", server.HTTPListenAddr))
	if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		utils.Logger.Fatal("server failed to start", zap.Error(err))
	}
}
