package controllers

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/emanueldf/form3-accounts/models"
	"github.com/emanueldf/form3-accounts/utils"
)

var (
	privateKeyFilePath = "/certs/test_private_key.pem"
	publicKeyFilePath  = "/certs/test_public_key.pem"
)

func Init() {
	path := "/v1/organisation/accounts"
	host := "api.form3.tech"
	baseUrl := "https://" + host
	fullPath := baseUrl + path

	var Country = "GB"

	var payload = models.NewAccountData{
		Type:           "accounts",
		ID:             utils.Generate(),
		OrganisationID: utils.Generate(),
		Attributes: &models.NewAccountAttributes{
			Country:                &Country,
			BaseCurrency:           "GBP",
			BankID:                 "400302",
			BankIDCode:             "GBDSC",
			Bic:                    "NWBKGB42",
			ProcessingService:      "ABC Bank",
			UserDefinedInformation: "Some important info",
			ValidationType:         "card",
			ReferenceMask:          "############",
			AcceptanceQualifier:    "same_day",
		},
	}

	out, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	bodyDigest := utils.Hash(crypto.SHA256, []byte(strings.Join(flag.Args(), string(out))))
	base64BodyDigest := base64.StdEncoding.EncodeToString([]byte(bodyDigest))

	signatureDate := time.Now().Format(time.RFC1123)

	outLength := len(string(out))
	contentLength := strconv.FormatInt(int64(outLength), 10)

	signature := "(request-target): post " + path
	signature += " host: " + host
	signature += " date: " + signatureDate
	signature += " accept: application/vnd.api+json content-type: application/vnd.api+json "
	signature += " content-length: " + contentLength
	signature += " digest: SHA-256=" + base64BodyDigest

	headers := "(request-target) host date content-type accept digest content-length"
	signedSignature := utils.Sign(privateKeyFilePath, signature, crypto.SHA256)
	base64SignedSignature := base64.StdEncoding.EncodeToString([]byte(signedSignature))

	authorizationHeader := "Signature keyId=" + utils.GetPublicKeyContent(publicKeyFilePath)
	authorizationHeader += ",algorithm=\"rsa-sha256\""
	authorizationHeader += ",headers=" + headers
	authorizationHeader += ",signature=" + base64SignedSignature

	client := &http.Client{}
	reqbody, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", fullPath, bytes.NewBuffer(reqbody))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Host", host)
	req.Header.Add("Date", time.Now().Format(time.RFC1123))
	req.Header.Add("Accept", "application/vnd.api+json")
	req.Header.Add("Content-Type", "application/vnd.api+json")
	req.Header.Add("Content-Length", contentLength)
	req.Header.Add("Authorization", authorizationHeader)
	req.Header.Add("Digest", base64BodyDigest)

	fmt.Println("\nRequest header: ", req.Header)
	fmt.Println("\nRequest body: ", req.Body)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("\nResponse body: ", string(body))

}
