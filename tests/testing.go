package main

import (
	"fmt"
	"net/http"

	rc "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

type transport struct {
	token               string
	underlyingTransport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	return t.underlyingTransport.RoundTrip(req)
}

func NewAccountAPIClient(port int) *account_api.Client {
	config := client.DefaultTransportConfig().
		WithHost(fmt.Sprintf("%s:%d", "localhost", port)).
		WithSchemes([]string{"http"})

	config.WithBasePath("/v1")

	transport := &transport{
		underlyingTransport: http.DefaultTransport,
	}
	h := &http.Client{Transport: transport}

	rt := rc.NewWithClient(config.Host, config.BasePath, config.Schemes, h)
	rt.SetDebug(true)
	return account_api.New(rt, strfmt.Default)
}
