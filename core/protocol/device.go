package protocol

import "time"

// DeviceClass clasifica el tipo de dispositivo conectado a la red IPv7-IEU.
// Cada clase tiene perfiles distintos de latencia, ancho de banda, MTU,
// presupuesto DHT y comportamiento del MoE Expert Dispatcher.
type DeviceClass int

const (
	// Genérico / desconocido
	DeviceUnknown DeviceClass = iota

	// ── Infraestructura de Red ──────────────────────────────────────────────
	DeviceRouter      // Routers, gateways, firewalls (Mikrotik, Cisco, TP-Link)
	DeviceNAT         // Dispositivos NAT puro / CGNAT / carrier-grade NAT
	DeviceAccessPoint // Puntos de acceso WiFi 5/6/6E, small cells 5G

	// ── Servidores ──────────────────────────────────────────────────────────
	DeviceServer    // Servidores web, DB, API, cloud (AWS, VPS, bare metal)
	DeviceNAS       // Almacenamiento en red (Synology, QNAP, TrueNAS)
	DeviceEdge      // Computación de borde (Raspberry Pi, Bmax, NVIDIA Jetson)

	// ── Dispositivos de Usuario ─────────────────────────────────────────────
	DeviceDesktop   // PCs de escritorio (alto rendimiento, fibra fija)
	DeviceNotebook  // Laptops/notebooks con batería
	DeviceMobile    // Smartphones (iOS/Android) y tablets con LTE/5G/WiFi
	DeviceWearable  // Smartwatches, pulseras, auriculares BT con WiFi

	// ── Entretenimiento y Hogar ─────────────────────────────────────────────
	DeviceSmartTV   // Smart TVs, Chromecast, Apple TV, Fire Stick
	DeviceConsole   // Consolas de videojuegos (PS5, Xbox, Switch)

	// ── Periféricos de Red ──────────────────────────────────────────────────
	DevicePrinter   // Impresoras de red (láser, inyección, multifunción)
	DeviceCamera    // Cámaras IP, cámaras de seguridad, dashcams en red

	// ── IoT ─────────────────────────────────────────────────────────────────
	DeviceIoTSensor   // Sensores ultra-bajo consumo (temperatura, CO2, movimiento)
	DeviceIoTActuator // Actuadores, válvulas, relés, controladores industriales
	DeviceSmartHome   // Domótica (termostatos, enchufes, bombillas, alarmas)
	DeviceIndustrial  // SCADA, PLCs, equipos industriales, ICS/OT

	// ── Conectividad Satelital y Remota ─────────────────────────────────────
	DeviceSatelliteLEO // Starlink, OneWeb, Amazon Kuiper (órbita baja, 20-40ms)
	DeviceSatelliteGEO // Satélite geoestacionario clásico (600-800ms RTT)
	DeviceSatelliteMEO // Órbita media (O3b, SES) - 120-150ms RTT

	// ── Vehículos y Movilidad ────────────────────────────────────────────────
	DeviceVehicle    // Automóviles conectados, buses, trenes con 5G/LTE
	DeviceDrone      // UAVs/drones con enlace de datos
)

// DeviceProfile describe las características de red de una clase de dispositivo.
type DeviceProfile struct {
	Class               DeviceClass
	Name                string
	TypicalLatencyMs    float64 // RTT mediano esperado (ms)
	MaxBandwidthKbps    int     // Ancho de banda de subida máximo típico (Kbps)
	MTUBytes            int     // MTU recomendado para este dispositivo
	BatteryPowered      bool    // Si es a batería, reducir agresividad de red
	RecommendedDHTBudget int    // RPCs/minuto para el Task Budget del DHT
	RetryMax            int     // Máximo de reintentos de paquete por el MoE
	LatencyThresholdMs  float64 // Umbral de alta latencia para activar ExpertSatellite
	BudgetWindow        time.Duration
}

