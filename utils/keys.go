package utils

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func ReadPublicKey() string {
	data, err := ioutil.ReadFile("/home/emanuel/go/src/github.com/emanueldf/form3-accounts/certs/test_public_key.pem")
	if err != nil {
		fmt.Print("Error on reading public key.")
	}
	return strings.TrimSpace(strings.Trim(strings.TrimSpace(string(data)), "- BEGIN PUBLIC KEY END"))
}

func ReadPrivateKey() string {
	data, err := ioutil.ReadFile("/home/emanuel/go/src/github.com/emanueldf/form3-accounts/certs/test_private_key.pem")
	if err != nil {
		fmt.Print("Error on reading private key.")
	}
	return strings.TrimSpace(strings.Trim(strings.TrimSpace(string(data)), "- BEGIN RSA PRIVATE KEY END"))
}
