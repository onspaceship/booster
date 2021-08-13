package config

import (
	"errors"
	"io/ioutil"

	"github.com/apex/log"
	"github.com/spf13/viper"
)

const (
	DefaultGroundControlHost = "wss://ground-control.onspaceship.com/socket/websocket"
)

type SocketOptions struct {
	Token   string
	AgentId string

	Host      string
	Namespace string
}

func NewSocketOptions() (*SocketOptions, error) {
	options := &SocketOptions{}
	err := options.Configure()

	return options, err
}

func (options *SocketOptions) Configure() error {
	options.Token = viper.GetString("token")
	if options.Token == "" {
		log.Fatal("An agent token must be provided. Please check your Spaceship account for more details.")
	}

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

	return nil
}

func init() {
	viper.SetDefault("ground_control_host", DefaultGroundControlHost)
}