// deviceProfiles es el registro completo de perfiles de dispositivo.
var deviceProfiles = map[DeviceClass]DeviceProfile{
	DeviceUnknown: {
		Name: "Dispositivo Desconocido", TypicalLatencyMs: 50, MaxBandwidthKbps: 10000,
		MTUBytes: 1400, BatteryPowered: false, RecommendedDHTBudget: 60, RetryMax: 2,
		LatencyThresholdMs: 100, BudgetWindow: time.Minute,
	},

	// ── Infraestructura ─────────────────────────────────────────────────────
	DeviceRouter: {
		Name: "Router / Gateway", TypicalLatencyMs: 2, MaxBandwidthKbps: 1000000,
		MTUBytes: 1500, BatteryPowered: false, RecommendedDHTBudget: 200, RetryMax: 1,
		LatencyThresholdMs: 15, BudgetWindow: time.Minute,
	},
	DeviceNAT: {
		Name: "NAT / CGNAT", TypicalLatencyMs: 5, MaxBandwidthKbps: 500000,
		MTUBytes: 1400, BatteryPowered: false, RecommendedDHTBudget: 150, RetryMax: 2,
		LatencyThresholdMs: 20, BudgetWindow: time.Minute,
	},
	DeviceAccessPoint: {
		Name: "Access Point WiFi/5G", TypicalLatencyMs: 3, MaxBandwidthKbps: 500000,
		MTUBytes: 1500, BatteryPowered: false, RecommendedDHTBudget: 120, RetryMax: 2,
		LatencyThresholdMs: 20, BudgetWindow: time.Minute,
	},

	// ── Servidores ───────────────────────────────────────────────────────────
	DeviceServer: {
		Name: "Servidor (Web/DB/API)", TypicalLatencyMs: 5, MaxBandwidthKbps: 1000000,
		MTUBytes: 9000, BatteryPowered: false, RecommendedDHTBudget: 240, RetryMax: 1,
		LatencyThresholdMs: 20, BudgetWindow: time.Minute,
	},
	DeviceNAS: {
		Name: "NAS (Almacenamiento)", TypicalLatencyMs: 8, MaxBandwidthKbps: 125000,
		MTUBytes: 9000, BatteryPowered: false, RecommendedDHTBudget: 30, RetryMax: 3,
		LatencyThresholdMs: 30, BudgetWindow: time.Minute,
	},
	DeviceEdge: {
		Name: "Edge Computing (Bmax/Jetson/Pi)", TypicalLatencyMs: 10, MaxBandwidthKbps: 100000,
		MTUBytes: 1400, BatteryPowered: true, RecommendedDHTBudget: 40, RetryMax: 2,
		LatencyThresholdMs: 50, BudgetWindow: time.Minute,
	},

	// ── Dispositivos de Usuario ──────────────────────────────────────────────
	DeviceDesktop: {
		Name: "PC de Escritorio", TypicalLatencyMs: 10, MaxBandwidthKbps: 500000,
		MTUBytes: 1500, BatteryPowered: false, RecommendedDHTBudget: 120, RetryMax: 2,
		LatencyThresholdMs: 50, BudgetWindow: time.Minute,
	},
	DeviceNotebook: {
		Name: "Notebook / Laptop", TypicalLatencyMs: 15, MaxBandwidthKbps: 100000,
		MTUBytes: 1400, BatteryPowered: true, RecommendedDHTBudget: 60, RetryMax: 2,
		LatencyThresholdMs: 60, BudgetWindow: time.Minute,
	},
	DeviceMobile: {
		Name: "Smartphone / Tablet", TypicalLatencyMs: 30, MaxBandwidthKbps: 50000,
		MTUBytes: 1280, BatteryPowered: true, RecommendedDHTBudget: 20, RetryMax: 3,
		LatencyThresholdMs: 80, BudgetWindow: time.Minute,
	},
	DeviceWearable: {
		Name: "Wearable (Smartwatch/Pulsera)", TypicalLatencyMs: 50, MaxBandwidthKbps: 1000,
		MTUBytes: 512, BatteryPowered: true, RecommendedDHTBudget: 5, RetryMax: 1,
		LatencyThresholdMs: 100, BudgetWindow: time.Minute,
	},

	// ── Entretenimiento ──────────────────────────────────────────────────────
	DeviceSmartTV: {
		Name: "Smart TV / Streaming", TypicalLatencyMs: 20, MaxBandwidthKbps: 100000,
		MTUBytes: 1500, BatteryPowered: false, RecommendedDHTBudget: 15, RetryMax: 2,
		LatencyThresholdMs: 60, BudgetWindow: time.Minute,
	},
	DeviceConsole: {
		Name: "Consola de Videojuegos", TypicalLatencyMs: 25, MaxBandwidthKbps: 200000,
		MTUBytes: 1400, BatteryPowered: false, RecommendedDHTBudget: 30, RetryMax: 1,
		LatencyThresholdMs: 40, BudgetWindow: time.Minute,
	},

	// ── Periféricos ──────────────────────────────────────────────────────────
	DevicePrinter: {
		Name: "Impresora de Red", TypicalLatencyMs: 20, MaxBandwidthKbps: 5000,
		MTUBytes: 1400, BatteryPowered: false, RecommendedDHTBudget: 3, RetryMax: 5,
		LatencyThresholdMs: 200, BudgetWindow: time.Minute,
	},
	DeviceCamera: {
		Name: "Cámara IP / Seguridad", TypicalLatencyMs: 15, MaxBandwidthKbps: 8000,
		MTUBytes: 1400, BatteryPowered: false, RecommendedDHTBudget: 10, RetryMax: 1,
		LatencyThresholdMs: 100, BudgetWindow: time.Minute,
	},

	// ── IoT ──────────────────────────────────────────────────────────────────
	DeviceIoTSensor: {
		Name: "Sensor IoT (Ultra-bajo consumo)", TypicalLatencyMs: 100, MaxBandwidthKbps: 250,
		MTUBytes: 256, BatteryPowered: true, RecommendedDHTBudget: 2, RetryMax: 1,
		LatencyThresholdMs: 500, BudgetWindow: 5 * time.Minute,
	},
	DeviceIoTActuator: {
		Name: "Actuador IoT / Controlador", TypicalLatencyMs: 50, MaxBandwidthKbps: 1000,
		MTUBytes: 512, BatteryPowered: false, RecommendedDHTBudget: 10, RetryMax: 3,
		LatencyThresholdMs: 200, BudgetWindow: time.Minute,
	},
	DeviceSmartHome: {
		Name: "Domótica (Smart Home)", TypicalLatencyMs: 40, MaxBandwidthKbps: 2000,
		MTUBytes: 512, BatteryPowered: true, RecommendedDHTBudget: 5, RetryMax: 2,
		LatencyThresholdMs: 300, BudgetWindow: 2 * time.Minute,
	},
	DeviceIndustrial: {
		Name: "Industrial (SCADA/PLC/ICS)", TypicalLatencyMs: 5, MaxBandwidthKbps: 10000,
		MTUBytes: 1400, BatteryPowered: false, RecommendedDHTBudget: 30, RetryMax: 5,
		LatencyThresholdMs: 20, BudgetWindow: time.Minute,
	},

	// ── Satelital ────────────────────────────────────────────────────────────
	DeviceSatelliteLEO: {
		Name: "Satélite LEO (Starlink/OneWeb/Kuiper)", TypicalLatencyMs: 30, MaxBandwidthKbps: 200000,
		MTUBytes: 1400, BatteryPowered: false, RecommendedDHTBudget: 20, RetryMax: 3,
		LatencyThresholdMs: 60, BudgetWindow: time.Minute,
	},
	DeviceSatelliteGEO: {
		Name: "Satélite GEO (Viasat/HughesNet)", TypicalLatencyMs: 700, MaxBandwidthKbps: 25000,
		MTUBytes: 1200, BatteryPowered: false, RecommendedDHTBudget: 5, RetryMax: 1,
		LatencyThresholdMs: 500, BudgetWindow: 5 * time.Minute,
	},
	DeviceSatelliteMEO: {
		Name: "Satélite MEO (O3b/SES)", TypicalLatencyMs: 140, MaxBandwidthKbps: 100000,
		MTUBytes: 1300, BatteryPowered: false, RecommendedDHTBudget: 15, RetryMax: 2,
		LatencyThresholdMs: 200, BudgetWindow: 2 * time.Minute,
	},

	// ── Movilidad ────────────────────────────────────────────────────────────
	DeviceVehicle: {
		Name: "Vehículo Conectado (5G/LTE)", TypicalLatencyMs: 20, MaxBandwidthKbps: 50000,
		MTUBytes: 1280, BatteryPowered: false, RecommendedDHTBudget: 30, RetryMax: 3,
		LatencyThresholdMs: 70, BudgetWindow: time.Minute,
	},
	DeviceDrone: {
		Name: "Drone / UAV", TypicalLatencyMs: 25, MaxBandwidthKbps: 15000,
		MTUBytes: 1000, BatteryPowered: true, RecommendedDHTBudget: 10, RetryMax: 2,
		LatencyThresholdMs: 60, BudgetWindow: time.Minute,
	},
}

