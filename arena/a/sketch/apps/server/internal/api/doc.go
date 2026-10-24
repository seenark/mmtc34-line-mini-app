// Package api holds the JSON shapes and enum vocabularies shared with the web app.
// It imports nothing from this module, so every other package may import it.
// apps/web/src/api/types.gen.ts is generated from it by tygo and is git-ignored.
//
// Money is whole Thai baht. Times are RFC3339 strings.
package api

type Config struct {
	LiffID   string `json:"liffId"`
	DevLogin bool   `json:"devLogin"`
}

// ErrorResponse is the body of every non-2xx JSON response. Error is Thai, shown verbatim.
type ErrorResponse struct {
	Error string `json:"error"`
}

type ImageUpload struct {
	Name string `json:"name"`
}
