package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ReadRequestContent() string {
	data, err := ioutil.ReadFile("/home/emanuel/go/src/github.com/emanueldf/form3-accounts/certs/request_signing_test_request_content.txt")
	if err != nil {
		fmt.Print("Error on reading request content.")
	}
	return strings.TrimSpace(string(data))
}

func ReadSignature() string {
	data, err := ioutil.ReadFile("/home/emanuel/go/src/github.com/emanueldf/form3-accounts/certs/request_signing_test_signature.txt")
	if err != nil {
		fmt.Print("Error on reading signature.")
	}
	return strings.TrimSpace(string(data))
}

func hash(hashFunc crypto.Hash, data []byte) []byte {
	h := hashFunc.New()
	h.Write(data)
	return h.Sum(nil)
}

func Sign(fname string, hashFunc crypto.Hash) string {
	key := getPrivateKey(fname)
	data := []byte(strings.Join(flag.Args(), " "))
	sign, errSign := rsa.SignPKCS1v15(rand.Reader, key, hashFunc, hash(hashFunc, data))
	if errSign != nil {
		exit(errSign)
	}
	return base64.StdEncoding.EncodeToString(sign)
}

func ReadPublicKey() string {
	data, err := ioutil.ReadFile("/home/emanuel/go/src/github.com/emanueldf/form3-accounts/certs/test_public_key.pem")
	if err != nil {
		fmt.Print("Error on reading public key.")
	}

	return (strings.TrimSpace(strings.Trim(strings.TrimSpace(string(data)), "- BEGIN PUBLIC KEY END")))
}

func getPublicKey(fname string) *rsa.PublicKey {
	buf, errRead := ioutil.ReadFile(fname)
	if errRead != nil {
		exit(errRead)
	}
	block, _ := pem.Decode(buf)
	if block == nil {
		exit(errors.New("public key error"))
	}
	pub, errParse := x509.ParsePKIXPublicKey(block.Bytes)
	if errParse != nil {
		exit(errParse)
	}
	return pub.(*rsa.PublicKey)
}

func ReadPrivateKeyContent() string {
	data, err := ioutil.ReadFile("/home/emanuel/go/src/github.com/emanueldf/form3-accounts/certs/test_private_key.pem")
	if err != nil {
		fmt.Print("Error on reading private key.")
	}
	return (strings.TrimSpace(strings.Trim(strings.TrimSpace(string(data)), "- BEGIN RSA PRIVATE KEY END")))
}

func ReadPrivateKey() rsa.PrivateKey {
	data, err := ioutil.ReadFile("/home/emanuel/go/src/github.com/emanueldf/form3-accounts/certs/test_private_key.pem")
	if err != nil {
		fmt.Print("Error on reading private key.")
	}
	return *getPrivateKey(strings.TrimSpace(strings.Trim(strings.TrimSpace(string(data)), "- BEGIN RSA PRIVATE KEY END")))
}

func getPrivateKey(fname string) *rsa.PrivateKey {
	buf, errRead := ioutil.ReadFile(fname)
	if errRead != nil {
		exit(errRead)
	}
	block, _ := pem.Decode(buf)
	if block == nil {
		exit(errors.New("private key error"))
	}
	key, errParse := x509.ParsePKCS1PrivateKey(block.Bytes)
	if errParse != nil {
		exit(errParse)
	}
	return key
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
