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
