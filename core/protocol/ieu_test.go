package protocol

import (
	"testing"
)

func TestNewIPv7_Inyectividad(t *testing.T) {
	a := NewIPv7(56, 1, 42)
	b := NewIPv7(56, 1, 42)
	c := NewIPv7(56, 1, 43)

	if !a.Equals(b) {
		t.Error("Misma (r,s,d) debe dar misma ResolvedIP")
	}
	if a.Equals(c) {
		t.Error("Diferente DeviceID debe dar ResolvedIP distinta")
	}
	if a.ResolvedIP == 0 {
		t.Error("ResolvedIP no debe ser cero con valores positivos")
	}
}

func TestNextHop_GradienteLogaritmico(t *testing.T) {
	n := &Node{
		Name:    "test-node",
		Latency: 50.0,
	}

	n1 := &Node{Name: "n1", Latency: 30.0}
	n2 := &Node{Name: "n2", Latency: 80.0}
	n3 := &Node{Name: "n3", Latency: 10.0} // mejor latencia

	n.Neighbors = []*Node{n1, n2, n3}

	best := n.NextHop()
	if best != n3 {
		t.Errorf("Deberia elegir el vecino con menor latencia (n3). Obtuvo: %s", best.Name)
	}
}

func TestNextHop_ProteccionLogCero(t *testing.T) {
	n := &Node{Latency: 50.0}
	n.Neighbors = []*Node{
		{Name: "zero", Latency: 0.0},
	}

	best := n.NextHop()
	if best == nil {
		t.Error("NextHop no debe fallar con latencia cero")
	}
}

// PQC Tests - ML-DSA-65 (FIPS-204)

func TestPQC_GenerateSignature_NotNil(t *testing.T) {
	data := []byte("QUANTUM_HANDSHAKE_SYNC")
	sig := GenerateSignature(data)
	if sig == nil {
		t.Fatal("GenerateSignature devolvio nil: las claves PQC no estan inicializadas")
	}
	if len(sig) == 0 {
		t.Fatal("GenerateSignature devolvio firma vacia")
	}
}

func TestPQC_SignVerify_RoundTrip(t *testing.T) {
	data := []byte("paquete-de-prueba-ipv7-ieu")
	sig := GenerateSignature(data)
	if sig == nil {
		t.Fatal("No se pudo generar firma PQC")
	}

	pkBytes := GetPublicKey()
	if pkBytes == nil {
		t.Fatal("GetPublicKey devolvio nil")
	}

	if !VerifySignature(pkBytes, data, sig) {
		t.Error("VerifySignature fallo en un mensaje valido (round-trip roto)")
	}
}

func TestPQC_Tamper_Detection(t *testing.T) {
	data := []byte("mensaje-original")
	sig := GenerateSignature(data)
	pkBytes := GetPublicKey()

	tampered := []byte("mensaje-ALTERADO")
	if VerifySignature(pkBytes, tampered, sig) {
		t.Error("VerifySignature acepto mensaje alterado: PQC no detecta tampering")
	}
}

func TestPQC_Tamper_Signature(t *testing.T) {
	data := []byte("mensaje-valido")
	sig := GenerateSignature(data)
	pkBytes := GetPublicKey()

	if len(sig) > 10 {
		sig[0] ^= 0xFF
		sig[5] ^= 0xAB
	}
	if VerifySignature(pkBytes, data, sig) {
		t.Error("VerifySignature acepto firma corrupta: el sistema es vulnerable")
	}
}

func TestPQC_SerializeHeader_And_Back(t *testing.T) {
	addr := NewIPv7(56, 1, 100)
	header := addr.SerializeHeader()

	// Header ahora es 10 bytes: [2B Region][2B Subnet][4B DeviceID][2B SubPort]
	if len(header) != 10 {
		t.Errorf("SerializeHeader debe producir exactamente 10 bytes, obtuvo %d", len(header))
	}

	parsed := ParseHeader(header)
	if parsed.Region != addr.Region || parsed.Subnet != addr.Subnet || parsed.DeviceID != addr.DeviceID {
		t.Errorf("ParseHeader no reconstruyo correctamente: original=%v parsed=%v", addr, parsed)
	}
}

// ─── Sub-Port Tests ──────────────────────────────────────────────────────────

func TestSubPort_Default_Zero(t *testing.T) {
	addr := NewIPv7(56, 1, 100)
	if addr.SubPort != 0 {
		t.Errorf("SubPort por defecto debe ser 0, obtuvo %d", addr.SubPort)
	}
}

