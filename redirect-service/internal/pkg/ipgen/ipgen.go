package ipgen

import (
	"math/rand"
	"net"
	"time"
)

type IPGenerator struct {
	rng *rand.Rand
}

var IpGenerator *IPGenerator

func Init() {
	IpGenerator = &IPGenerator{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// 已知的公网 IP 段（这些是真实的 ISP 地址段）
var publicIPRanges = []struct {
	network string
	mask    string
}{
	{"1.0.0.0", "255.0.0.0"},   // APNIC
	{"8.0.0.0", "255.0.0.0"},   // Level 3
	{"23.0.0.0", "255.0.0.0"},  // Akamai
	{"31.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"37.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"45.0.0.0", "255.0.0.0"},  // ARIN
	{"49.0.0.0", "255.0.0.0"},  // APNIC
	{"58.0.0.0", "255.0.0.0"},  // APNIC
	{"59.0.0.0", "255.0.0.0"},  // APNIC
	{"61.0.0.0", "255.0.0.0"},  // APNIC
	{"62.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"63.0.0.0", "255.0.0.0"},  // ARIN
	{"77.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"78.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"79.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"80.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"81.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"82.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"83.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"84.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"85.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"86.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"87.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"88.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"89.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"90.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"91.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"92.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"93.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"94.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"95.0.0.0", "255.0.0.0"},  // RIPE NCC
	{"108.0.0.0", "255.0.0.0"}, // ARIN
	{"173.0.0.0", "255.0.0.0"}, // ARIN
	{"174.0.0.0", "255.0.0.0"}, // ARIN
	{"184.0.0.0", "255.0.0.0"}, // ARIN
	{"199.0.0.0", "255.0.0.0"}, // ARIN
	{"204.0.0.0", "255.0.0.0"}, // ARIN
	{"205.0.0.0", "255.0.0.0"}, // ARIN
	{"206.0.0.0", "255.0.0.0"}, // ARIN
	{"207.0.0.0", "255.0.0.0"}, // ARIN
	{"208.0.0.0", "255.0.0.0"}, // ARIN
	{"209.0.0.0", "255.0.0.0"}, // ARIN
	{"216.0.0.0", "255.0.0.0"}, // ARIN
}

// 生成随机公网 IP
func (g *IPGenerator) GeneratePublicIP() string {
	// 随机选择一个 IP 段
	ipRange := publicIPRanges[g.rng.Intn(len(publicIPRanges))]

	// 解析网络地址
	ip := net.ParseIP(ipRange.network).To4()
	mask := net.IPMask(net.ParseIP(ipRange.mask).To4())

	// 生成随机 IP（在网段内）
	network := &net.IPNet{
		IP:   ip.Mask(mask),
		Mask: mask,
	}

	return g.generateIPInRange(network)
}

// 在指定网段内生成随机 IP
func (g *IPGenerator) generateIPInRange(network *net.IPNet) string {
	ip := make(net.IP, len(network.IP))
	copy(ip, network.IP)

	// 随机化主机部分
	for i := 0; i < len(ip); i++ {
		if network.Mask[i] == 0 {
			ip[i] = byte(g.rng.Intn(256))
		} else if network.Mask[i] < 255 {
			// 部分掩码，需要特殊处理（这里简化处理）
			ip[i] = byte(g.rng.Intn(256)) & network.Mask[i]
		}
	}

	return ip.String()
}
