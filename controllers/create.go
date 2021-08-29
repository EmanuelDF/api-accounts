package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/emanueldf/form3-accounts/models"
	"github.com/emanueldf/form3-accounts/utils"
)

func Create() {

	var (
		accountClassification = "Personal"
		country               = "GB"
	)

	uuid := utils.Generate()

	var payload = models.AccountData{
		Attributes: &models.AccountAttributes{
			AccountClassification:   &accountClassification,
			AccountMatchingOptOut:   new(bool),
			AccountNumber:           "10000004",
			AlternativeNames:        []string{},
			BankID:                  "400302",
			BankIDCode:              "GBDSC",
			BaseCurrency:            "GBP",
			Bic:                     "NWBKGB42",
			Country:                 &country,
			Iban:                    "GB28NWBK40030212764204",
			JointAccount:            new(bool),
			Name:                    []string{},
			SecondaryIdentification: "234",
			Status:                  new(string),
			Switched:                new(bool),
		},
		ID:             "",
		OrganisationID: uuid,
		Type:           "accounts",
		Version:        new(int64),
	}

	client := &http.Client{}
	reqbody, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	url := "https://api.staging-form3.tech/v1/organisation/accounts"
	method := "POST"

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqbody))
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
	req.Header.Add("Authorization", "Signature keyId=\"75a8ba12-fff2-4a52-ad8a-e8b34c5ccec8\","+
		"algorithm=\"rsa-sha256\","+
		"headers=\"(request-target) host date digest content-length content-type\","+
		"signature=\"dOO1gnywk/Awo2Z0vSxcxrcoPZ51wKbMG8JYIBJ+xn4MUVjDy/ooP7l7EzsDQVZPj8ylkzLMvYoDnzyKA1xNaphoujRpfs1wBoqe4DCMFMeVaDZLsGXbpgijYICCdriYoLo0agjbpeDh+zeyh2b/+wGLevy7oHB+KrBtBHtYDXo=\"")

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
