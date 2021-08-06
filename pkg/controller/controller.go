package controller

import (
	"github.com/onspaceship/booster/pkg/config"

	"github.com/apex/log"
	kpackapi "github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Options = config.ControllerOptions

func StartController(exit chan<- bool) {
	log.Info("Controller Lift Off!")

	kubeConfig := ctrl.GetConfigOrDie()
	ctrlCtx := ctrl.SetupSignalHandler()

	options, err := config.NewControllerOptions()
	if err != nil {
		log.WithError(err).Fatal("failed to configure controller")
	}

	mgr, err := ctrl.NewManager(kubeConfig, ctrl.Options{
		LeaderElection:     options.LeaderElection,
		LeaderElectionID:   "booster",
		MetricsBindAddress: "0",
	})
	if err != nil {
		log.WithError(err).Fatal("Unable to create manager")
	}

	err = kpackapi.AddToScheme(mgr.GetScheme())
	if err != nil {
		log.WithError(err).Fatal("Unable to add scheme")
	}

	err = ctrl.NewControllerManagedBy(mgr).
		For(&kpackapi.Build{}).
		Complete(&Reconciler{
			Options: options,
			Client:  mgr.GetClient(),
			Scheme:  mgr.GetScheme(),
		})
	if err != nil {
		log.WithError(err).Fatal("Unable to create controller")
	}

	log.Info("Starting manager")
	if err := mgr.Start(ctrlCtx); err != nil {
		log.WithError(err).Fatal("Problem running manager")
	}

	exit <- true
}
