package httpget

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/slurdge/goeland/version"
)

var userAgent string = "multiple:goeland:" + version.Version + " (commit id:" + version.GitCommit + ") (by /u/goelandrss)"
var defaultClient http.Client
var insecureClient http.Client

func GetHTTPRessourceGeneric(url string, client http.Client) (body []byte, err error) {
	var request *http.Request
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", userAgent)
	request.Header.Set("Accept", "*/*")

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return body, fmt.Errorf("received error code %d", resp.StatusCode)
	}
	return body, err
}

// GetHTTPRessource gets the bytes corresponding to an URL
func GetHTTPRessource(url string) (body []byte, err error) {
	return GetHTTPRessourceGeneric(url, defaultClient)
}

// GetHTTPRessource gets the bytes corresponding to an URL, without checking
func GetHTTPRessourceInsecure(url string) (body []byte, err error) {
	return GetHTTPRessourceGeneric(url, insecureClient)
}

func init() {
	//this one is needed because of incompatibility between latest golang and reddit
	defaultClient = http.Client{
		Transport: &http.Transport{
			//TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
			TLSClientConfig: &tls.Config{
				//KeyLogWriter: dmp,
			},
			//ForceAttemptHTTP2: false,
		},
	}
	insecureClient = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}
