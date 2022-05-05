package httpget

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"

	"github.com/slurdge/goeland/version"
)

var UserAgent string = "multiple:goeland:" + version.Version + " (commit id:" + version.GitCommit + ") (by /u/goelandrss)"

func GetHTTPRessource(url string) (body []byte, err error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	//this one is needed because of incompatibility between latest golang and reddit
	request.Header.Set("User-Agent", UserAgent)
	request.Header.Set("Accept", "*/*")
	var defaultClient = http.Client{
		Transport: &http.Transport{
			TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
		},
	}
	resp, err := defaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
