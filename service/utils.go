package service

import (
	"fmt"
	"time"

	"github.com/bluele/slack"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/jinzhu/gorm"

	log "github.com/sirupsen/logrus"

	"github.com/unacademy/kubernetes-pod-monitor/sessions"
)

const createdFormat = "2006-01-02 15:04:05" //"Jan 2, 2006 at 3:04pm (MST)"

var defaultSlackChannel = viper.GetString("slack.channel")

// getCurrentTimeInSeconds returns current time in seconds.
func getCurrentTimeInSeconds() int64 {
	return time.Now().UnixNano() / int64(time.Second)
}

// getCurrentTimeInMillis returns current time in milliseconds.
func getCurrentTimeInMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func getSavedObject(obj *ContainerInfo) *RestartValue {
	sqlClient := sessions.GetSqlClient()
	row := RestartValue{}
	err := sqlClient.Raw(`SELECT restartcount, retries FROM `+tableName+` WHERE clustername=? AND namespace=? AND podname=? AND containername=?`,
		obj.ClusterName, obj.Namespace, obj.PodName, obj.ContainerName).First(&row).Error
	switch {
	case err == gorm.ErrRecordNotFound:
	case err != nil:
		log.Errorf("query error: %v\n", err)
	}
	return &row
}

func updateObject(obj *ContainerInfo, newVal *RestartValue) {
	sqlClient := sessions.GetSqlClient()
	currTime := getCurrentTimeInSeconds()
	err := sqlClient.Exec(`INSERT INTO `+tableName+` (clustername, namespace, podname, containername, restartcount, retries, edited_at) VALUES(?,?,?,?,?,?,?)
	ON DUPLICATE KEY UPDATE restartcount=?, retries=?, edited_at=?`,
		obj.ClusterName, obj.Namespace, obj.PodName, obj.ContainerName, newVal.RestartCount, newVal.Retries, currTime,
		newVal.RestartCount, newVal.Retries, currTime).Error

	if err != nil {
		log.Errorf("update query error: %v\n", err)
	}
}

func sendESLogs(containerInfo *ContainerInfo, restartCount int32, logs *string, lastTerminationState *string) {
	max_length := viper.GetInt64("max_crashlog_length")
	log_len := int64(len(*logs))
	if log_len >= max_length {
		*logs = (*logs)[log_len-max_length:]
	}
	LogChannel <- LogMessage{
		Namespace:        containerInfo.Namespace,
		PodName:          containerInfo.PodName,
		ContainerName:    containerInfo.ContainerName,
		CreatedAt:        getCurrentTimeInMillis(),
		ClusterName:      containerInfo.ClusterName,
		Logs:             *logs,
		RestartCount:     restartCount,
		TerminationState: *lastTerminationState,
	}
}

func notifyOnSlack(containerInfo *ContainerInfo, lastTerminationState *string, percent int) {
	token := viper.GetString("slack.channel")
	if shouldNotify(containerInfo.ClusterName, containerInfo.Namespace, containerInfo.ContainerName) == false || viper.GetBool("slack.notify") == false {
		log.Info("not notifying...")
		return
	}
	channelName := getSlackChannel(containerInfo.ClusterName, containerInfo.Namespace)
	api := slack.New(token)

	log.Info("notifying...")
	link := viper.GetString("elasticsearch.dashboard")

	msg := fmt.Sprintf("*Cluster Name*:- %s\n*Namespace*:- %s\n*Container Name*:- %s\n *Reason*:- %s\n `Kibana dashboard`:- <%s|Dashboard>",
		containerInfo.ClusterName, containerInfo.Namespace, containerInfo.ContainerName, containerInfo.Reason, link)

	err := api.ChatPostMessage(channelName, msg, nil)
	if err != nil {
		log.Error(err)
	}
}

func getSlackChannel(clusterName string, namespace string) string {
	sqlClient := sessions.GetSqlClient()
	var slackChannel string

	rows, err := sqlClient.Raw(`select slack_channel FROM k8s_pod_crash_notify WHERE clustername=? AND namespace=?`,
		clusterName, namespace).Rows()
	if err != nil {
		log.Error(err)
		return defaultSlackChannel
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&slackChannel)
		if err != nil {
			log.Error(err)
			return defaultSlackChannel
		}
		return slackChannel
	}
	return defaultSlackChannel
}

func shouldNotify(clusterName string, namespace string, containername string) bool {
	sqlClient := sessions.GetSqlClient()
	var exists string

	rows, err := sqlClient.Raw(`select count(*) FROM k8s_crash_ignore_notify WHERE clustername=? AND namespace=? AND containername=?`,
		clusterName, namespace, containername).Rows()
	if err != nil {
		log.Error(err)
		return true
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&exists)
		if err != nil {
			log.Error(err)
			return true
		}
		if exists == "1" {
			return false
		}
	}
	return true
}

func persistPodCrash(containerInfo *ContainerInfo, restartCount int32) {
	sqlClient := sessions.GetSqlClient()
	currTime := time.Unix(getCurrentTimeInSeconds(), 0).Format(createdFormat)
	err := sqlClient.Exec(`INSERT INTO k8s_pod_crash (clustername, namespace, containername, restartcount, date) VALUES(?,?,?,?,?)`,
		containerInfo.ClusterName, containerInfo.Namespace, containerInfo.ContainerName, restartCount, currTime).Error

	if err != nil {
		log.Errorf("update query error: %v\n", err)
	}
}

func getPercentPodCrash(namespace string, clientset *kubernetes.Clientset, containerName string) int {
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get pods:", err)
	}

	podCount := 0
	terminatedPods := 0
	for _, pod := range pods.Items {
		podCopy := pod.DeepCopy()

		for _, _container := range podCopy.Status.ContainerStatuses {
			_containerName := _container.Name
			if _containerName == containerName {
				podCount++
			}
			if _container.LastTerminationState.Terminated != nil {
				terminatedPods++
			}
		}
	}

	return (terminatedPods * 100) / podCount
}
