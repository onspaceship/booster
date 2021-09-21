package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/url"

	"github.com/spf13/viper"
)

const (
	DefaultCoreBaseURL = "https://core.onspaceship.com/"
)

type ClientOptions struct {
	CoreBaseURL *url.URL
	CoreJWTKey  *rsa.PrivateKey
}

func NewClientOptions() (*ClientOptions, error) {
	options := &ClientOptions{}
	err := options.Configure()

	return options, err
}

func (options *ClientOptions) Configure() error {
	coreBaseURL, err := url.Parse(viper.GetString("core_base_url"))
	if err != nil {
		return errors.New("invalid core_base_url")
	}
	options.CoreBaseURL = coreBaseURL

	privPem, _ := pem.Decode([]byte(viper.GetString("core_jwt_key")))
	signKey, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if err != nil {
		return errors.New("invalid core_jwt_key")
	}

	options.CoreJWTKey = signKey

	return nil
}

func init() {
	viper.SetDefault("core_base_url", DefaultCoreBaseURL)
}
