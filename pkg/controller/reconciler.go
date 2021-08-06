package controller

import (
	"context"

	"github.com/apex/log"
	kpackapi "github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler struct {
	client.Client
	*Options
	Scheme *runtime.Scheme
}

func (rec *Reconciler) Reconcile(_ context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.WithField("build", req.NamespacedName)
	logger.Info("Reconciling")

	ctx := context.Background()

	build := &kpackapi.Build{}
	if err := rec.Get(ctx, req.NamespacedName, build); err != nil {
		log.WithError(err).Error("Unable to get build")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}
