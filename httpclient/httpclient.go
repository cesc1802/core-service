package httpclient

import (
	"github.com/cesc1802/core-service/config"
	"github.com/go-resty/resty/v2"
	"time"
)

func NewRestyClient(c config.HttpClientConfig, baseUrl string) (*resty.Client, error) {
	// Create a Resty Client
	client := resty.New()

	// Unique settings at Client level
	//--------------------------------
	// Enable debug mode
	client.SetDebug(c.Debug)

	// Set client timeout as per your need
	duration, err := time.ParseDuration(c.Timeout)
	if err != nil {
		return nil, err
	}
	client.SetTimeout(duration)

	// You can override all below settings and options at request level if you want to
	//--------------------------------------------------------------------------------
	// Host URL for all request. So you can use relative URL in the request
	client.SetHostURL(baseUrl)

	// Headers for all request
	client.SetHeader("Accept", "application/json")
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
	})

	return client, nil
}
