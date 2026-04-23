import { useState, useEffect } from "react";

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
  const [selectedNode, setSelectedNode] = useState(NODES[5]);
  const [showNodes, setShowNodes] = useState(false);
  const [elapsed, setElapsed] = useState(0);
  const [stats, setStats] = useState({ down: 0, up: 0 });
  const [view, setView] = useState("main");
  const [pqcPulse, setPqcPulse] = useState(false);

  useEffect(() => {
    if (!connected) { setElapsed(0); return; }
    const t = setInterval(() => setElapsed(e => e + 1), 1000);
    return () => clearInterval(t);
  }, [connected]);

  useEffect(() => {
    if (!connected) { setStats({ down: 0, up: 0 }); return; }
    const t = setInterval(() => {
      setStats({
        down: +(Math.random() * 80 + 20).toFixed(1),
        up: +(Math.random() * 30 + 5).toFixed(1),
      });
    }, 1500);
    return () => clearInterval(t);
  }, [connected]);

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

  const formatTime = (s) => {
    const h = Math.floor(s / 3600).toString().padStart(2, "0");
    const m = Math.floor((s % 3600) / 60).toString().padStart(2, "0");
    const sec = (s % 60).toString().padStart(2, "0");
    return `${h}:${m}:${sec}`;
  };

  const handleConnect = () => {
    if (connected) { setConnected(false); setConnecting(false); return; }
    setConnecting(true);
    setTimeout(() => { setConnecting(false); setConnected(true); }, 2200);
  };

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
            {["main", "nodes", "settings"].map(v => (
              <button key={v} onClick={() => { setView(v); setShowNodes(v === "nodes"); }}
                style={{
                  background: view === v ? "rgba(0,229,255,0.1)" : "transparent",
                  border: `1px solid ${view === v ? "rgba(0,229,255,0.2)" : "transparent"}`,
                  color: view === v ? cyan : textSec, borderRadius: 8,
                  padding: "6px 12px", fontSize: 11, fontWeight: 600, cursor: "pointer",
                  transition: "all 0.2s",
                }}>
                {v === "main" ? "⚡" : v === "nodes" ? "🌐" : "⚙️"}
              </button>
            ))}
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
                <span style={{ fontSize: 20 }}>{selectedNode.flag}</span>
                <div style={{ textAlign: "left" }}>
                  <div style={{ fontWeight: 600, fontSize: 13 }}>{selectedNode.name}</div>
                  <div style={{ fontSize: 11, color: textSec }}>{selectedNode.country}</div>
                </div>
              </div>
              <div style={{ textAlign: "right" }}>
                <div style={{ fontSize: 12, color: cyan, fontFamily: "monospace" }}>
                  {selectedNode.latency}ms
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
                      {stats.down}
                    </div>
                    <div style={{ fontSize: 10, color: textSec }}>Mbps</div>
                  </div>
                  <div style={{
                    flex: 1, padding: "12px", borderRadius: 12,
                    background: "rgba(0,229,255,0.04)", border: `1px solid ${border}`,
                  }}>
                    <div style={{ fontSize: 10, color: textSec, marginBottom: 4 }}>↑ SUBIDA</div>
                    <div style={{ fontSize: 20, fontWeight: 700, fontFamily: "monospace", color: blue }}>
                      {stats.up}
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
                    {formatTime(elapsed)}
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
              Nodos Disponibles
            </div>
            <div style={{ fontSize: 11, color: textSec, marginBottom: 16 }}>
              Selecciona un nodo. El gradiente optimizará la ruta automáticamente.
            </div>
            {NODES.map(node => (
              <button key={node.id} onClick={() => { setSelectedNode(node); setView("main"); }}
                style={{
                  width: "100%", padding: "12px 14px", borderRadius: 12, marginBottom: 8,
                  background: selectedNode.id === node.id ? "rgba(0,229,255,0.06)" : "rgba(255,255,255,0.02)",
                  border: `1px solid ${selectedNode.id === node.id ? cyan + "30" : border}`,
                  color: textPri, display: "flex", justifyContent: "space-between",
                  alignItems: "center", cursor: "pointer", transition: "all 0.2s",
                }}>
                <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
                  <span style={{ fontSize: 22 }}>{node.flag}</span>
                  <div style={{ textAlign: "left" }}>
                    <div style={{ fontWeight: 600, fontSize: 13 }}>{node.name}</div>
                    <div style={{ fontSize: 11, color: textSec }}>{node.country}</div>
                  </div>
                </div>
                <div style={{ textAlign: "right" }}>
                  <div style={{
                    fontSize: 12, fontFamily: "monospace", fontWeight: 600,
                    color: node.latency < 20 ? green : node.latency < 50 ? cyan : node.latency < 150 ? "#ffab00" : red,
                  }}>{node.latency}ms</div>
                  <div style={{ fontSize: 10, color: textSec }}>
                    <span style={{
                      display: "inline-block", width: 40, height: 4, borderRadius: 2,
                      background: "#ffffff10", position: "relative", overflow: "hidden",
                    }}>
                      <span style={{
                        position: "absolute", left: 0, top: 0, height: "100%",
                        width: `${node.load}%`, borderRadius: 2,
                        background: node.load > 70 ? red : node.load > 50 ? "#ffab00" : green,
                      }} />
                    </span>
                  </div>
                </div>
              </button>
            ))}
          </div>
        )}

        {/* Settings View */}
        {view === "settings" && (
          <div style={{ padding: "16px 20px" }}>
            <div style={{ fontSize: 14, fontWeight: 700, marginBottom: 16 }}>Configuración</div>
            {[
              { label: "Kill Switch", desc: "Cortar internet si FluxVPN se desconecta", on: true },
              { label: "Auto-conectar", desc: "Conectar al iniciar el sistema", on: false },
              { label: "Modo Starlink", desc: "Optimizar para redes satelitales LEO", on: true },
              { label: "PQC Estricto", desc: "Solo conexiones con nodos ML-DSA-65", on: true },
              { label: "GhostUpdater", desc: "Actualizaciones automáticas con SHA-256", on: true },
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
                  background: s.on ? cyan : "#ffffff15",
                  display: "flex", alignItems: "center",
                  justifyContent: s.on ? "flex-end" : "flex-start", padding: 2,
                  transition: "all 0.2s",
                }}>
                  <div style={{
                    width: 18, height: 18, borderRadius: "50%",
                    background: s.on ? bg : "#555",
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
