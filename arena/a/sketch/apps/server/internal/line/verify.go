// Package line verifies LIFF ID tokens with LINE. It satisfies server.Verifier.
package line

import (
	"context"
	"net/http"
)

type Client struct {
	ChannelID string
	HTTP      *http.Client // 10s timeout
}

// Verify POSTs form {id_token, client_id} to https://api.line.me/oauth2/v2.1/verify.
// Non-200 -> error. Returns the response's `sub` (the LINE userId). LINE checks aud/exp/iss.
func (c Client) Verify(ctx context.Context, idToken string) (lineUserID string, err error) {
	panic("not implemented")
}
