package main

import (
	"crypto/tls"
	"io/ioutil"

	"github.com/pkg/errors"
)

/*
	This file is provides a versatile way of loading key pairs auth keys.

	We currently have a way of loading keypairs from local files
	stored within the root of the application.

	If one needed to load keypairs from another location, say,
	AWS parameter store, just create a new certFunc that creates
	an AWS SSM client and retrieves the keypair.

	Same for any other cloud provider.
*/

func loadFromFile(fp string) (string, error) {
	v, err := ioutil.ReadFile(fp)
	if err != nil {
		return "", errors.Wrapf(err, "read file: %s", fp)
	}

	return string(v), nil
}

type certFunc func() (string, string, error)

func loadCert() (tls.Certificate, error) {
	certFuncs := []certFunc{
		loadLocalCert,
	}

	for _, f := range certFuncs {
		cert, key, err := f()
		if err != nil {
			log.Println(err)
			continue
		}
		// Load TLS Certificate from public/private key pair
		tlscert, err := tls.X509KeyPair([]byte(cert), []byte(key))
		if err != nil {
			return tls.Certificate{}, errors.Wrap(err, "parse x509 key pair")
		}

		return tlscert, nil
	}

	return tls.Certificate{}, errors.New("couldn't load certificate")
}

func mustLoadCert() tls.Certificate {
	tlscert, err := loadCert()
	if err != nil {
		panic(err)
	}

	return tlscert
}

func loadLocalCert() (string, string, error) {
	cert, err := loadFromFile("certificate.pem")
	if err != nil {
		return "", "", errors.Wrap(err, "load local cert")
	}

	key, err := loadFromFile("key.pem")
	if err != nil {
		return "", "", errors.Wrap(err, "load local key")
	}

	return cert, key, nil
}

type authKeyFunc func() (string, error)

func loadAuthKey() (string, error) {
	keyFuncs := []authKeyFunc{
		loadLocalKey,
	}

	for _, f := range keyFuncs {
		key, err := f()
		if err != nil {
			log.Println(err)
			continue
		}

		return key, nil
	}

	return "", errors.New("couldn't load auth key")
}

func mustLoadAuthKey() string {
	key, err := loadAuthKey()
	if err != nil {
		panic(err)
	}

	return key
}

func loadLocalKey() (string, error) {
	key, err := loadFromFile("auth.pem")
	if err != nil {
		return "", errors.Wrap(err, "load local auth key")
	}

	return key, nil
}