func TestSubPort_WithSubPort_Constructor(t *testing.T) {
	addr := NewIPv7WithSubPort(56, 1, 200, 8080)
	if addr.SubPort != 8080 {
		t.Errorf("SubPort debe ser 8080, obtuvo %d", addr.SubPort)
	}
	// La identidad base no debe cambiar
	base := NewIPv7(56, 1, 200)
	if !addr.Equals(base) {
		t.Error("NewIPv7WithSubPort debe producir la misma ResolvedIP que NewIPv7 (sub-puerto es canal logico, no identidad)")
	}
}

func TestSubPort_SerializeDeserialize(t *testing.T) {
	original := NewIPv7WithSubPort(56, 1, 300, 65535)
	header := original.SerializeHeader()

	if len(header) != 10 {
		t.Fatalf("Header con sub-puerto debe ser 10 bytes, obtuvo %d", len(header))
	}

	parsed := ParseHeader(header)
	if parsed.SubPort != 65535 {
		t.Errorf("SubPort 65535 no serializo/deserializo correctamente, obtuvo %d", parsed.SubPort)
	}
}

func TestSubPort_Max_Value(t *testing.T) {
	addr := NewIPv7WithSubPort(1, 1, 1, 65535)
	header := addr.SerializeHeader()
	parsed := ParseHeader(header)
	if parsed.SubPort != 65535 {
		t.Errorf("SubPort maximo (65535) fallo, obtuvo %d", parsed.SubPort)
	}
}

func TestSubPort_Zero_Is_CatchAll(t *testing.T) {
	// Sub-puerto 0 = catch-all (cualquier canal)
	addr := NewIPv7WithSubPort(10, 1, 99, 0)
	if addr.SubPort != 0 {
		t.Error("SubPort 0 debe ser el canal catch-all por defecto")
	}
}

func TestSubPort_IndependentChannels(t *testing.T) {
	// Dos nodos con misma identidad pero distintos sub-puertos son canales distintos
	addrREST := NewIPv7WithSubPort(56, 1, 100, 7780)
	addrCoAP := NewIPv7WithSubPort(56, 1, 100, 5683)
	addrMQTT := NewIPv7WithSubPort(56, 1, 100, 1883)

	if addrREST.SubPort == addrCoAP.SubPort {
		t.Error("REST y CoAP deben tener sub-puertos distintos")
	}
	if addrCoAP.SubPort == addrMQTT.SubPort {
		t.Error("CoAP y MQTT deben tener sub-puertos distintos")
	}
	// Pero todos tienen la misma identidad IPv7
	if !addrREST.Equals(addrCoAP) || !addrCoAP.Equals(addrMQTT) {
		t.Error("Nodos con mismo (r,s,d) deben tener la misma ResolvedIP (identidad) independiente del sub-puerto")
	}
}


func TestPQC_SignatureSize(t *testing.T) {
	data := []byte("test-size-check")
	sig := GenerateSignature(data)
	if sig == nil {
		t.Fatal("No se pudo generar firma")
	}
	// ML-DSA-65 produce firmas de exactamente 3309 bytes
	if len(sig) != 3309 {
		t.Errorf("Tamano de firma ML-DSA-65 incorrecto: esperado 3309, obtuvo %d", len(sig))
	}
}

func TestPQC_MultipleMessages_Independent(t *testing.T) {
	msgA := []byte("mensaje-A")
	msgB := []byte("mensaje-B")

	sigA := GenerateSignature(msgA)
	sigB := GenerateSignature(msgB)
	pkBytes := GetPublicKey()

	if !VerifySignature(pkBytes, msgA, sigA) {
		t.Error("Mensaje A no verifica con su propia firma")
	}
	if !VerifySignature(pkBytes, msgB, sigB) {
		t.Error("Mensaje B no verifica con su propia firma")
	}
	// Cross-firma: firma de A no debe verificar mensaje B
	if VerifySignature(pkBytes, msgB, sigA) {
		t.Error("La firma de A no debe verificar el mensaje B")
	}
}

// ─── Tests Reed-Solomon ECC (mat.md §Códigos ECC) ────────────────────────────

func TestRS_Encode_Length(t *testing.T) {
	msg := []byte("paquete-ipv7-test")
	nsym := 4
	encoded := RSEncode(msg, nsym)
	if len(encoded) != len(msg)+nsym {
		t.Errorf("RS encode: esperado %d bytes, obtuvo %d", len(msg)+nsym, len(encoded))
	}
}

func TestRS_NoError_Detection(t *testing.T) {
	msg := []byte("paquete-sin-errores")
	encoded := RSEncode(msg, 6)
	if RSCanDetectError(encoded, 6) {
		t.Error("RS: mensaje sin errores reporto error (falso positivo)")
	}
}

func TestRS_Error_Detected(t *testing.T) {
	msg := []byte("paquete-con-error")
	encoded := RSEncode(msg, 6)
	// Corromper un byte de paridad
	encoded[len(encoded)-1] ^= 0xFF
	if !RSCanDetectError(encoded, 6) {
		t.Error("RS: error de byte de paridad no fue detectado")
	}
}

