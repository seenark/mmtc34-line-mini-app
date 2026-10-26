package api

type Config struct {
	LiffID     string `json:"liffId"`
	DevLogin   bool   `json:"devLogin"`
	AdminLogin bool   `json:"adminLogin"` // Super Admin username + password login is enabled
	Banks      []Bank `json:"banks"`      // the registry a Campaign's bankCode is picked from
}

// Bank is a Thai bank: its 3-digit Bank of Thailand code and Thai display name.
type Bank struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// ErrorResponse is the body of every non-2xx JSON response. Error is Thai, shown verbatim.
type ErrorResponse struct {
	Error string `json:"error"`
}

type ImageUpload struct {
	Name string `json:"name"`
}
