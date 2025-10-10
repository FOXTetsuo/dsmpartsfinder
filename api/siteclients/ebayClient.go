package siteclients

import (
	"net/http"
)

type EbayClient struct {
	baseURL    string
	httpClient *http.Client
	siteID     int
}
