package protocol

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"
	"sync"
)

type IPv7Address struct {
	Region     float64
	Subnet     float64
	DeviceID   float64
	ResolvedIP float64
	SubPort    uint16
}

func NewIPv7(r, s, d float64) IPv7Address {
	if r <= 0 || s <= 0 || d <= 0 {
		return IPv7Address{Region: r, Subnet: s, DeviceID: d, ResolvedIP: 0}
	}
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%.0f:%.0f:%.0f", r, s, d)))
	val := float64(h.Sum64() & 0x1FFFFFFFFFFFFF)
	return IPv7Address{Region: r, Subnet: s, DeviceID: d, ResolvedIP: val, SubPort: 0}
}

func NewIPv7WithSubPort(r, s, d float64, subPort uint16) IPv7Address {
	addr := NewIPv7(r, s, d)
	addr.SubPort = subPort
	return addr
}

func (a IPv7Address) Equals(b IPv7Address) bool {
	return math.Abs(a.ResolvedIP-b.ResolvedIP) < 1e-9
}

func (a IPv7Address) SerializeHeader() []byte {
	buf := make([]byte, 10)
	binary.BigEndian.PutUint16(buf[0:2], uint16(a.Region))
	binary.BigEndian.PutUint16(buf[2:4], uint16(a.Subnet))
	binary.BigEndian.PutUint32(buf[4:8], uint32(a.DeviceID))
	binary.BigEndian.PutUint16(buf[8:10], a.SubPort)
	return buf
}

func ParseHeader(buf []byte) IPv7Address {
	if len(buf) < 10 {
		if len(buf) < 8 {
			return IPv7Address{}
		}
		r := float64(binary.BigEndian.Uint16(buf[0:2]))
		s := float64(binary.BigEndian.Uint16(buf[2:4]))
		d := float64(binary.BigEndian.Uint32(buf[4:8]))
		return NewIPv7(r, s, d)
	}
	r := float64(binary.BigEndian.Uint16(buf[0:2]))
	s := float64(binary.BigEndian.Uint16(buf[2:4]))
	d := float64(binary.BigEndian.Uint32(buf[4:8]))
	sp := binary.BigEndian.Uint16(buf[8:10])
	addr := NewIPv7(r, s, d)
	addr.SubPort = sp
	return addr
}

const HeaderSize = 10

const latencyRingSize = 8
const emaAlpha = 0.25
const trendPenalty = 0.15

type LatencyHistory struct {
	samples [latencyRingSize]float64
	head    int
	filled  bool
	ema     float64
}

func (h *LatencyHistory) Push(latency float64) {
	l := math.Max(latency, 1e-9)
	if !h.filled && h.samples[0] == 0 {
		h.ema = l
	} else {
		h.ema = emaAlpha*l + (1-emaAlpha)*h.ema
	}
	h.samples[h.head] = l
	h.head = (h.head + 1) % latencyRingSize
	if h.head == 0 {
		h.filled = true
	}
}

func (h *LatencyHistory) Trend() float64 {
	n := latencyRingSize
	if !h.filled {
		n = h.head
	}
	if n < 2 {
		return 0
	}
	var sumX, sumXI float64
	for i := 0; i < n; i++ {
		idx := (h.head - n + i + latencyRingSize) % latencyRingSize
		sumX += h.samples[idx]
		sumXI += float64(i) * h.samples[idx]
	}
	meanX := sumX / float64(n)
	meanI := float64(n-1) / 2
	var num, den float64
	for i := 0; i < n; i++ {
		idx := (h.head - n + i + latencyRingSize) % latencyRingSize
		num += (float64(i) - meanI) * (h.samples[idx] - meanX)
		den += (float64(i) - meanI) * (float64(i) - meanI)
	}
	_ = sumXI
	if den < 1e-12 {
		return 0
	}
	return num / den
}

type Node struct {
	Name      string
	Address   IPv7Address
	Latency   float64
	Neighbors []*Node
	Mu        sync.RWMutex

	latencyHistory       []LatencyHistory
	stochasticPredictors []*StochasticLatencyPredictor
}

func (n *Node) UpdateNeighborLatency(neighborIdx int, latency float64) {
	n.Mu.Lock()
	defer n.Mu.Unlock()
	if neighborIdx >= len(n.latencyHistory) {
		extra := make([]LatencyHistory, neighborIdx-len(n.latencyHistory)+1)
		n.latencyHistory = append(n.latencyHistory, extra...)
	}
	n.latencyHistory[neighborIdx].Push(latency)
	for len(n.stochasticPredictors) <= neighborIdx {
		n.stochasticPredictors = append(n.stochasticPredictors, NewStochasticPredictor(16))
	}
	n.stochasticPredictors[neighborIdx].Push(latency)
}

func (n *Node) NextHop() *Node {
	n.Mu.RLock()
	myLatency := math.Max(n.Latency, 1e-9)
	neighborsCopy := make([]*Node, len(n.Neighbors))
	copy(neighborsCopy, n.Neighbors)
	historyCopy := make([]LatencyHistory, len(n.latencyHistory))
	copy(historyCopy, n.latencyHistory)
	predCopy := make([]*StochasticLatencyPredictor, len(n.stochasticPredictors))
	copy(predCopy, n.stochasticPredictors)
	n.Mu.RUnlock()

	var bestNeighbor *Node
	maxScore := -math.MaxFloat64

	for i, neighbor := range neighborsCopy {
		neighbor.Mu.RLock()
		nl := math.Max(neighbor.Latency, 1e-9)
		neighbor.Mu.RUnlock()

		effectiveLatency := nl
		var trend float64
		if i < len(historyCopy) && historyCopy[i].ema > 1e-9 {
			effectiveLatency = historyCopy[i].ema
			trend = historyCopy[i].Trend()
		}

		baseGradient := math.Log(myLatency) - math.Log(effectiveLatency)
		degradationPenalty := trendPenalty * math.Tanh(trend/10.0)

		var stochasticPenalty float64
		if i < len(predCopy) && predCopy[i] != nil {
			risk := predCopy[i].IsHighRisk(5, effectiveLatency*2)
			stochasticPenalty = 0.3 * risk
		}

		score := baseGradient - degradationPenalty - stochasticPenalty
		if score > maxScore {
			maxScore = score
			bestNeighbor = neighbor
		}
	}
	return bestNeighbor
}

func (n *Node) NextHopWithPrediction(degradationThresholdMs float64) (*Node, float64, float64, float64) {
	best := n.NextHop()
	if best == nil {
		return nil, 0, 0, 0
	}
	n.Mu.RLock()
	defer n.Mu.RUnlock()
	for i, nb := range n.Neighbors {
		if nb != best {
			continue
		}
		if i < len(n.stochasticPredictors) && n.stochasticPredictors[i] != nil {
			mean, std := n.stochasticPredictors[i].Predict(5)
			risk := n.stochasticPredictors[i].IsHighRisk(5, degradationThresholdMs)
			return best, mean, std, risk
		}
		break
	}
	return best, 0, 0, 0
}
