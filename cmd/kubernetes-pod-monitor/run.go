package kubernetes_pod_monitor

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/unacademy/kubernetes-pod-monitor/http"
	"github.com/unacademy/kubernetes-pod-monitor/service"
	"github.com/unacademy/kubernetes-pod-monitor/sessions"

	log "github.com/sirupsen/logrus"
)

var done = make(chan bool)
var gracefulStop = make(chan os.Signal)

func setup() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	viper.SetConfigName("application")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error in config file: %s \n", err))
	}
	log.SetFormatter(&log.JSONFormatter{})
	gin.SetMode(gin.ReleaseMode)
	switch viper.GetString("log.level") {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
		break
	case "INFO":
		log.SetLevel(log.InfoLevel)
		break
	default:
		log.SetLevel(log.ErrorLevel)
		break
	}
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)
	setupApp()
}

func cleanupOnSignal(cleanup func()) {
	go func() {
		sig := <-gracefulStop
		log.Info(fmt.Sprintf("caught sig: %+v. waiting for goroutines to finish", sig))
		cleanup()
		log.Info("goroutines finished. exiting")
		os.Exit(0)
	}()
}

func setupApp() {
	service.Initialize()
	http.Initialize()
	sessions.HealthOrPanic()
	cleanupOnSignal(cleanup)
}

func cleanup() {
	service.Shutdown()
	http.Shutdown()
	done <- true
}

func Run() {
	setup()
	go http.Run()
	go service.Run()
	<-done
}
