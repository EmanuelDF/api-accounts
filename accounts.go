package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Account represents an account in the form3 org section.
// See https://api-docs.form3.tech/api.html#organisation-accounts for
// more information about fields.
type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`
}

type AccountAttributes struct {
	AccountClassification   *string  `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  *string  `json:"status,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
}

func main() {
	create()
	fetch()
	delete()
}

func create() {

	var accountClassification string = "Personal"
	var country string = "GB"

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	payload := AccountData{
		Attributes: &AccountAttributes{
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

	var contentLength int = len(reqbody)

	req.Header.Add("Host", "api.form3.tech")
	req.Header.Add("Date", time.Now().Format(time.RFC1123))
	req.Header.Add("Accept", "application/vnd.api+json")
	req.Header.Add("Content-Type", "application/vnd.api+json")
	req.Header.Add("Content-Length", strconv.FormatInt(int64(contentLength), 10))
	req.Header.Add("Digest", "SHA-256=WllU95a/P37KDBmTedpEIIvVtBgRqDdYrHz06NXDuvk=")
	req.Header.Add("Authorization",
		"Signature keyId=\"75a8ba12-fff2-4a52-ad8a-e8b34c5ccec8\",algorithm=\"rsa-sha256\",headers=\"(request-target) host date accept digest content-length content-type\",signature=\"sEl9KI0sK1NTxFYpVa+u8NBxnQx12zDEHSo/ijfvqi9z8zt5O1aXjoy8fyLvg/ICXaHoogb9oJ4C4i1iJDP1RCiTpW0OvwNPP4t0XlGnKlKX4iyLV4CofR8H9o/X5mcsiv/tVP7qCgP92efaisLCVjE9MKMPjDaA7Tj3gBbeYnI=\"")

	fmt.Println("Header: ", req.Header)
	fmt.Println("Body: ", req.Body)

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
	fmt.Println(string(body))
}

func fetch() {

	url := "https://api.staging-form3.tech/v1/organisation/accounts/{{account_id}}"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization",
		"Signature keyId=\"75a8ba12-fff2-4a52-ad8a-e8b34c5ccec8\",algorithm=\"rsa-sha256\",headers=\"(request-target) host date content-type accept digest content-length\",signature=\"sEl9KI0sK1NTxFYpVa+u8NBxnQx12zDEHSo/ijfvqi9z8zt5O1aXjoy8fyLvg/ICXaHoogb9oJ4C4i1iJDP1RCiTpW0OvwNPP4t0XlGnKlKX4iyLV4CofR8H9o/X5mcsiv/tVP7qCgP92efaisLCVjE9MKMPjDaA7Tj3gBbeYnI=\"")
	req.Header.Add("Accept", "application/vnd.api+json")
	req.Header.Add("Date", time.Now().Format(time.RFC1123))

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
	fmt.Println(string(body))
}

func delete() {

	url := "https://api.staging-form3.tech/v1/organisation/accounts/{{account_id}}?version=0"
	method := "DELETE"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization",
		"Signature keyId=\"75a8ba12-fff2-4a52-ad8a-e8b34c5ccec8\",algorithm=\"rsa-sha256\",headers=\"(request-target) host date content-type accept digest content-length\",signature=\"sEl9KI0sK1NTxFYpVa+u8NBxnQx12zDEHSo/ijfvqi9z8zt5O1aXjoy8fyLvg/ICXaHoogb9oJ4C4i1iJDP1RCiTpW0OvwNPP4t0XlGnKlKX4iyLV4CofR8H9o/X5mcsiv/tVP7qCgP92efaisLCVjE9MKMPjDaA7Tj3gBbeYnI=\"")
	req.Header.Add("Date", time.Now().Format(time.RFC1123))

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
	fmt.Println(string(body))
}
