package service

import (
	"sync"
	"time"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
	"github.com/unacademy/kubernetes-pod-monitor/sessions"
)

var (
	serviceDone      chan bool
	serviceWaitGroup sync.WaitGroup
)

func Initialize() {
	sessions.InitClientset()
	InitializeLogger()
	serviceDone = make(chan bool)
	clusterName = viper.GetString("CLUSTER_NAME")
	tableName = viper.GetString("sql.table_name")
	StartLogger()
}

func Run() {
	serviceWaitGroup.Add(1)
	go startService()
}

func Shutdown() {
	log.Info("not checking pods anymore ...")
	serviceDone <- true
	serviceWaitGroup.Wait()
	ShutdownLogger()
}

func startService() {
	defer serviceWaitGroup.Done()
	Execute()

	for {
		select {
		case <-time.After(time.Duration(viper.GetInt64("check_interval")) * time.Second):
			log.Info("checking if any pods crashed...")
			Execute()
			log.Info("finished checking pods...")

		case <-serviceDone:
			log.Info("exiting the monitor.")
			return
		}
	}
}
