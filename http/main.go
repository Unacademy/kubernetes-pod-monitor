package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	"github.com/unacademy/kubernetes-pod-monitor/http/health"
)

var (
	healthController *health.Controller
	router           = gin.New()
	server           *http.Server
)

// Initialize initializes http server
func Initialize() {
	healthController = health.NewController()
	router.Use(ginlogrus.Logger(log.StandardLogger()), gin.Recovery())
	setupRoutes()
}

func Shutdown() {
	if err := server.Close(); err != nil {
		log.Fatal("Server Close:", err)
	}
}

// Run starts the http server
func Run() {
	server = &http.Server{
		Addr:    ":" + viper.GetString("server.port"),
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect")
		}
	}

	log.Println("Server exiting")
}

func setupRoutes() {
	router.GET("/health", healthController.HealthHandler)
}
