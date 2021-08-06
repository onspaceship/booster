package config

import (
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

type ControllerOptions struct {
	Namespace      string
	LeaderElection bool
}

func NewControllerOptions() (*ControllerOptions, error) {
	options := &ControllerOptions{}
	err := options.Configure()

	return options, err
}

func (options *ControllerOptions) Configure() error {
	namespace, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		namespace = []byte(viper.GetString("namespace"))
	}
	options.Namespace = string(namespace)

	_, leaderElection := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	options.LeaderElection = leaderElection

	return nil
}
