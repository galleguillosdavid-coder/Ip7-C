package protocol

// ─── Traffic Classes (QoS) — Necesidades.md §Imposibilidad Técnica QoS ─────────────────────
//
// Clasifica el tráfico para que el MoE Dispatcher y el Network Slice puedan
// priorizar de extremo a extremo, independientemente del tipo de dispositivo.
// Se codifica en los 3 bits más altos del SubPort (bits 15-13), dejando
// 8192 canales lógicos por clase (bits 12-0).
//
// Por qué 4 clases y no más: reflect DiffServ (RFC 2474) simplificado.
// DiffServ define 6 classes pero en redes overlay la distinción práctica es 4.

type TrafficClass uint8

const (
	// Implementación Decit Cuántico: 10 Niveles de Tráfico en lugar de 4.
	// Emula lógica Multivaluada (Álgebra de Post/Lukasiewicz).
	TC_Q0_CRITICAL TrafficClass = 0
	TC_Q1_REALTIME TrafficClass = 1
	TC_Q2_CONTROL  TrafficClass = 2
	TC_Q3_SIG      TrafficClass = 3
	TC_Q4_BULK     TrafficClass = 4
	TC_Q5_BACKGROUND TrafficClass = 5
	// TC_Q6 a Q9 reservados para experimentación interferencial AI agentica
	TC_Q6 TrafficClass = 6
	TC_Q7 TrafficClass = 7
	TC_Q8 TrafficClass = 8
	TC_Q9 TrafficClass = 9

	// Alias legados para Slices autogenerados
	TC_CONTROL    = TC_Q0_CRITICAL
	TC_REALTIME   = TC_Q1_REALTIME
	TC_BULK       = TC_Q4_BULK
	TC_BACKGROUND = TC_Q5_BACKGROUND
)

// TCFromSubPort extrae la TrafficClass del campo SubPort usando enmascaramiento Decit.
// Los bits 15-12 codifican la clase (4 bits = hasta 16 clases, usamos 10 decits).
func TCFromSubPort(subPort uint16) TrafficClass {
	return TrafficClass((subPort >> 12) & 0x0F)
}

// SubPortWithTC construye un SubPort codificando la TrafficClass P-bit en los bits 15-12
// y el canal lógico (0-4095) en bits 11-0.
func SubPortWithTC(tc TrafficClass, channel uint16) uint16 {
	return (uint16(tc) << 12) | (channel & 0x0FFF)
}

// ─── Network Slice Profiles (Previsiones.md §IMT-2030 6G Network Slicing) ─────────────────
//
// "Un corte de red para telecirugía blindará latencia <1ms; otro manejará firmware OTA;
//  un tercero sostendrá XR multi-gigabit — sin interrumpirse mutuamente."
//
// En IPv7-IEU un Slice combina DeviceClass + TrafficClass para configurar
// automáticamente el comportamiento del tunnel, DHT y MoE Dispatcher.

// SliceProfile define el comportamiento completo de un slice de red lógico.
type SliceProfile struct {
	Name              string
	Device            DeviceClass
	TC                TrafficClass
	MaxLatencyMs      float64 // Umbral de latencia para el MoE Dispatcher
	RSSymbols         int     // Bytes de paridad Reed-Solomon (0 = desactivado)
	DHTBudget         int     // RPCs/ventana para el Task Budget del DHT
	DHTWindow         int64   // Ventana en segundos
	DeltaHeader       bool    // Activar delta encoding del header
	Description       string
}