// GetDeviceProfile devuelve el perfil de un dispositivo por su clase.
// Si no se encuentra, retorna el perfil de DeviceUnknown.
func GetDeviceProfile(class DeviceClass) DeviceProfile {
	if p, ok := deviceProfiles[class]; ok {
		p.Class = class
		return p
	}
	p := deviceProfiles[DeviceUnknown]
	p.Class = DeviceUnknown
	return p
}

// ParseDeviceClass convierte un string en DeviceClass.
// Permite usar el flag --device desde la línea de comandos.
func ParseDeviceClass(name string) DeviceClass {
	mapping := map[string]DeviceClass{
		"unknown":     DeviceUnknown,
		"router":      DeviceRouter,
		"nat":         DeviceNAT,
		"ap":          DeviceAccessPoint,
		"server":      DeviceServer,
		"nas":         DeviceNAS,
		"edge":        DeviceEdge,
		"desktop":     DeviceDesktop,
		"notebook":    DeviceNotebook,
		"laptop":      DeviceNotebook,
		"mobile":      DeviceMobile,
		"phone":       DeviceMobile,
		"tablet":      DeviceMobile,
		"wearable":    DeviceWearable,
		"smartwatch":  DeviceWearable,
		"smarttv":     DeviceSmartTV,
		"tv":          DeviceSmartTV,
		"console":     DeviceConsole,
		"printer":     DevicePrinter,
		"camera":      DeviceCamera,
		"iot-sensor":  DeviceIoTSensor,
		"sensor":      DeviceIoTSensor,
		"iot":         DeviceIoTSensor,
		"actuator":    DeviceIoTActuator,
		"smarthome":   DeviceSmartHome,
		"industrial":  DeviceIndustrial,
		"scada":       DeviceIndustrial,
		"plc":         DeviceIndustrial,
		"leo":         DeviceSatelliteLEO,
		"starlink":    DeviceSatelliteLEO,
		"geo":         DeviceSatelliteGEO,
		"meo":         DeviceSatelliteMEO,
		"vehicle":     DeviceVehicle,
		"car":         DeviceVehicle,
		"drone":       DeviceDrone,
		"uav":         DeviceDrone,
	}
	if class, ok := mapping[name]; ok {
		return class
	}
	return DeviceUnknown
}

// ListDeviceClasses devuelve todos los nombres de dispositivo aceptados por ParseDeviceClass.
func ListDeviceClasses() []string {
	return []string{
		"router", "nat", "ap", "server", "nas", "edge",
		"desktop", "notebook", "laptop", "mobile", "phone", "tablet",
		"wearable", "smartwatch", "smarttv", "console", "printer", "camera",
		"iot-sensor", "sensor", "iot", "actuator", "smarthome",
		"industrial", "scada", "plc",
		"leo", "starlink", "geo", "meo",
		"vehicle", "car", "drone", "uav",
	}
}
