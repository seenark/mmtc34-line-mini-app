package server_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

// send is a request with raw headers, which the Client helpers cannot set.
func (h *H) send(method, path, secFetchSite string, body string) *http.Response {
	h.t.Helper()
	req, err := http.NewRequest(method, h.srv.URL+h.base+path, strings.NewReader(body))
	if err != nil {
		h.t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Sec-Fetch-Site", secFetchSite)
	resp, err := h.Anonymous().http.Do(req)
	if err != nil {
		h.t.Fatal(err)
	}
	resp.Body.Close()
	return resp
}

func TestOversizedJSONBodyIs413(t *testing.T) {
	h := New(t)

	res := h.Anonymous().POST("/api/auth/dev", api.DevLoginInput{LineUserID: strings.Repeat("a", 100<<10)}).Status(413)

	Equal(t, res.Error(), "ข้อมูลใหญ่เกินกำหนด")
}

func TestCrossSiteWritesAreRejected(t *testing.T) {
	h := New(t)
	const login = `{"lineUserId":"U-nok"}`

	Equal(t, h.send(http.MethodPost, "/api/auth/dev", "same-site", login).StatusCode, 403)
	Equal(t, h.send(http.MethodPost, "/api/auth/dev", "cross-site", login).StatusCode, 403)
	Equal(t, h.send(http.MethodPost, "/api/auth/dev", "same-origin", login).StatusCode, 204)
	Equal(t, h.send(http.MethodGet, "/api/config", "cross-site", "").StatusCode, 200)
}

func TestWriteWithoutSecFetchSiteIsAllowed(t *testing.T) {
	h := New(t)

	h.Anonymous().POST("/api/auth/dev", api.DevLoginInput{LineUserID: "U-nok"}).Status(204)
}