// Catálogo de slices predefinidos — mapeo Device+TC → comportamiento
var SliceProfiles = map[string]SliceProfile{
	// ── Infraestructura ────────────────────────────────────────────────────────
	"router-control": {
		Name: "Router / Control Plane", Device: DeviceRouter, TC: TC_CONTROL,
		MaxLatencyMs: 5, RSSymbols: 0, DHTBudget: 300, DHTWindow: 60,
		DeltaHeader: true, Description: "Señalización de alta frecuencia entre routers",
	},
	"router-bulk": {
		Name: "Router / Data Plane", Device: DeviceRouter, TC: TC_BULK,
		MaxLatencyMs: 15, RSSymbols: 4, DHTBudget: 200, DHTWindow: 60,
		DeltaHeader: true, Description: "Enrutamiento de datos a velocidad de línea",
	},

	// ── Servidores ─────────────────────────────────────────────────────────────
	"server-realtime": {
		Name: "Servidor / Tiempo Real", Device: DeviceServer, TC: TC_REALTIME,
		MaxLatencyMs: 20, RSSymbols: 2, DHTBudget: 240, DHTWindow: 60,
		DeltaHeader: false, Description: "APIs de latencia crítica y streaming",
	},
	"server-bulk": {
		Name: "Servidor / Transferencia", Device: DeviceServer, TC: TC_BULK,
		MaxLatencyMs: 50, RSSymbols: 8, DHTBudget: 120, DHTWindow: 60,
		DeltaHeader: true, Description: "Transferencia masiva de datos",
	},

	// ── Dispositivos móviles ───────────────────────────────────────────────────
	"mobile-realtime": {
		Name: "Móvil / Tiempo Real", Device: DeviceMobile, TC: TC_REALTIME,
		MaxLatencyMs: 60, RSSymbols: 4, DHTBudget: 20, DHTWindow: 60,
		DeltaHeader: true, Description: "Videollamadas y juegos en smartphone",
	},
	"mobile-background": {
		Name: "Móvil / Background", Device: DeviceMobile, TC: TC_BACKGROUND,
		MaxLatencyMs: 500, RSSymbols: 0, DHTBudget: 5, DHTWindow: 300,
		DeltaHeader: true, Description: "Sincronización de apps en segundo plano",
	},

	// ── IoT ────────────────────────────────────────────────────────────────────
	"iot-sensor-bg": {
		Name: "IoT Sensor / Background", Device: DeviceIoTSensor, TC: TC_BACKGROUND,
		MaxLatencyMs: 1000, RSSymbols: 2, DHTBudget: 2, DHTWindow: 300,
		DeltaHeader: true, Description: "Telemetría pasiva de sensores ultra-bajo consumo",
	},
	"iot-actuator-rt": {
		Name: "IoT Actuador / Tiempo Real", Device: DeviceIoTActuator, TC: TC_REALTIME,
		MaxLatencyMs: 50, RSSymbols: 4, DHTBudget: 15, DHTWindow: 60,
		DeltaHeader: false, Description: "Comandos críticos a actuadores industriales",
	},
	"industrial-control": {
		Name: "Industrial / SCADA Control", Device: DeviceIndustrial, TC: TC_CONTROL,
		MaxLatencyMs: 5, RSSymbols: 8, DHTBudget: 40, DHTWindow: 60,
		DeltaHeader: false, Description: "Señalización SCADA de tiempo crítico",
	},

	// ── Satelital ──────────────────────────────────────────────────────────────
	"leo-bulk": {
		Name: "LEO / Transferencia", Device: DeviceSatelliteLEO, TC: TC_BULK,
		MaxLatencyMs: 60, RSSymbols: 12, DHTBudget: 20, DHTWindow: 60,
		DeltaHeader: true, Description: "Datos en masa sobre Starlink LEO",
	},
	"geo-background": {
		Name: "GEO / Background", Device: DeviceSatelliteGEO, TC: TC_BACKGROUND,
		MaxLatencyMs: 800, RSSymbols: 16, DHTBudget: 3, DHTWindow: 300,
		DeltaHeader: true, Description: "Telemetría sobre satélite geoestacionario",
	},

	// ── Notebook / Desktop ─────────────────────────────────────────────────────
	"notebook-bulk": {
		Name: "Notebook / Transferencia", Device: DeviceNotebook, TC: TC_BULK,
		MaxLatencyMs: 80, RSSymbols: 4, DHTBudget: 60, DHTWindow: 60,
		DeltaHeader: true, Description: "Descargas y transferencias desde laptop",
	},
	"desktop-realtime": {
		Name: "Desktop / Tiempo Real", Device: DeviceDesktop, TC: TC_REALTIME,
		MaxLatencyMs: 30, RSSymbols: 2, DHTBudget: 100, DHTWindow: 60,
		DeltaHeader: false, Description: "Aplicaciones de tiempo real en PC",
	},

	// ── Impresora ──────────────────────────────────────────────────────────────
	"printer-bulk": {
		Name: "Impresora / Datos", Device: DevicePrinter, TC: TC_BULK,
		MaxLatencyMs: 500, RSSymbols: 8, DHTBudget: 3, DHTWindow: 60,
		DeltaHeader: true, Description: "Trabajos de impresión en red",
	},

	// ── Cámara IP ──────────────────────────────────────────────────────────────
	"camera-realtime": {
		Name: "Cámara IP / Stream", Device: DeviceCamera, TC: TC_REALTIME,
		MaxLatencyMs: 80, RSSymbols: 4, DHTBudget: 10, DHTWindow: 60,
		DeltaHeader: false, Description: "Stream de video de vigilancia en tiempo real",
	},

	// ── Smart Home ─────────────────────────────────────────────────────────────
	"smarthome-control": {
		Name: "Smart Home / Control", Device: DeviceSmartHome, TC: TC_CONTROL,
		MaxLatencyMs: 200, RSSymbols: 2, DHTBudget: 5, DHTWindow: 120,
		DeltaHeader: true, Description: "Comandos a dispositivos domóticos",
	},

	// ── Vehículo / Drone ───────────────────────────────────────────────────────
	"vehicle-realtime": {
		Name: "Vehículo / Tiempo Real", Device: DeviceVehicle, TC: TC_REALTIME,
		MaxLatencyMs: 30, RSSymbols: 4, DHTBudget: 30, DHTWindow: 60,
		DeltaHeader: false, Description: "Telemetría vehicular y V2X en tiempo real",
	},
	"drone-control": {
		Name: "Drone / Control", Device: DeviceDrone, TC: TC_CONTROL,
		MaxLatencyMs: 20, RSSymbols: 6, DHTBudget: 10, DHTWindow: 60,
		DeltaHeader: false, Description: "Señales de control de vuelo de UAV",
	},
}

// GetSliceProfile devuelve el perfil de slice para una combinación Device+TC.
// Busca el match más específico; si no existe, genera uno genérico.
func GetSliceProfile(device DeviceClass, tc TrafficClass) SliceProfile {
	devProfile := GetDeviceProfile(device)
	for _, sp := range SliceProfiles {
		if sp.Device == device && sp.TC == tc {
			return sp
		}
	}
	// Fallback genérico
	return SliceProfile{
		Name:         devProfile.Name + " / Genérico",
		Device:       device,
		TC:           tc,
		MaxLatencyMs: devProfile.LatencyThresholdMs,
		RSSymbols:    4,
		DHTBudget:    devProfile.RecommendedDHTBudget,
		DHTWindow:    60,
		DeltaHeader:  true,
		Description:  "Perfil genérico auto-generado",
	}
}
