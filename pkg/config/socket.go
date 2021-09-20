package config

import (
	"errors"
	"io/ioutil"

	"github.com/spf13/viper"
)

const (
	DefaultGroundControlHost   = "wss://ground-control.onspaceship.com/socket/websocket"
	DefaultBuildRegistryURL    = "us.gcr.io/onspaceship"
	DefaultBuildServiceAccount = "booster"
)

type SocketOptions struct {
	AgentId string

	Host      string
	Namespace string

	BuildRegistryURL    string
	BuildServiceAccount string
}

func NewSocketOptions() (*SocketOptions, error) {
	options := &SocketOptions{}
	err := options.Configure()

	return options, err
}

func (options *SocketOptions) Configure() error {
	options.AgentId = viper.GetString("agent_id")

	options.Host = viper.GetString("ground_control_host")
	if options.Host == "" {
		return errors.New("invalid ground_control_host configuration")
	}

	namespace, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		namespace = []byte(viper.GetString("namespace"))
	}
	options.Namespace = string(namespace)

	options.BuildRegistryURL = viper.GetString("build_registry_url")
	if options.BuildRegistryURL == "" {
		return errors.New("invalid build_registry_url configuration")
	}

	options.BuildServiceAccount = viper.GetString("build_service_account")
	if options.BuildServiceAccount == "" {
		return errors.New("invalid build_service_account configuration")
	}

	return nil
}

func init() {
	viper.SetDefault("ground_control_host", DefaultGroundControlHost)
	viper.SetDefault("build_registry_url", DefaultBuildRegistryURL)
	viper.SetDefault("build_service_account", DefaultBuildServiceAccount)
}
