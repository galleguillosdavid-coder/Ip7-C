//go:build !windows

package main

// ensureSecurityExceptions es un no-op en plataformas que no son Windows.
// En Linux/macOS, los permisos se gestionan via capabilities (CAP_NET_ADMIN) o sudo.
func ensureSecurityExceptions() {}

// extractWintun es un no-op fuera de Windows — el driver TUN es nativo del kernel.
func extractWintun(kernelMode string) {}
