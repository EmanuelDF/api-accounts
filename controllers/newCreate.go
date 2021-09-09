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
	"strings"
	"time"

	"github.com/emanueldf/form3-accounts/models"
	"github.com/emanueldf/form3-accounts/utils"
)

func Init() {

	privateKeyFilePath := "/Users/emanuel/go/src/github.com/emanueldf/form3-accounts/certs/test_private_key.pem"
	publicKeyFilePath := "/Users/emanuel/go/src/github.com/emanueldf/form3-accounts/certs/test_public_key.pem"

	path := "/v1/organisation/accounts"
	host := "api.staging-form3.tech"
	baseUrl := "https://" + host
	fullPath := baseUrl + path

	var Country = "GB"

	var payload = models.NewAccountData{
		Type:           "accounts",
		ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
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

	signature := "(request-target): post" + path
	signature += "host:" + host
	signature += "date:" + signatureDate
	signature += " content-type: application/json accept: application/json digest: SHA-256=" + base64BodyDigest
	signature += " content-length:" + string(rune(len(string(out))))

	headers := "(request-target) host date content-type accept digest content-length"

	signedSignature := utils.Sign(privateKeyFilePath, signature, crypto.SHA256)

	base64SignedSignature := base64.StdEncoding.EncodeToString([]byte(signedSignature))

	authorizationHeader := "Signature: keyId=" + utils.GetPublicKeyContent(publicKeyFilePath) +
		",algorithm=\"rsa-sha256\",headers=" + headers + ",signature=" + base64SignedSignature

	base64Signature := base64.StdEncoding.EncodeToString([]byte(signature))

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

	req.Header.Add("Accept", "application/vnd.api+json")
	req.Header.Add("Content-Type", "application/vnd.api+json")
	req.Header.Add("Date", time.Now().Format(time.RFC1123))
	req.Header.Add("Authorization", authorizationHeader)
	req.Header.Add("Digest", base64BodyDigest)
	req.Header.Add("Signature-Debug", base64Signature)

	fmt.Println("\nHeader request: ", req.Header)
	fmt.Println("\nBody request: ", req.Body, "\n ")

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
	fmt.Println("Create response: ", string(body))

}
