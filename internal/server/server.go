package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"

	"github.com/upper-institute/flipbook/internal/helpers"
)

const (
	keyBitSize_flag     = "tls.key.bit.size"
	certDuration_flag   = "tls.cert.duration"
	enableTls_flag      = "tls.enable"
	certPath_flag       = "tls.cert.path"
	privateKeyPath_flag = "tls.key.path"
	serverAddress_flag  = "server.address"
)

func BindOptions(binder helpers.FlagBinder) {

	binder.BindInt64(keyBitSize_flag, 2048, "Bit size for private key (RSA 1024, 2048 or 4096)")
	binder.BindDuration(certDuration_flag, 315360000*time.Second, "Bit size for private key (RSA 1024, 2048 or 4096)")
	binder.BindBool(enableTls_flag, true, "Enable TLS on server")
	binder.BindString(certPath_flag, "", "Path to TLS certificate to use, if not provided a in memory certificate will be generated")
	binder.BindString(privateKeyPath_flag, "", "Path to private key for TLS encipherment")
	binder.BindString(serverAddress_flag, "0.0.0.0:6333", "Bind address for gRPC server listener")

}

func generateInMemoryTlsCertificate(getter helpers.FlagGetter) tls.Certificate {

	privateKey, err := rsa.GenerateKey(rand.Reader, int(getter.GetInt64(keyBitSize_flag)))
	if err != nil {
		panic(err)
	}

	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	caCert := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Organization:       []string{"Upper Institute"},
			Country:            []string{"BR"},
			Province:           []string{},
			Locality:           []string{"São Paulo"},
			OrganizationalUnit: []string{"flipbook"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(getter.GetDuration(certDuration_flag)),
		IsCA:      true,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
		KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment,
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, caCert, caCert, &privateKey.PublicKey, privateKey)
	if err != nil {
		panic(err)
	}

	caBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}

	tlsCert, err := tls.X509KeyPair(
		pem.EncodeToMemory(caBlock),
		pem.EncodeToMemory(privateKeyBlock),
	)
	if err != nil {
		panic(err)
	}

	return tlsCert

}

func loadTlsConfig(getter helpers.FlagGetter) *tls.Config {

	if !getter.GetBool(enableTls_flag) {
		return nil
	}

	certPath := getter.GetString(certPath_flag)

	var (
		cert tls.Certificate
		err  error
	)

	if len(certPath) == 0 {
		cert = generateInMemoryTlsCertificate(getter)
	} else {
		cert, err = tls.LoadX509KeyPair(
			certPath,
			getter.GetString(privateKeyPath_flag),
		)
		if err != nil {
			panic(err)
		}
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		ClientAuth:         tls.RequestClientCert,
		NextProtos:         []string{"h2"},
	}

	return tlsConfig

}

func CreateListener(getter helpers.FlagGetter) net.Listener {

	addr := getter.GetString(serverAddress_flag)

	if getter.GetBool(enableTls_flag) {

		config := loadTlsConfig(getter)

		lis, err := tls.Listen("tcp", addr, config)
		if err != nil {
			panic(err)
		}

		return lis

	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	return lis

}
