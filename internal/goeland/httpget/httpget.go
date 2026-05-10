package httpget

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/slurdge/goeland/version"
	"github.com/spf13/viper"
)

var userAgent string = "multiple:goeland:" + version.Version + " (commit id:" + version.GitCommit + ") (by /u/goelandrss)"
var defaultClient http.Client
var defaultDialer net.Dialer
var insecureClient http.Client

var errUnsafeURL = errors.New("unsafe url")

func GetHTTPRessourceGeneric(url string, client http.Client) (body []byte, err error) {
	var request *http.Request
	request, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if err = validateURL(request.URL); err != nil {
		return nil, err
	}
	config := viper.GetViper()
	if config.GetString("user-agent") != "" {
		userAgent = config.GetString("user-agent")
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

func validateURL(parsedURL *url.URL) error {
	if parsedURL == nil {
		return fmt.Errorf("%w: empty url", errUnsafeURL)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("%w: unsupported scheme %q", errUnsafeURL, parsedURL.Scheme)
	}
	if parsedURL.Hostname() == "" {
		return fmt.Errorf("%w: missing host", errUnsafeURL)
	}
	return nil
}

func isSafeIP(ip net.IP) bool {
	return ip.IsGlobalUnicast() && !ip.IsPrivate()
}

func filterDialContext(ctx context.Context, network string, address string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	resolver := net.DefaultResolver
	addresses, err := resolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}

	for _, resolvedAddress := range addresses {
		if !isSafeIP(resolvedAddress.IP) {
			return nil, fmt.Errorf("%w: blocked address %s for host %s", errUnsafeURL, resolvedAddress.IP.String(), host)
		}
	}
	return defaultDialer.DialContext(ctx, network, address)
}

func init() {
	originalCheckRedirect := http.DefaultClient.CheckRedirect
	defaultDialer = net.Dialer{}
	checkRedirect := func(request *http.Request, via []*http.Request) error {
		if err := validateURL(request.URL); err != nil {
			return err
		}
		if originalCheckRedirect != nil {
			return originalCheckRedirect(request, via)
		}
		return nil
	}

	//this one is needed because of incompatibility between latest golang and reddit
	defaultClient = http.Client{
		Transport: &http.Transport{
			DialContext:     filterDialContext,
			TLSClientConfig: &tls.Config{},
		},
		CheckRedirect: checkRedirect,
	}
	insecureClient = http.Client{
		Transport: &http.Transport{
			DialContext: filterDialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		CheckRedirect: checkRedirect,
	}
}
