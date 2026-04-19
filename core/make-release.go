//go:build ignore

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	binName := fmt.Sprintf("ipv7-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	
	fmt.Println("🚀 Iniciando empaquetado de Release Oficial IPv7-IEU...")

	// 1. Compilar el binario
	fmt.Printf("🔨 Compilando ejecutable: %s\n", binName)
	cmd := exec.Command("go", "build", "-ldflags", "-s -w", "-o", binName, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "." // asume que se corre dentro de core/
	
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ Error al compilar: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Binario %s compilado con éxito.\n", binName)

	// 2. Leer el binario para verificar
	data, err := os.ReadFile(binName)
	if err != nil {
		fmt.Printf("❌ Error al leer %s: %v\n", binName, err)
		os.Exit(1)
	}

	// 3. Generar SHA-256
	fmt.Println("✍️ Generando hash SHA-256...")
	hasher := sha256.New()
	hasher.Write(data)
	hashSum := hex.EncodeToString(hasher.Sum(nil))

	// 4. Escribir archivo sum
	sigPath := binName + ".sha256"
	if err := os.WriteFile(sigPath, []byte(hashSum+"  "+binName+"\n"), 0644); err != nil {
		fmt.Printf("❌ Error guardando hash: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Hash SHA-256 guardado exitosamente en %s\n", sigPath)

	fmt.Println("\n🎉 PROCESO COMPLETADO. Tienes tus artefactos listos:")
	absPath, _ := filepath.Abs(binName)
	absSig, _ := filepath.Abs(sigPath)
	fmt.Printf("- %s\n", absPath)
	fmt.Printf("- %s\n", absSig)
}
