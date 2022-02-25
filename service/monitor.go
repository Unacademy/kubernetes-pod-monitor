package service

import (
	"bytes"
	"io"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
	"github.com/unacademy/kubernetes-pod-monitor/sessions"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ContainerInfo struct {
	PodName       string `gorm:"column:podname"`
	Namespace     string `gorm:"column:namespace"`
	ContainerName string `gorm:"column:containername"`
	ClusterName   string `gorm:"column:clustername"`
	Reason        string `gorm:"column:reason"`
}

type RestartValue struct {
	RestartCount int32 `gorm:"column:restartcount"`
	Retries      int32 `gorm:"column:retries"`
}

var (
	concurrencyCh chan bool
	clusterName   string
	tableName     string
)

func checkPod(pod *corev1.Pod) {
	clientset := sessions.GetClientset()
	namespace := pod.GetNamespace()
	podName := pod.GetName()
	containers := pod.Status.ContainerStatuses

	for _, container := range containers {
		containerInfo := ContainerInfo{
			PodName:       podName,
			Namespace:     namespace,
			ContainerName: container.Name,
			ClusterName:   clusterName,
		}
		restartInfo := getSavedObject(&containerInfo)
		if container.RestartCount <= restartInfo.RestartCount {
			continue
		} else if restartInfo.Retries >= viper.GetInt32("max_retries") {
			restartInfo = &RestartValue{
				RestartCount: container.RestartCount,
				Retries:      0,
			}
			updateObject(&containerInfo, restartInfo)
			continue
		} else {
			restartInfo = &RestartValue{
				RestartCount: restartInfo.RestartCount,
				Retries:      restartInfo.Retries + 1,
			}
			updateObject(&containerInfo, restartInfo)
		}

		log.Infof("Pod: %s in Namespace: %s crashed.", podName, namespace)

		lastState := container.LastTerminationState.Terminated
		lastTerminationState := "Not Available"
		if lastState != nil {
			lastTerminationState = lastState.String()
		}

		podLogOpts := corev1.PodLogOptions{Container: container.Name, Previous: true}
		req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &podLogOpts)

		podLogs, err := req.Stream()
		if err != nil {
			log.Errorln("failed to get pods:", err)
			continue
		}

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, podLogs)
		if err != nil {
			log.Errorln("error in copy information from podLogs to buf:", err)
			continue
		}
		logs := buf.String()
		containerInfo.Reason = container.LastTerminationState.Terminated.Reason

		sendESLogs(&containerInfo, container.RestartCount, &logs, &lastTerminationState)

		percent := getPercentPodCrash(namespace, clientset, container.Name)
		notifyOnSlack(&containerInfo, &lastTerminationState, percent)

		persistPodCrash(&containerInfo, container.RestartCount)
		podLogs.Close()

		restartInfo = &RestartValue{
			RestartCount: container.RestartCount,
			Retries:      0,
		}
		updateObject(&containerInfo, restartInfo)
	}
	<-concurrencyCh
}

func Execute() {
	concurrencyCh = make(chan bool, viper.GetInt64("max_parallelism"))
	clientset := sessions.GetClientset()
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get pods:", err)
	}
	for _, pod := range pods.Items {
		podCopy := pod.DeepCopy()
		concurrencyCh <- true
		go checkPod(podCopy)
	}
}
