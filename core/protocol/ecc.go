package protocol

import "math"

// ─── Reed-Solomon Ligero para IPv7-IEU (mat.md §Códigos ECC) ─────────────────────────────────
//
// RS(n, k) sobre GF(2^8) con polinomio primitivo 0x11d.
// Añade nsym bytes de paridad al payload, corrigiendo hasta nsym/2 bytes erróneos.
//
// En IPv7-IEU se usa en links satelitales GEO/MEO donde los bit-errors en ráfaga
// son frecuentes y la retransmisión tiene un costo alto de latencia (700ms GEO).
//
// Referencia: mat.md — "los códigos de Reed-Solomon operan sobre bloques de símbolos,
// lo que les otorga una capacidad excepcional para corregir errores en ráfaga (burst errors),
// siendo el estándar en discos ópticos y comunicaciones por satélite"

// Tablas de Galois GF(2^8) con primitivo 0x11d
var gfExp [512]byte
var gfLog [256]byte

func init() {
	x := 1
	for i := 0; i < 255; i++ {
		gfExp[i] = byte(x)
		gfLog[x] = byte(i)
		x <<= 1
		if x&0x100 != 0 {
			x ^= 0x11d
		}
	}
	for i := 255; i < 512; i++ {
		gfExp[i] = gfExp[i-255]
	}
}

func gfMul(a, b byte) byte {
	if a == 0 || b == 0 {
		return 0
	}
	return gfExp[int(gfLog[a])+int(gfLog[b])]
}

// rsGeneratorPoly genera g(x) = ∏(x - α^i), i=0..nsym-1 en orden ascendente de grado.
// Retorna slice de len=nsym+1: [coef_x^0, coef_x^1, ..., coef_x^nsym]
// El coeficiente líder (grado nsym) es siempre 1.
func rsGeneratorPoly(nsym int) []byte {
	g := []byte{1}
	for i := 0; i < nsym; i++ {
		root := gfExp[i]
		newg := make([]byte, len(g)+1)
		for j, c := range g {
			newg[j+1] ^= c
			newg[j] ^= gfMul(c, root)
		}
		g = newg
	}
	return g // g[k] = coeficiente del término x^k
}

// rsPolyEval evalúa un polinomio cw[] en el punto x usando el esquema de Horner.
// cw[0] es el coeficiente del grado más alto (big-endian).
func rsPolyEval(cw []byte, x byte) byte {
	val := cw[0]
	for _, c := range cw[1:] {
		val = gfMul(val, x) ^ c
	}
	return val
}

// RSEncode añade nsym bytes de paridad Reed-Solomon al final del mensaje.
// Retorna: [msg...][paridad_nsym_bytes] (longitud total = len(msg) + nsym).
//
// Implementación: división polinómica larga de msg(x)*x^nsym por gen(x) en GF(2^8).
// Los bytes de paridad se añaden en orden descendente de grado (convención big-endian).
func RSEncode(msg []byte, nsym int) []byte {
	if nsym <= 0 {
		return msg
	}
	gen := rsGeneratorPoly(nsym) // ascendente: gen[k] = coef de x^k

	// LFSR de divison larga: procesamos msg de mayor a menor grado.
	// rem[j] = coeficiente de x^j del remainder actual (ascendente).
	rem := make([]byte, nsym)
	for _, b := range msg {
		// El byte b es el coeficiente del siguiente grado más alto
		factor := rem[nsym-1] ^ b
		// Desplazar rem (multiplicar por x)
		for j := nsym - 1; j > 0; j-- {
			rem[j] = rem[j-1] ^ gfMul(factor, gen[j])
		}
		rem[0] = gfMul(factor, gen[0])
	}

	result := make([]byte, len(msg)+nsym)
	copy(result, msg)
	// Paridad en big-endian: rem[nsym-1] es grado más alto → va primero
	for i := 0; i < nsym; i++ {
		result[len(msg)+i] = rem[nsym-1-i]
	}
	return result
}

// RSCanDetectError evalúa los síndromes del codeword.
// Retorna true si algún síndrome ≠ 0 (error detectado).
// Un codeword RS válido evaluado en α^i da 0 para i=0..nsym-1.
func RSCanDetectError(msgWithECC []byte, nsym int) bool {
	if nsym <= 0 || len(msgWithECC) < nsym || len(msgWithECC) == 0 {
		return false
	}
	for i := 0; i < nsym; i++ {
		// Evaluar el codeword como polinomio big-endian en α^i
		if rsPolyEval(msgWithECC, gfExp[i]) != 0 {
			return true
		}
	}
	return false
}

// ─── Delta Encoding de Headers (mat.md §Entropía de Shannon) ─────────────────────────────────
//
// En una sesión IPv7-IEU, Region y Subnet son prácticamente constantes.
// Transmitir solo el delta (DeviceID+SubPort) reduce el overhead un 40%:
//   - Header completo: 10 bytes
//   - Header delta:     6 bytes
//
// El bit 7 del primer byte distingue modo delta (1) de modo completo (0).