func TestRS_ZeroSymbols_Passthrough(t *testing.T) {
	msg := []byte("mensaje-directo")
	result := RSEncode(msg, 0)
	if len(result) != len(msg) {
		t.Error("RS con nsym=0 debe retornar el mensaje sin modificar")
	}
}

// ─── Tests Delta Encoding de Headers (mat.md §Entropía de Shannon) ───────────

func TestDeltaHeader_FirstEncode_IsFullHeader(t *testing.T) {
	state := &HeaderDeltaState{}
	addr := NewIPv7WithSubPort(56, 1, 100, 8080)
	buf := state.EncodeDelta(addr)
	// Primera codificación siempre debe ser el header completo (10 bytes)
	if len(buf) != 10 {
		t.Errorf("Primera codificacion debe ser 10 bytes, obtuvo %d", len(buf))
	}
}

func TestDeltaHeader_SameRegionSubnet_IsDelta(t *testing.T) {
	state := &HeaderDeltaState{}
	addr1 := NewIPv7WithSubPort(56, 1, 100, 0)
	addr2 := NewIPv7WithSubPort(56, 1, 200, 0) // Mismo Region+Subnet, distinto DeviceID

	state.EncodeDelta(addr1) // Primera: inicializa estado
	buf2 := state.EncodeDelta(addr2)

	// Segunda con misma Region+Subnet debe ser delta (6 bytes)
	if len(buf2) != 6 {
		t.Errorf("Delta encoding debe producir 6 bytes, obtuvo %d", len(buf2))
	}
}

func TestDeltaHeader_Roundtrip(t *testing.T) {
	encState := &HeaderDeltaState{}
	decState := &HeaderDeltaState{}

	original := NewIPv7WithSubPort(56, 1, 999, 7780)

	buf1 := encState.EncodeDelta(original)
	decoded1 := decState.DecodeDelta(buf1)
	if !decoded1.Equals(original) {
		t.Errorf("Roundtrip completo fallo: original=%.0f decoded=%.0f", original.ResolvedIP, decoded1.ResolvedIP)
	}

	// Segunda vuelta (delta)
	original2 := NewIPv7WithSubPort(56, 1, 777, 7780)
	buf2 := encState.EncodeDelta(original2)
	decoded2 := decState.DecodeDelta(buf2)
	if !decoded2.Equals(original2) {
		t.Errorf("Roundtrip delta fallo: original=%.0f decoded=%.0f", original2.ResolvedIP, decoded2.ResolvedIP)
	}
}

// ─── Tests Predictor Estocástico Fokker-Planck (mat.md §Kolmogorov) ──────────

func TestStochastic_Predict_StableSignal(t *testing.T) {
	p := NewStochasticPredictor(16)
	// Señal estable: latencia constante de 30ms
	for i := 0; i < 12; i++ {
		p.Push(30.0)
	}
	mean, std := p.Predict(5)
	if mean < 20 || mean > 40 {
		t.Errorf("Señal estable: media predicha deberia ser ~30ms, obtuvo %.2f", mean)
	}
	if std > 10 {
		t.Errorf("Señal estable: desviacion deberia ser baja, obtuvo %.2f", std)
	}
}

func TestStochastic_Predict_RisingSignal(t *testing.T) {
	p := NewStochasticPredictor(16)
	// Señal creciente: latencia subiendo de 20 a 80ms
	for i := 0; i < 12; i++ {
		p.Push(float64(20 + i*5))
	}
	mean, _ := p.Predict(3)
	// Media predicha deberia ser mayor que la ultima muestra (75ms)
	if mean < 75 {
		t.Errorf("Señal creciente: media predicha deberia ser >75ms, obtuvo %.2f", mean)
	}
}

func TestStochastic_HighRisk_HighLatency(t *testing.T) {
	p := NewStochasticPredictor(16)
	// Señal que ya supera el umbral
	for i := 0; i < 10; i++ {
		p.Push(200.0)
	}
	risk := p.IsHighRisk(1, 100.0) // Umbral 100ms, ya estamos en 200ms
	if risk < 0.5 {
		t.Errorf("Latencia en 200ms deberia tener alto riesgo de superar 100ms, obtuvo %.2f", risk)
	}
}

func TestStochastic_LowRisk_LowLatency(t *testing.T) {
	p := NewStochasticPredictor(16)
	for i := 0; i < 10; i++ {
		p.Push(10.0)
	}
	risk := p.IsHighRisk(5, 500.0) // Umbral 500ms, estamos en 10ms
	if risk > 0.1 {
		t.Errorf("Latencia baja deberia tener bajo riesgo de superar 500ms, obtuvo %.2f", risk)
	}
}
