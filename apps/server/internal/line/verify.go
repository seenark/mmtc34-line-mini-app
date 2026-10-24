// Package line verifies LIFF ID tokens with LINE. It satisfies server.Verifier.
package line

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const verifyURL = "https://api.line.me/oauth2/v2.1/verify"

type Client struct {
	ChannelID string
	HTTP      *http.Client // nil -> a client with a 10s timeout
}

var defaultHTTP = &http.Client{Timeout: 10 * time.Second}

// Verify POSTs form {id_token, client_id} to LINE and returns the response's `sub` (the LINE
// userId). LINE checks signature, aud, iss and exp; any non-2xx or missing sub is an error.
func (c Client) Verify(ctx context.Context, idToken string) (string, error) {
	hc := c.HTTP
	if hc == nil {
		hc = defaultHTTP
	}
	form := url.Values{"id_token": {idToken}, "client_id": {c.ChannelID}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, verifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("line verify: status %d", resp.StatusCode)
	}
	var body struct {
		Sub string `json:"sub"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("line verify: %w", err)
	}
	if body.Sub == "" {
		return "", errors.New("line verify: response has no sub")
	}
	return body.Sub, nil
}
