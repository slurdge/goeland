package fetch

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/internal/goeland/i18n"
	"github.com/slurdge/goeland/log"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	i18n.Init("en-US")
	log.SetDefaultLogger(log.NewLogger(viper.New()))
	os.Exit(m.Run())
}

type minifluxRoundTripper func(*http.Request) (*http.Response, error)

func (f minifluxRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type minifluxResponseBody struct {
	io.ReadCloser
	closes int
}

func (b *minifluxResponseBody) Close() error {
	b.closes++
	return b.ReadCloser.Close()
}

func TestFetchMinifluxClosesMarkAsReadResponse(t *testing.T) {
	for _, status := range []int{http.StatusNoContent, http.StatusInternalServerError} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.URL.Path != "/v1/entries" {
					t.Errorf("request path = %q, want /v1/entries", req.URL.Path)
				}
				if req.Header.Get("X-Auth-Token") != "test-token" {
					t.Error("request is missing the configured token")
				}
				switch req.Method {
				case http.MethodGet:
					w.Header().Set("Content-Type", "application/json")
					io.WriteString(w, `{"total":1,"entries":[{"id":42,"hash":"entry-hash","title":"Entry","url":"https://example.com/entry","content":"<p>Body</p>"}]}`)
				case http.MethodPut:
					var payload minifluxEntries
					if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
						t.Errorf("decode mark-as-read payload: %v", err)
					}
					if !reflect.DeepEqual(payload.IDs, []int64{42}) || payload.Status != "read" {
						t.Errorf("mark-as-read payload = %+v", payload)
					}
					w.WriteHeader(status)
					if status != http.StatusNoContent {
						io.WriteString(w, "could not mark entries read")
					}
				default:
					t.Errorf("unexpected request method %s", req.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			}))
			t.Cleanup(server.Close)

			originalTransport := http.DefaultTransport
			transport := server.Client().Transport.(*http.Transport).Clone()
			bodies := make(map[string]*minifluxResponseBody)
			http.DefaultTransport = minifluxRoundTripper(func(req *http.Request) (*http.Response, error) {
				res, err := transport.RoundTrip(req)
				if err == nil {
					body := &minifluxResponseBody{ReadCloser: res.Body}
					bodies[req.Method] = body
					res.Body = body
				}
				return res, err
			})
			t.Cleanup(func() {
				http.DefaultTransport = originalTransport
				for _, body := range bodies {
					if body.closes == 0 {
						body.ReadCloser.Close()
					}
				}
				transport.CloseIdleConnections()
			})

			source := &goeland.Source{Name: "test", Title: "Test"}
			err := fetchMiniflux(source, server.URL+"/v1/entries?status=unread", "test-token", false, true, false)
			if status == http.StatusNoContent {
				if err != nil {
					t.Fatalf("fetchMiniflux: %v", err)
				}
			} else if err == nil || !strings.Contains(err.Error(), "non-204 status 500") {
				t.Errorf("fetchMiniflux error = %v, want non-204 status 500", err)
			}
			if len(source.Entries) != 1 || source.Entries[0].UID != "miniflux-entry-hash" {
				t.Errorf("parsed entries = %+v", source.Entries)
			}
			for _, method := range []string{http.MethodGet, http.MethodPut} {
				body, ok := bodies[method]
				if !ok {
					t.Errorf("missing %s response", method)
				} else if body.closes != 1 {
					t.Errorf("%s response body closed %d times, want 1", method, body.closes)
				}
			}
		})
	}
}
