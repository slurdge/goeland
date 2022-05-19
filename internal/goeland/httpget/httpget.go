package httpget

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/slurdge/goeland/version"
)

var UserAgent string = "multiple:goeland:" + version.Version + " (commit id:" + version.GitCommit + ") (by /u/goelandrss)"
var defaultClient http.Client

func GetHTTPRessource(url string) (body []byte, err error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", UserAgent)
	request.Header.Set("Accept", "*/*")

	resp, err := defaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return body, fmt.Errorf("received error code %d", resp.StatusCode)
	}
	return body, err
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
}
