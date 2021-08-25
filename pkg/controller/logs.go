package controller

import (
	"context"
	"io"
	"strings"

	"github.com/apex/log"
	buildapi "github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (rec *Reconciler) getBuildLogs(build *buildapi.Build) string {
	ctx := context.Background()

	var readyContainers []string
	pod, err := rec.client.CoreV1().Pods(build.Namespace).Get(ctx, build.Status.PodName, metav1.GetOptions{})
	if err != nil {
		log.WithError(err).Info("Unable to get build pod")
		return ""
	}

	for _, container := range pod.Status.InitContainerStatuses {
		if container.State.Waiting == nil && (container.Name == "detect" || container.Name == "build") {
			readyContainers = append(readyContainers, container.Name)
		}
	}

	for _, container := range pod.Status.ContainerStatuses {
		if container.State.Waiting == nil {
			readyContainers = append(readyContainers, container.Name)
		}
	}

	var logs strings.Builder
	for _, container := range readyContainers {
		logStream, err := rec.client.CoreV1().Pods(build.Namespace).GetLogs(build.Status.PodName, &corev1.PodLogOptions{Container: container}).Stream(ctx)
		if err != nil {
			log.WithError(err).WithField("container", container).Info("Unable to get build logs")
			return ""
		}

		defer logStream.Close()

		buf, err := io.ReadAll(logStream)
		if err != nil {
			log.WithError(err).WithField("container", container).Info("Unable to read build log")
			return ""
		}

		logs.WriteString(string(buf))
	}

	return logs.String()
}
