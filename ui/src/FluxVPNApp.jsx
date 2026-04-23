import { useState, useEffect, useCallback } from "react";
import axios from "axios";
import { EventSource } from "eventsource";
import { Line } from "react-chartjs-2";
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend } from "chart.js";

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend);

const NODES = [
  { id: 1, name: "Santiago", country: "Chile", flag: "🇨🇱", latency: 8, load: 34 },
  { id: 2, name: "São Paulo", country: "Brasil", flag: "🇧🇷", latency: 22, load: 67 },
  { id: 3, name: "Miami", country: "EEUU", flag: "🇺🇸", latency: 45, load: 52 },
  { id: 4, name: "Frankfurt", country: "Alemania", flag: "🇩🇪", latency: 142, load: 41 },
  { id: 5, name: "Tokio", country: "Japón", flag: "🇯🇵", latency: 198, load: 28 },
  { id: 6, name: "Starlink LEO", country: "Satelital", flag: "🛰️", latency: 14, load: 12 },
];

export default function FluxVPNApp() {
  const [connected, setConnected] = useState(false);
  const [connecting, setConnecting] = useState(false);
  const [selectedNode, setSelectedNode] = useState(null);
  const [elapsed, setElapsed] = useState(0);
  const [stats, setStats] = useState({ down: 0, up: 0 });
  const [view, setView] = useState("main");
  const [pqcPulse, setPqcPulse] = useState(false);
  const [settings, setSettings] = useState({
    killSwitch: true,
    autoConnect: false,
    starlinkMode: true,
    pqcStrict: true,
    ghostUpdater: true,
  });
  const [peers, setPeers] = useState([]);
  const [status, setStatus] = useState(null);
  const [metricsHistory, setMetricsHistory] = useState([]);

  const formatTime = (s) => {
    const h = Math.floor(s / 3600).toString().padStart(2, "0");
    const m = Math.floor((s % 3600) / 60).toString().padStart(2, "0");
    const sec = (s % 60).toString().padStart(2, "0");
    return `${h}:${m}:${sec}`;
  };

  const handleConnect = useCallback(async () => {
    if (connected) {
      await disconnectVPN();
      setConnecting(false);
      return;
    }
    setConnecting(true);
    try {
      await connectVPN();
      setConnecting(false);
    } catch (err) {
      setConnecting(false);
    }
  }, [connected]);

  useEffect(() => {
    const es = new EventSource(`${API_BASE}/telemetry`);
    es.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setStatus(prev => ({ ...prev, ...data }));
      setMetricsHistory(prev => [...prev.slice(-59), data]); // Mantener últimos 60 puntos
    };
    return () => es.close();
  }, []);

  // PQC re-sign pulse every 30s
  useEffect(() => {
    if (!connected) return;
    const t = setInterval(() => {
      setPqcPulse(true);
      setTimeout(() => setPqcPulse(false), 2000);
    }, 30000);
    // Initial pulse after 3s
    const init = setTimeout(() => {
      setPqcPulse(true);
      setTimeout(() => setPqcPulse(false), 2000);
    }, 3000);
    return () => { clearInterval(t); clearTimeout(init); };
  }, [connected]);

  useEffect(() => {
    if (settings.autoConnect && !connected && !connecting) {
      handleConnect();
    }
  }, [settings.autoConnect, connected, connecting, handleConnect]);

  useEffect(() => {
    const handleKeyDown = (e) => {
      if (e.ctrlKey || e.metaKey) return; // Allow Ctrl shortcuts
      switch (e.key) {
        case 'c':
          handleConnect();
          break;
        case 'n':
          setView("nodes");
          break;
        case 's':
          setView("settings");
          break;
        case 'm':
          setView("main");
          break;
        case 'Escape':
          setView("main");
          break;
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [view, handleConnect]);

  const bg = "#06080d";
  const card = "#0d1219";
  const border = "rgba(0,229,255,0.08)";
  const cyan = "#00e5ff";
  const blue = "#0066ff";
  const textPri = "#e0e6f0";
  const textSec = "#5a6e88";
  const green = "#00e676";
  const red = "#ff3d5a";

  return (
    <div style={{
      width: "100%", minHeight: "100vh", background: bg,
      fontFamily: "'Segoe UI', system-ui, sans-serif", color: textPri,
      display: "flex", justifyContent: "center", alignItems: "center", padding: 20,
    }}>
      <div style={{
        width: 380, background: card, borderRadius: 24,
        border: `1px solid ${border}`, overflow: "hidden",
        boxShadow: "0 20px 60px rgba(0,0,0,0.5)",
      }}>
        {/* Header */}
        <div style={{
          padding: "16px 20px", display: "flex", justifyContent: "space-between",
          alignItems: "center", borderBottom: `1px solid ${border}`,
        }}>
          <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
            <span style={{
              fontFamily: "'Courier New', monospace", fontWeight: 700,
              fontSize: 18, letterSpacing: -0.5,
            }}>
              Flux<span style={{ color: cyan }}>VPN</span>
            </span>
            <span style={{
              fontSize: 9, padding: "2px 6px", borderRadius: 4,
              background: "rgba(0,229,255,0.1)", color: cyan, fontWeight: 600,
            }}>v1.5.7</span>
          </div>
          <div style={{ display: "flex", gap: 4 }}>
            {["main", "nodes", "dashboard", "expert", "settings"].map(v => (
              <button key={v} onClick={() => setView(v)}
                style={{
                  background: view === v ? "rgba(0,229,255,0.1)" : "transparent",
                  border: `1px solid ${view === v ? "rgba(0,229,255,0.2)" : "transparent"}`,
                  color: view === v ? cyan : textSec, borderRadius: 8,
                  padding: "6px 12px", fontSize: 11, fontWeight: 600, cursor: "pointer",
                  transition: "all 0.2s",
                }}>
                {v === "main" ? "⚡" : v === "nodes" ? "🌐" : v === "dashboard" ? "📊" : v === "expert" ? "🔧" : "⚙️"}
              </button>
            ))}
            <button onClick={() => alert("Atajos:\nC: Conectar\nN: Nodos\nS: Config\nM: Main\nEsc: Main")}
              style={{
                background: "transparent", border: `1px solid transparent`,
                color: textSec, borderRadius: 8, padding: "6px 8px", fontSize: 11,
                cursor: "pointer", transition: "all 0.2s",
              }}>
              ?
            </button>
          </div>
        </div>

        {/* Main View */}
        {view === "main" && (
          <div style={{ padding: "30px 20px", textAlign: "center" }}>
            {/* Connection Orb */}
            <div style={{
              width: 160, height: 160, borderRadius: "50%", margin: "0 auto 24px",
              background: connected
                ? `radial-gradient(circle at 40% 40%, ${cyan}22, ${cyan}08)`
                : connecting
                ? `radial-gradient(circle at 40% 40%, ${blue}22, ${blue}08)`
                : `radial-gradient(circle at 40% 40%, #ffffff08, #ffffff03)`,
              border: `2px solid ${connected ? cyan + "40" : connecting ? blue + "40" : "#ffffff10"}`,
              display: "flex", flexDirection: "column", justifyContent: "center",
              alignItems: "center", cursor: "pointer", transition: "all 0.5s ease",
              boxShadow: connected ? `0 0 40px ${cyan}15` : "none",
              animation: connecting ? "pulse-connect 1.5s ease-in-out infinite" : "none",
            }} onClick={handleConnect}>
              <div style={{
                width: 48, height: 48, borderRadius: 12,
                background: connected
                  ? `linear-gradient(135deg, ${cyan}, ${blue})`
                  : connecting ? `linear-gradient(135deg, ${blue}88, ${cyan}44)` : "#ffffff10",
                display: "flex", alignItems: "center", justifyContent: "center",
                fontSize: 22, marginBottom: 8, transition: "all 0.5s",
              }}>
                {connected ? "🔒" : connecting ? "⟳" : "⏻"}
              </div>
              <span style={{
                fontSize: 13, fontWeight: 700, letterSpacing: 1,
                color: connected ? cyan : connecting ? blue : textSec,
              }}>
                {connected ? "PROTEGIDO" : connecting ? "CONECTANDO..." : "CONECTAR"}
              </span>
            </div>

            {/* Selected Node */}
            <button onClick={() => setView("nodes")}
              style={{
                width: "100%", padding: "12px 16px", borderRadius: 12,
                background: "rgba(255,255,255,0.03)", border: `1px solid ${border}`,
                color: textPri, display: "flex", justifyContent: "space-between",
                alignItems: "center", cursor: "pointer", marginBottom: 16,
                transition: "all 0.2s", fontSize: 14,
              }}>
              <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
                <span style={{ fontSize: 20 }}>🌐</span>
                <div style={{ textAlign: "left" }}>
                  <div style={{ fontWeight: 600, fontSize: 13 }}>{selectedNode ? (selectedNode.name || selectedNode.id) : "Seleccionar Peer"}</div>
                  <div style={{ fontSize: 11, color: textSec }}>{selectedNode ? selectedNode.addr : "Ningún peer seleccionado"}</div>
                </div>
              </div>
              <div style={{ textAlign: "right" }}>
                <div style={{ fontSize: 12, color: cyan, fontFamily: "monospace" }}>
                  {selectedNode ? (selectedNode.latency || 0) : 0}ms
                </div>
                <div style={{ fontSize: 10, color: textSec }}>▾ Cambiar</div>
              </div>
            </button>

            {/* Stats */}
            {connected && (
              <div style={{ animation: "fadeIn 0.5s ease-out" }}>
                <div style={{
                  display: "flex", gap: 10, marginBottom: 12,
                }}>
                  <div style={{
                    flex: 1, padding: "12px", borderRadius: 12,
                    background: "rgba(0,229,255,0.04)", border: `1px solid ${border}`,
                  }}>
                    <div style={{ fontSize: 10, color: textSec, marginBottom: 4 }}>↓ DESCARGA</div>
                    <div style={{ fontSize: 20, fontWeight: 700, fontFamily: "monospace", color: cyan }}>
                      {status ? (status.downMbps || 0).toFixed(1) : 0}
                    </div>
                    <div style={{ fontSize: 10, color: textSec }}>Mbps</div>
                  </div>
                  <div style={{
                    flex: 1, padding: "12px", borderRadius: 12,
                    background: "rgba(0,229,255,0.04)", border: `1px solid ${border}`,
                  }}>
                    <div style={{ fontSize: 10, color: textSec, marginBottom: 4 }}>↑ SUBIDA</div>
                    <div style={{ fontSize: 20, fontWeight: 700, fontFamily: "monospace", color: blue }}>
                      {status ? (status.upMbps || 0).toFixed(1) : 0}
                    </div>
                    <div style={{ fontSize: 10, color: textSec }}>Mbps</div>
                  </div>
                </div>

                {/* Timer */}
                <div style={{
                  padding: "10px", borderRadius: 12,
                  background: "rgba(0,229,255,0.04)", border: `1px solid ${border}`,
                  marginBottom: 12,
                }}>
                  <div style={{ fontSize: 10, color: textSec, marginBottom: 2 }}>TIEMPO CONECTADO</div>
                  <div style={{ fontSize: 24, fontWeight: 700, fontFamily: "monospace", letterSpacing: 2 }}>
                    {status ? formatTime(status.connectedTime || 0) : "00:00:00"}
                  </div>
                </div>

                {/* PQC Status */}
                <div style={{
                  padding: "10px 14px", borderRadius: 12, display: "flex",
                  justifyContent: "space-between", alignItems: "center",
                  background: pqcPulse ? "rgba(0,229,255,0.08)" : "rgba(0,229,255,0.02)",
                  border: `1px solid ${pqcPulse ? cyan + "30" : border}`,
                  transition: "all 0.5s",
                }}>
                  <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
                    <div style={{
                      width: 8, height: 8, borderRadius: "50%",
                      background: pqcPulse ? cyan : green,
                      boxShadow: pqcPulse ? `0 0 8px ${cyan}` : `0 0 4px ${green}44`,
                      transition: "all 0.3s",
                    }} />
                    <span style={{ fontSize: 11, fontWeight: 600, color: pqcPulse ? cyan : textPri }}>
                      {pqcPulse ? "🔐 Re-firma PQC ML-DSA-65" : "🛡️ HMAC activo"}
                    </span>
                  </div>
                  <span style={{ fontSize: 10, color: textSec, fontFamily: "monospace" }}>
                    FIPS-204
                  </span>
                </div>
              </div>
            )}
          </div>
        )}

        {/* Nodes View */}
        {view === "nodes" && (
          <div style={{ padding: "16px 20px" }}>
            <div style={{ fontSize: 14, fontWeight: 700, marginBottom: 12 }}>
              Peers P2P Activos
            </div>
            <div style={{ fontSize: 11, color: textSec, marginBottom: 16 }}>
              Peers conectados en la red IPv7. El gradiente optimiza rutas automáticamente.
            </div>
            {peers.length > 0 ? peers.map(peer => (
              <button key={peer.id} onClick={() => setSelectedNode(peer)}
                style={{
                  width: "100%", padding: "12px 14px", borderRadius: 12, marginBottom: 8,
                  background: selectedNode?.id === peer.id ? "rgba(0,229,255,0.06)" : "rgba(255,255,255,0.02)",
                  border: `1px solid ${selectedNode?.id === peer.id ? cyan + "30" : border}`,
                  color: textPri, display: "flex", justifyContent: "space-between",
                  alignItems: "center", cursor: "pointer", transition: "all 0.2s",
                }}>
                <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
                  <span style={{ fontSize: 22 }}>🌐</span>
                  <div style={{ textAlign: "left" }}>
                    <div style={{ fontWeight: 600, fontSize: 13 }}>{peer.name || peer.id}</div>
                    <div style={{ fontSize: 11, color: textSec }}>{peer.addr}</div>
                  </div>
                </div>
                <div style={{ textAlign: "right" }}>
                  <div style={{
                    fontSize: 12, fontFamily: "monospace", fontWeight: 600,
                    color: peer.latency < 20 ? green : peer.latency < 50 ? cyan : peer.latency < 150 ? "#ffab00" : red,
                  }}>{peer.latency || 0}ms</div>
                  <div style={{ fontSize: 10, color: textSec }}>
                    {peer.load || 0}%
                  </div>
                </div>
              </button>
            )) : (
              <div style={{ textAlign: "center", color: textSec, padding: 20 }}>
                No hay peers conectados. Inicia el daemon como master para crear red P2P.
              </div>
            )}
          </div>
        )}

        {/* Dashboard View */}
        {view === "dashboard" && (
          <div style={{ padding: "16px 20px" }}>
            <div style={{ fontSize: 14, fontWeight: 700, marginBottom: 12 }}>
              Dashboard de Telemetría
            </div>
            <div style={{ fontSize: 11, color: textSec, marginBottom: 16 }}>
              Métricas en tiempo real de la red IPv7.
            </div>
            {metricsHistory.length > 0 && (
              <div style={{ height: 200 }}>
                <Line
                  data={{
                    labels: metricsHistory.map((_, i) => `${i * 2}s`),
                    datasets: [{
                      label: "Latencia (ms)",
                      data: metricsHistory.map(m => m.latency || 0),
                      borderColor: cyan,
                      backgroundColor: cyan + "20",
                    }, {
                      label: "Descarga (Mbps)",
                      data: metricsHistory.map(m => m.downMbps || 0),
                      borderColor: blue,
                      backgroundColor: blue + "20",
                    }]
                  }}
                  options={{
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                      y: { beginAtZero: true }
                    }
                  }}
                />
              </div>
            )}
          </div>
        )}

        {/* Expert View */}
        {view === "expert" && (
          <div style={{ padding: "16px 20px" }}>
            <div style={{ fontSize: 14, fontWeight: 700, marginBottom: 12 }}>
              Modo Experto
            </div>
            <div style={{ fontSize: 11, color: textSec, marginBottom: 16 }}>
              Detalles técnicos avanzados de la red IPv7.
            </div>
            {status && (
              <div style={{ spaceY: 12 }}>
                <div style={{ padding: "12px", borderRadius: 12, background: "rgba(255,255,255,0.02)", border: `1px solid ${border}` }}>
                  <div style={{ fontSize: 12, fontWeight: 600, marginBottom: 4 }}>Estado del Daemon</div>
                  <div style={{ fontSize: 11, color: textSec }}>Puerto: {status.apiPort || 7781} | Peers: {peers.length} | Conectado: {connected ? "Sí" : "No"}</div>
                </div>
                <div style={{ padding: "12px", borderRadius: 12, background: "rgba(255,255,255,0.02)", border: `1px solid ${border}` }}>
                  <div style={{ fontSize: 12, fontWeight: 600, marginBottom: 4 }}>Métricas de Red</div>
                  <div style={{ fontSize: 11, color: textSec }}>
                    Latencia: {status.latency || 0}ms | Pérdida: {status.loss || 0}% | Jitter: {status.jitter || 0}ms
                  </div>
                </div>
                <div style={{ padding: "12px", borderRadius: 12, background: "rgba(255,255,255,0.02)", border: `1px solid ${border}` }}>
                  <div style={{ fontSize: 12, fontWeight: 600, marginBottom: 4 }}>Criptografía</div>
                  <div style={{ fontSize: 11, color: textSec }}>
                    PQC: ML-DSA-65 | HMAC: SHA-256 | ECC: Curve25519
                  </div>
                </div>
              </div>
            )}
          </div>
        )}

        {/* Settings View */}
        {view === "settings" && (
          <div style={{ padding: "16px 20px" }}>
            <div style={{ fontSize: 14, fontWeight: 700, marginBottom: 16 }}>Configuración</div>
            {[
              { key: "killSwitch", label: "Kill Switch", desc: "Cortar internet si FluxVPN se desconecta" },
              { key: "autoConnect", label: "Auto-conectar", desc: "Conectar al iniciar el sistema" },
              { key: "starlinkMode", label: "Modo Starlink", desc: "Optimizar para redes satelitales LEO" },
              { key: "pqcStrict", label: "PQC Estricto", desc: "Solo conexiones con nodos ML-DSA-65" },
              { key: "ghostUpdater", label: "GhostUpdater", desc: "Actualizaciones automáticas con SHA-256" },
            ].map((s, i) => (
              <div key={i} style={{
                display: "flex", justifyContent: "space-between", alignItems: "center",
                padding: "14px 0", borderBottom: `1px solid ${border}`,
              }}>
                <div>
                  <div style={{ fontSize: 13, fontWeight: 600 }}>{s.label}</div>
                  <div style={{ fontSize: 11, color: textSec, marginTop: 2 }}>{s.desc}</div>
                </div>
                <div style={{
                  width: 40, height: 22, borderRadius: 11, cursor: "pointer",
                  background: settings[s.key] ? cyan : "#ffffff15",
                  display: "flex", alignItems: "center",
                  justifyContent: settings[s.key] ? "flex-end" : "flex-start", padding: 2,
                  transition: "all 0.2s",
                }} onClick={() => {
                  const newSettings = { ...settings, [s.key]: !settings[s.key] };
                  setSettings(newSettings);
                  axios.post(`${API_BASE}/config`, newSettings).catch(err => console.error("Error updating config:", err));
                }}>
                  <div style={{
                    width: 18, height: 18, borderRadius: "50%",
                    background: settings[s.key] ? bg : "#555",
                  }} />
                </div>
              </div>
            ))}

            <div style={{
              marginTop: 20, padding: "12px", borderRadius: 12,
              background: "rgba(0,229,255,0.03)", border: `1px solid ${border}`,
              fontSize: 11, color: textSec, lineHeight: 1.6,
            }}>
              <div style={{ fontWeight: 600, color: textPri, marginBottom: 4 }}>Protocolo</div>
              IPv7-IEU v1.5.7 · ML-DSA-65 FIPS-204<br/>
              Enrutamiento: −∇ ln(L) gradiente logarítmico<br/>
              DHT: Kademlia MicroDHT · P2P puro
            </div>
          </div>
        )}

        {/* Bottom bar */}
        <div style={{
          padding: "12px 20px", borderTop: `1px solid ${border}`,
          display: "flex", justifyContent: "space-between", alignItems: "center",
        }}>
          <div style={{ display: "flex", alignItems: "center", gap: 6 }}>
            <div style={{
              width: 6, height: 6, borderRadius: "50%",
              background: connected ? green : red,
              boxShadow: connected ? `0 0 6px ${green}44` : "none",
            }} />
            <span style={{ fontSize: 11, color: textSec }}>
              {connected ? "Conectado · P2P" : "Desconectado"}
            </span>
          </div>
          <span style={{ fontSize: 10, color: textSec, fontFamily: "monospace" }}>
            did:ipv7:flux
          </span>
        </div>
      </div>

      <style>{`
        @keyframes pulse-connect {
          0%, 100% { opacity: 0.8; transform: scale(1); }
          50% { opacity: 1; transform: scale(1.03); }
        }
        @keyframes fadeIn {
          from { opacity: 0; transform: translateY(8px); }
          to { opacity: 1; transform: translateY(0); }
        }
      `}</style>
    </div>
  );
}