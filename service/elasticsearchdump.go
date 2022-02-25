package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"

	elastic "github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/unacademy/kubernetes-pod-monitor/config"
	"github.com/unacademy/kubernetes-pod-monitor/sessions"
)

var (
	LogChannel      chan LogMessage
	LogBufferSize   int = 100
	LogBatchSize    int = 10
	LogFlushSeconds int = 30
	logDone         chan bool
	logWaitGroup    sync.WaitGroup
)

type LogMessage struct {
	Namespace        string `json:"namespace"`
	PodName          string `json:"pod_name"`
	ContainerName    string `json:"container_name"`
	CreatedAt        int64  `json:"created_at"`
	ClusterName      string `json:"cluster_name"`
	Logs             string `json:"logs"`
	RestartCount     int32  `json:"restart_count"`
	TerminationState string `json:"termination_state"`
}

func CheckAndCreateESIndex(indexName string) error {
	client := sessions.GetElasticSearchClient()
	ctx := context.Background()

	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		log.Infof("Creating ES Log index: %s", indexName)
		docType := viper.GetString("elasticsearch.doc_type")

		var createIndex *elastic.IndicesCreateResult
		var err error

		if viper.GetBool("elasticsearch.v7") {
			createIndex, err = client.CreateIndex(indexName).BodyString(config.Mappingv7).Do(ctx)
		} else {
			createIndex, err = client.CreateIndex(indexName).BodyString(fmt.Sprintf(config.Mapping, docType)).Do(ctx)
		}
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			log.Warnf("ES Log Index creation not acknowledged: %v", createIndex)
		}
	}
	return nil
}

func BulkInsertES(logBatch []LogMessage) {
	log.Infof("Indexing %d logs in ES Client.", len(logBatch))
	indexName := viper.GetString("elasticsearch.index") + time.Now().Format("2006.01.02")
	err := CheckAndCreateESIndex(indexName)
	if err != nil {
		log.Errorf("Error in checking and creating Index in ES: %s", err)
		return
	}
	client := sessions.GetElasticSearchClient()
	ctx := context.Background()

	bulkRequest := client.Bulk()

	for _, msg := range logBatch {
		if viper.GetBool("elasticsearch.v7") {
			bulkRequest = bulkRequest.Add(elastic.NewBulkIndexRequest().Index(indexName).Doc(msg))
		} else {
			docType := viper.GetString("elasticsearch.doc_type")
			bulkRequest = bulkRequest.Add(elastic.NewBulkIndexRequest().Index(indexName).Type(docType).Doc(msg))
		}
	}

	_, err = bulkRequest.Do(ctx)
	if err != nil {
		log.Error(err)
	}
}

func LogMessageBatcher() {
	defer logWaitGroup.Done()
	logBatch := make([]LogMessage, 0, LogBatchSize)
	flushLogs := func(size int) {
		if len(logBatch) >= size {
			BulkInsertES(logBatch)
			logBatch = make([]LogMessage, 0, LogBatchSize)
		}
	}

	for {
		select {
		case msg := <-LogChannel:
			logBatch = append(logBatch, msg)
			flushLogs(LogBatchSize)

		case <-time.After(time.Duration(LogFlushSeconds) * time.Second):
			log.Info("Timeout hence sending logs to ES.")
			flushLogs(1)

		case <-logDone:
			log.Info("exiting the Logger.")
			flushLogs(1)
			return
		}
	}
}

func InitializeLogger() {
	LogBufferSize = viper.GetInt("elasticsearch.buffer_size")
	LogBatchSize = viper.GetInt("elasticsearch.batch_size")
	LogFlushSeconds = viper.GetInt("elasticsearch.flush_seconds")
	LogChannel = make(chan LogMessage, LogBufferSize)
	logDone = make(chan bool)
}

func StartLogger() {
	logWaitGroup.Add(1)
	go LogMessageBatcher()
}

func ShutdownLogger() {
	log.Info("not reading any further messages from Log ...")
	logDone <- true
	close(LogChannel)
	logWaitGroup.Wait()
	log.Info("all log messages processed. exiting ...")
}
