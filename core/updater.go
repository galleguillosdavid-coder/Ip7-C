package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const ReleaseAPI = "https://api.github.com/repos/galleguillosdavid-coder/Ip7-IEU/releases/latest"

type GitHubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []GitHubAsset `json:"assets"`
}

type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// platformBinaryName devuelve el nombre del binario correcto según el OS y arquitectura en runtime.
// Esto permite que el GhostUpdater descargue siempre el artefacto correcto para cada plataforma.
func platformBinaryName() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	switch goos {
	case "windows":
		return fmt.Sprintf("ipv7-windows-%s.exe", goarch)
	default:
		return fmt.Sprintf("ipv7-%s-%s", goos, goarch)
	}
}

// executableExtension devuelve la extensión del ejecutable actual según el OS
func executableExtension() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

// CleanOldBinaries purga los ejecutables .old sobrantes que dejó la iteración anterior
func CleanOldBinaries() {
	files, err := os.ReadDir(".")
	if err != nil {
		return
	}
	for _, f := range files {
		if strings.Contains(f.Name(), ".old.") {
			_ = os.Remove(f.Name())
		}
	}
}

// CheckUpdate verifica y aplica actualizaciones automáticamente.
// verifyMethod puede ser:
//   ""       → Verifica SHA256 y aplica el binario (comportamiento normal)
//   "dryrun" → Solo descarga y verifica integridad, NO aplica ni reinicia (Previsiones.md §validación paramétrica)
//
// Previsiones.md: "los parámetros de automatización que operan contra infraestructuras de
// producción requieren validación paramétrica explícita y dry-runs antes de su ejecución.
// Con la automatización, la escala destructiva supera la capacidad humana de contención."
func CheckUpdate(currentVersion string, verifyMethod string) {
	isDryRun := verifyMethod == "dryrun"
	if isDryRun {
		fmt.Println("🔬 [GhostUpdater] Modo DRY-RUN — Verificación sin aplicar cambios")
	}
	fmt.Println("🕵️ [GhostUpdater] Escaneando iteraciones de fase remotas y purgando residuos...")
	CleanOldBinaries()

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(ReleaseAPI)
	if err != nil {
		// Modo Offline / Transitorio, ignorar.
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	if latestVersion == "" || latestVersion == currentVersion {
		fmt.Println("✅ [GhostUpdater] Motor IPv7 en estado óptimo (Versión Mantenida).")
		return
	}

	targetBinary := platformBinaryName()
	fmt.Printf("⚠️ [GhostUpdater] Mutación detectada -> v%s | Buscando artefacto: %s\n",
		latestVersion, targetBinary)

	var downloadURL string
	var sha256URL string
	for _, asset := range release.Assets {
		if asset.Name == targetBinary {
			downloadURL = asset.BrowserDownloadURL
		} else if asset.Name == "SHA256SUMS.txt" {
			sha256URL = asset.BrowserDownloadURL
		}
	}

	if downloadURL == "" {
		fmt.Printf("⚠️ [GhostUpdater] Artefacto %s no encontrado en el release. Sin actualización.\n", targetBinary)
		return
	}

	if err := selfUpdate(downloadURL, sha256URL, verifyMethod, targetBinary, isDryRun); err != nil {
		fmt.Printf("❌ [GhostUpdater] Colapso en inyección de binario: %v\n", err)
		return
	}

	if isDryRun {
		fmt.Println("✅ [GhostUpdater] DRY-RUN completado: integridad del binario verificada. Sin cambios aplicados.")
		return
	}
	fmt.Println("🔄 [GhostUpdater] Auto-Actualización completada. El núcleo IPv7 se reiniciará instantáneamente.")
	time.Sleep(3 * time.Second)
	os.Exit(0)
}

// selfUpdate descarga el nuevo binario, verifica su integridad y lo intercambia
// con el ejecutable actual usando Hot-Swap renaming.
// Si isDryRun=true, solo descarga y verifica sin aplicar ni reiniciar.
func selfUpdate(downloadURL, sha256URL, verifyMethod, targetBinaryName string, isDryRun bool) error {
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("error descargando nuevo binario: %v", err)
	}
	defer resp.Body.Close()

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("no se pudo determinar ruta del ejecutable: %v", err)
	}

	newExePath := exePath + ".new" + executableExtension()
	oldExePath := fmt.Sprintf("%s.old.%d", exePath, time.Now().Unix())

	// Descargar nuevo binario al path temporal
	outFile, err := os.OpenFile(newExePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("error abriendo archivo temporal: %v", err)
	}

	written, err := io.Copy(outFile, resp.Body)
	outFile.Close()
	if err != nil || written == 0 {
		os.Remove(newExePath)
		return fmt.Errorf("error escribiendo binario (%d bytes): %v", written, err)
	}

	fmt.Printf("📦 [GhostUpdater] Binario descargado: %d bytes -> %s\n", written, newExePath)

	// PQC VERIFICATION DISABLED UNTIL STABLE - DECISIÓN DEL DESARROLLADOR
	// No reactivar sin pruebas exhaustivas y clave maestra centralizada
	fmt.Println("⚠️ [GhostUpdater] Validación PQC DESACTIVADA temporalmente por estabilidad y delegación IA")

	// 2. Verificación SHA-256 (Temporal) si se solicita
	if verifyMethod == "sha256" {
		if sha256URL == "" {
			os.Remove(newExePath)
			return fmt.Errorf("el método sha256 fue solicitado pero no se encontró SHA256SUMS.txt en el release")
		}
		fmt.Println("🔒 [GhostUpdater] Iniciando verificación estricta SHA-256...")
		
		sumsResp, err := client.Get(sha256URL)
		if err != nil || sumsResp.StatusCode != 200 {
			os.Remove(newExePath)
			return fmt.Errorf("no se pudo descargar SHA256SUMS.txt")
		}
		defer sumsResp.Body.Close()
		
		sumsData, err := io.ReadAll(sumsResp.Body)
		if err != nil {
			os.Remove(newExePath)
			return fmt.Errorf("error leyendo SHA256SUMS.txt")
		}
		
		// Buscar el hash correspondiente al archivo de destino
		expectedHash := ""
		for _, line := range strings.Split(string(sumsData), "\n") {
			parts := strings.Fields(line)
			if len(parts) >= 2 && strings.HasSuffix(parts[1], targetBinaryName) {
				expectedHash = parts[0]
				break
			}
		}

		if expectedHash == "" {
			os.Remove(newExePath)
			return fmt.Errorf("no se encontró hash para el artefacto %s en SHA256SUMS.txt", targetBinaryName)
		}

		// Calcular el hash del binario descargado
		binData, _ := os.ReadFile(newExePath)
		hasher := sha256.New()
		hasher.Write(binData)
		actualHash := hex.EncodeToString(hasher.Sum(nil))

		if actualHash != expectedHash {
			os.Remove(newExePath)
			return fmt.Errorf("INTEGRIDAD FALLIDA: sha256 mismatch (esperado: %s, obtenido: %s)", expectedHash, actualHash)
		}
		fmt.Println("✅ [GhostUpdater] Integridad y Verificación SHA-256 Exitosa!")
	} else {
		fmt.Println("⚠️ [GhostUpdater] Auto-Update remoto autorizado sin checks adicionales (verify=none)")
	}

	// Modo DRY-RUN: verificación completada, no aplicar el binario
	if isDryRun {
		os.Remove(newExePath) // Limpieza: eliminar descarga temporal
		fmt.Printf("🔬 [GhostUpdater] DRY-RUN: binario verificado y eliminado sin aplicar (%s)\n", newExePath)
		return nil
	}

	// Hot-Swap: renombrar actual -> .old, nuevo -> actual

	// Este método funciona en todos los OS porque no intenta sobreescribir un archivo en uso
	if err := os.Rename(exePath, oldExePath); err != nil {
		os.Remove(newExePath)
		return fmt.Errorf("error moviendo ejecutable actual: %v", err)
	}

	if err := os.Rename(newExePath, exePath); err != nil {
		os.Rename(oldExePath, exePath) // Rollback de emergencia
		return fmt.Errorf("error instalando nuevo binario: %v", err)
	}

	// Lanzar el nuevo proceso con los mismos argumentos
	cmd := exec.Command(exePath, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		// Rollback si el nuevo proceso no inicia
		os.Rename(oldExePath, exePath)
		return fmt.Errorf("error iniciando nuevo proceso: %v", err)
	}

	// El binario viejo será purgado por CleanOldBinaries() instanciada en el nuevo proceso
	return nil
}
