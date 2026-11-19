package feature

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/RedHatInsights/composer-api/internal/config"
)

type Client struct {
	HTTP    *http.Client
	BaseURL string
}

func New(cfg config.FeatureServiceConfig) (*Client, error) {
	cert, err := tls.LoadX509KeyPair(cfg.CertPath, cfg.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("loading feature service TLS keypair: %w", err)
	}

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	if cfg.CACertPath != "" {
		caCert, err := os.ReadFile(cfg.CACertPath)
		if err != nil {
			return nil, fmt.Errorf("reading feature service CA certificate: %w", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse feature service CA certificate")
		}
		tlsCfg.RootCAs = pool
	}

	transport := &http.Transport{
		TLSClientConfig: tlsCfg,
	}

	return &Client{
		HTTP: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
		BaseURL: cfg.URL,
	}, nil
}
