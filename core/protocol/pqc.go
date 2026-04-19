package protocol

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudflare/circl/sign/mldsa/mldsa65"
)

const (
	keyDir  = ".ipv7-pqc"
	privKey = "ipv7-pqc.priv"
	pubKey  = "ipv7-pqc.pub"
)

var pub *mldsa65.PublicKey
var priv *mldsa65.PrivateKey

// LoadOrGenerateKeys carga claves desde disco o genera y guarda nuevo par
func LoadOrGenerateKeys() error {
	keyPath := filepath.Join(keyDir, privKey)
	pubPath := filepath.Join(keyDir, pubKey)

	// Intentar cargar
	if data, err := os.ReadFile(keyPath); err == nil {
		priv = new(mldsa65.PrivateKey)
		if err := priv.UnmarshalBinary(data); err == nil {
			pub = priv.Public().(*mldsa65.PublicKey) // derivar pública
			fmt.Println("✅ Claves PQC cargadas desde disco")
			return nil
		}
	}

	// Generar nuevas
	fmt.Println("🔑 Generando nuevo par ML-DSA-65...")
	var err error
	pub, priv, err = mldsa65.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	// Guardar
	os.MkdirAll(keyDir, 0700)
	privBytes, _ := priv.MarshalBinary()
	pubBytes, _ := pub.MarshalBinary()

	os.WriteFile(keyPath, privBytes, 0600)
	os.WriteFile(pubPath, pubBytes, 0644)

	fmt.Println("✅ Nuevo par PQC generado y persistido en disco")
	return nil
}

func init() {
	if err := LoadOrGenerateKeys(); err != nil {
		fmt.Printf("❌ [PQC] Error crítico cargando claves ML-DSA-65: %v\n", err)
		fmt.Println("⚠️ [PQC] El nodo operará sin firma post-cuántica hasta resolver el error de disco/permisos.")
	}
}

// GenerateSignature ...
func GenerateSignature(data []byte) []byte {
	if priv == nil {
		return nil
	}
	sig, _ := priv.Sign(rand.Reader, data, nil)
	return sig
}

// VerifySignature ...
func VerifySignature(pkBytes, msg, sig []byte) bool {
	var pk mldsa65.PublicKey
	if err := pk.UnmarshalBinary(pkBytes); err != nil {
		return false
	}
	return mldsa65.Verify(&pk, msg, nil, sig)
}

// GetPublicKey ...
func GetPublicKey() []byte {
	if pub == nil {
		return nil
	}
	b, _ := pub.MarshalBinary()
	return b
}
