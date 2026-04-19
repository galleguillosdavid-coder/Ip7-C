//go:build !windows

package main

// wintunDLL es nil en plataformas que no son Windows.
// El adaptador TUN se gestiona via el driver del kernel del SO correspondiente.
var wintunDLL []byte
