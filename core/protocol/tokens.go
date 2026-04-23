package protocol

// TablaMaestraObjetos maps standard concepts to single-byte tokens for IEU compression
var TablaMaestraObjetos = map[string]byte{
	// Web Standards
	"text/html":        0x01,
	"application/json": 0x02,
	"text/plain":       0x03,
	"image/jpeg":       0x04,
	"image/png":        0x05,

	// HTTP Methods
	"GET":    0x10,
	"POST":   0x11,
	"PUT":    0x12,
	"DELETE": 0x13,

	// Status Codes
	"200": 0x20, // OK
	"404": 0x21, // Not Found
	"500": 0x22, // Internal Server Error

	// Accounting Standards (SII Chile)
	"balance_general":   0x30,
	"estado_resultados": 0x31,
	"flujo_caja":        0x32,
	"formulario_f29":    0x33,
	"impuesto_renta":    0x34,

	// Conceptual Tokens
	"identidad":  0x40,
	"proposito":  0x41,
	"estructura": 0x42,
	"estado":     0x43,
	"flujo":      0x44,
	"integridad": 0x45,
}

// ComprimirConcepto converts a standard string to its IEU token
func ComprimirConcepto(concepto string) (byte, bool) {
	token, exists := TablaMaestraObjetos[concepto]
	return token, exists
}

// ExpandirToken converts a token back to its concept
func ExpandirToken(token byte) (string, bool) {
	for concepto, t := range TablaMaestraObjetos {
		if t == token {
			return concepto, true
		}
	}
	return "", false
}
