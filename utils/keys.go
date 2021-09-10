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
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Sign(fname string, content string, hashFunc crypto.Hash) string {
	key := getPrivateKey(getCertsPath(fname))
	data := []byte(strings.Join(flag.Args(), content))
	sign, errSign := rsa.SignPKCS1v15(rand.Reader, key, hashFunc, Hash(hashFunc, data))
	if errSign != nil {
		exit(errSign)
	}
	return base64.StdEncoding.EncodeToString(sign)
}

func getCertsPath(fname string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]) + "/certs")
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(dir, fname)
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

func Hash(hashFunc crypto.Hash, data []byte) []byte {
	h := hashFunc.New()
	h.Write(data)
	return h.Sum(nil)
}

func GetPublicKeyContent(fname string) string {
	data, err := ioutil.ReadFile(getCertsPath(fname))
	if err != nil {
		fmt.Print("Error on reading public key.")
	}

	return (strings.TrimSpace(strings.Trim(strings.TrimSpace(string(data)), "- BEGIN PUBLIC KEY END")))
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
