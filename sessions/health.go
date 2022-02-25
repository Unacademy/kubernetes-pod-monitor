package sessions

import (
	"context"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func HealthOrPanic() {
	log.Info("Checking Health ...")

	log.Info("Checking Elasticsearch ...")
	esClient := GetElasticSearchClient()
	ctx := context.Background()
	log.Info("Got Elasticsearch client...")
	_, err := esClient.CatIndices().Do(ctx)
	if err != nil {
		panic(err)
	}
	log.Info("Got Elasticsearch response...")

	clientset := GetClientset()
	log.Info("Got Clientset...")
	_, err = clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	log.Info("Got Clientset response...")

	sqlClient := GetSqlClient()
	log.Info("Got SqlClient...")
	err = sqlClient.DB().Ping()
	if err != nil {
		panic(err)
	}
	log.Info("Got SqlClient response...")
}
