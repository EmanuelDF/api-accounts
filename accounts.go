package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Account struct
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

	url := "https://api.staging-form3.tech/v1/organisation/accounts"
	method := "POST"

	payload := strings.NewReader(`{
		"data": {
			"id": "{{random_guid}}",
			"organisation_id": "",
			"type": "accounts",
			"attributes": {
			"country": "GB",
				"base_currency": "GBP",
				"bank_id": "400302",
				"bank_id_code": "GBDSC",
				"account_number": "10000004",
				"customer_id": "234",
				"iban": "GB28NWBK40030212764204",
				"bic": "NWBKGB42",
				"account_classification": "Personal"
			}
		}
	}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "{{authorization}}")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Date", "{{request_date}}")
	req.Header.Add("Digest", "{{request_signing_digest}}")

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
	req.Header.Add("Authorization", "{{authorization}}")
	req.Header.Add("Accept", "application/vnd.api+json")
	req.Header.Add("Date", "{{request_date}}")

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
	req.Header.Add("Authorization", "{{authorization}}")
	req.Header.Add("Date", "{{request_date}}")

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
