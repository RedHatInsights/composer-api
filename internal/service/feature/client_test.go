package feature

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/RedHatInsights/composer-api/internal/config"
)

func writeSelfSignedCert(t *testing.T, dir string) (certPath, keyPath string) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}

	certPath = filepath.Join(dir, "tls.crt")
	certFile, err := os.Create(certPath)
	if err != nil {
		t.Fatal(err)
	}
	defer certFile.Close()
	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		t.Fatal(err)
	}

	keyPath = filepath.Join(dir, "tls.key")
	keyFile, err := os.Create(keyPath)
	if err != nil {
		t.Fatal(err)
	}
	defer keyFile.Close()
	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	if err := pem.Encode(keyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes}); err != nil {
		t.Fatal(err)
	}

	return certPath, keyPath
}

func TestNew(t *testing.T) {
	dir := t.TempDir()
	certPath, keyPath := writeSelfSignedCert(t, dir)

	client, err := New(config.FeatureServiceConfig{
		URL:      "https://feature.example.com",
		CertPath: certPath,
		KeyPath:  keyPath,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.BaseURL != "https://feature.example.com" {
		t.Errorf("expected BaseURL %q, got %q", "https://feature.example.com", client.BaseURL)
	}
	if client.HTTP == nil {
		t.Fatal("expected non-nil HTTP client")
	}
}

func TestNew_WithCACert(t *testing.T) {
	dir := t.TempDir()
	certPath, keyPath := writeSelfSignedCert(t, dir)

	caPath := filepath.Join(dir, "ca.pem")
	certPEM, _ := os.ReadFile(certPath)
	if err := os.WriteFile(caPath, certPEM, 0644); err != nil {
		t.Fatal(err)
	}

	client, err := New(config.FeatureServiceConfig{
		URL:        "https://feature.example.com",
		CertPath:   certPath,
		KeyPath:    keyPath,
		CACertPath: caPath,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.HTTP == nil {
		t.Fatal("expected non-nil HTTP client")
	}
}

func TestNew_InvalidCert(t *testing.T) {
	_, err := New(config.FeatureServiceConfig{
		URL:      "https://feature.example.com",
		CertPath: "/nonexistent/tls.crt",
		KeyPath:  "/nonexistent/tls.key",
	})
	if err == nil {
		t.Fatal("expected error for invalid cert path")
	}
}

func TestNew_InvalidCACert(t *testing.T) {
	dir := t.TempDir()
	certPath, keyPath := writeSelfSignedCert(t, dir)

	_, err := New(config.FeatureServiceConfig{
		URL:        "https://feature.example.com",
		CertPath:   certPath,
		KeyPath:    keyPath,
		CACertPath: "/nonexistent/ca.pem",
	})
	if err == nil {
		t.Fatal("expected error for invalid CA cert path")
	}
}
