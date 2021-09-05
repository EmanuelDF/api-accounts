package controllers

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/emanueldf/form3-accounts/models"
	"github.com/emanueldf/form3-accounts/utils"
)

func Init() {

	path := "/v1/organisation/accounts"
	host := "api.staging-form3.tech"
	base_url := "https://" + host
	full_path := base_url + path

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

	digest := base64.StdEncoding.EncodeToString([]byte(utils.Sign(utils.ReadPrivateKeyContent(), crypto.SHA256)))

	headers := "(request-target) host date content-type accept digest content-length"

	authorization_header := "Signature: keyId=" + utils.ReadPublicKey() + ",algorithm=\"rsa-sha256\",headers=" + headers + ",signature=" + digest

	client := &http.Client{}
	reqbody, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", full_path, bytes.NewBuffer(reqbody))
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Host", "api.form3.tech")
	req.Header.Add("Date", time.Now().Format(time.RFC1123))
	//req.Header.Add("Accept", "application/vnd.api+json")
	req.Header.Add("Digest", "SHA-256=WllU95a/P37KDBmTedpEIIvVtBgRqDdYrHz06NXDuvk=")
	req.Header.Add("Content-Length", strconv.FormatInt(int64(req.ContentLength), 10))
	req.Header.Add("Content-Type", "application/vnd.api+json")
	req.Header.Add("Authorization", authorization_header)

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
