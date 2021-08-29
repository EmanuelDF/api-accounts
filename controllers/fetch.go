package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func Fetch() {

	url := "https://api.staging-form3.tech/v1/organisation/accounts/{{account_id}}"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Host", "api.form3.tech")
	req.Header.Add("Date", time.Now().Format(time.RFC1123))
	//req.Header.Add("Accept", "application/vnd.api+json")
	req.Header.Add("Authorization",
		"Signature keyId=\"75a8ba12-fff2-4a52-ad8a-e8b34c5ccec8\","+
			"algorithm=\"rsa-sha256\","+
			"headers=\"(request-target) host date\","+
			"signature=\"sEl9KI0sK1NTxFYpVa+u8NBxnQx12zDEHSo/ijfvqi9z8zt5O1aXjoy8fyLvg/ICXaHoogb9oJ4C4i1iJDP1RCiTpW0OvwNPP4t0XlGnKlKX4iyLV4CofR8H9o/X5mcsiv/tVP7qCgP92efaisLCVjE9MKMPjDaA7Tj3gBbeYnI=\"")

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
	fmt.Println("Fetch response: ", string(body))
}