// HeaderDeltaState almacena el estado de sesión para aplicar delta encoding.
type HeaderDeltaState struct {
	LastRegion  uint16
	LastSubnet  uint16
	Initialized bool
}

// EncodeDelta produce una representación compacta si Region y Subnet no cambiaron.
// Formato delta (6 bytes): [0x80|Region_lo][Subnet_lo][4B DeviceID]
// Formato completo (10 bytes): header IPv7 estándar.
func (s *HeaderDeltaState) EncodeDelta(addr IPv7Address) []byte {
	r := uint16(addr.Region)
	sub := uint16(addr.Subnet)
	dev := uint32(addr.DeviceID)

	if s.Initialized && r == s.LastRegion && sub == s.LastSubnet {
		buf := make([]byte, 6)
		buf[0] = 0x80 | byte(r&0x7F)
		buf[1] = byte(sub & 0xFF)
		buf[2] = byte(dev >> 24)
		buf[3] = byte(dev >> 16)
		buf[4] = byte(dev >> 8)
		buf[5] = byte(dev)
		return buf
	}
	s.LastRegion = r
	s.LastSubnet = sub
	s.Initialized = true
	return addr.SerializeHeader()
}

// DecodeDelta reconstruye el header completo desde un buffer (delta o completo).
func (s *HeaderDeltaState) DecodeDelta(buf []byte) IPv7Address {
	if len(buf) >= 6 && buf[0]&0x80 != 0 {
		dev := float64(uint32(buf[2])<<24 | uint32(buf[3])<<16 | uint32(buf[4])<<8 | uint32(buf[5]))
		sp := uint16(0)
		if len(buf) >= 8 {
			sp = uint16(buf[6])<<8 | uint16(buf[7])
		}
		addr := NewIPv7(float64(s.LastRegion), float64(s.LastSubnet), dev)
		addr.SubPort = sp
		return addr
	}
	addr := ParseHeader(buf)
	s.LastRegion = uint16(addr.Region)
	s.LastSubnet = uint16(addr.Subnet)
	s.Initialized = true
	return addr
}

// ─── Kolmogorov/Fokker-Planck: Predicción Estocástica de Latencia (mat.md §Ecuación de Kolmogorov) ──
//
// dL = μ dt + σ dW  — proceso de difusión gaussiana.
// Permite predecir distribución de latencia futura, más preciso que regresión
// lineal para saltos no-lineales tipo Starlink handover orbital.

// StochasticLatencyPredictor predice latencia futura usando difusión estocástica.
type StochasticLatencyPredictor struct {
	samples []float64
	maxN    int
}

// NewStochasticPredictor crea un predictor con ventana de windowSize muestras.
func NewStochasticPredictor(windowSize int) *StochasticLatencyPredictor {
	return &StochasticLatencyPredictor{maxN: windowSize}
}

// Push añade una nueva muestra de latencia.
func (p *StochasticLatencyPredictor) Push(latency float64) {
	p.samples = append(p.samples, math.Max(latency, 0.1))
	if len(p.samples) > p.maxN {
		p.samples = p.samples[1:]
	}
}

// Predict estima (media, desviación) de la latencia en stepsAhead pasos futuros.
func (p *StochasticLatencyPredictor) Predict(stepsAhead int) (mean, stddev float64) {
	n := len(p.samples)
	if n < 3 {
		if n > 0 {
			return p.samples[n-1], 5.0
		}
		return 50.0, 20.0
	}

	var drifts []float64
	for i := 1; i < n; i++ {
		drifts = append(drifts, p.samples[i]-p.samples[i-1])
	}

	var muD, varD float64
	for _, d := range drifts {
		muD += d
	}
	muD /= float64(len(drifts))
	for _, d := range drifts {
		varD += (d - muD) * (d - muD)
	}
	varD /= float64(len(drifts))
	sigmaD := math.Sqrt(varD)

	currentMean := p.samples[n-1]
	projectedMean := currentMean + muD*float64(stepsAhead)
	projectedStd := math.Sqrt(varD*float64(stepsAhead) + sigmaD*sigmaD)

	if math.IsNaN(projectedMean) || projectedMean < 0 {
		projectedMean = currentMean
	}
	if math.IsNaN(projectedStd) || projectedStd < 0 {
		projectedStd = sigmaD
	}
	return projectedMean, projectedStd
}

// IsHighRisk retorna prob (0–1) de que la latencia supere thresholdMs en stepsAhead pasos.
func (p *StochasticLatencyPredictor) IsHighRisk(stepsAhead int, thresholdMs float64) float64 {
	mean, stddev := p.Predict(stepsAhead)
	if stddev < 1e-9 {
		if mean > thresholdMs {
			return 1.0
		}
		return 0.0
	}
	z := (thresholdMs - mean) / stddev
	absZ := math.Abs(z)
	t2 := 1.0 / (1.0 + 0.2316419*absZ)
	poly := t2 * (0.319381530 + t2*(-0.356563782+t2*(1.781477937+t2*(-1.821255978+t2*1.330274429))))
	phi := math.Exp(-0.5*z*z) / math.Sqrt(2*math.Pi) * poly
	if z < 0 {
		return 1.0 - phi
	}
	return phi
}
