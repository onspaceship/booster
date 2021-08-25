package controller

import (
	"context"

	"github.com/onspaceship/booster/pkg/client"

	"github.com/apex/log"
	buildapi "github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	kpackapi "github.com/pivotal/kpack/pkg/apis/core/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler struct {
	ctrlclient.Client
	*Options
	Scheme *runtime.Scheme

	client *kubernetes.Clientset
}

type BuildStatus string

const (
	BuildStatusBuilding BuildStatus = "building"
	BuildStatusComplete BuildStatus = "complete"
	BuildStatusError    BuildStatus = "error"
)

func (rec *Reconciler) Reconcile(_ context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.WithField("build", req.NamespacedName)
	logger.Info("Reconciling")

	ctx := context.Background()

	build := &buildapi.Build{}
	if err := rec.Get(ctx, req.NamespacedName, build); err != nil {
		log.WithError(err).Error("Unable to get build")
		return ctrl.Result{}, ctrlclient.IgnoreNotFound(err)
	}

	status := getStatus(build)

	err := client.NewClient().CoreBuildUpdate(build.Annotations[BuildIdAnnotation], string(status), build.Status.LatestImage)
	if err != nil {
		log.WithError(err).Error("Unable to update Core")
		return ctrl.Result{Requeue: false}, err
	}

	logger.WithField("status", status).WithField("image", build.Status.LatestImage).Info("Build updated!")

	if build.Status.PodName != "" {
		logs := rec.getBuildLogs(build)

		if logs != "" {
			err := client.NewClient().CoreBuildLogsUpdate(build.Annotations[BuildIdAnnotation], logs)

			if err != nil {
				log.WithError(err).Error("Unable to send logs to Core")
				return ctrl.Result{Requeue: false}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

func getStatus(build *buildapi.Build) BuildStatus {
	cond := build.Status.GetCondition(kpackapi.ConditionSucceeded)
	switch {
	case cond.IsTrue():
		return BuildStatusComplete
	case cond.IsFalse():
		return BuildStatusError
	case cond.IsUnknown():
		return BuildStatusBuilding
	default:
		return BuildStatusError
	}
}
