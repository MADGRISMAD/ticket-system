package utils

// StandardResponse define la estructura estándar de respuestas en JSON
type StandardResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"` // Si hay error, se agrega; si no, se omite
}
