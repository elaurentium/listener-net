package sub

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/elaurentium/listener-net/cmd"
	"github.com/elaurentium/listener-net/internal/domain/service"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type ARPReply struct {
	IP  net.IP
	MAC net.HardwareAddr
}

func Interfaces(userService *service.UserService) {
	ifaces, err := net.Interfaces()

	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, iface := range ifaces {
		wg.Add(1)
		go func(iface net.Interface) {
			defer wg.Done()
			if err := scan(&iface, userService); err != nil {
				cmd.Usage(fmt.Sprintf("interface %v: %v\n", iface.Name, err))
			}
		}(iface)
	}

	wg.Wait()
}

func scan(iface *net.Interface, userService *service.UserService) error {
	var addr *net.IPNet

	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				addr = &net.IPNet{
					IP:   ip4,
					Mask: ipnet.Mask[len(ipnet.Mask)-4:],
				}
				break
			}
		}
	}

	if addr == nil {
		return errors.New("no good IP network found")
	} else if addr.IP[0] == 127 {
		return errors.New("skipping localhost")
	} else if addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
		return errors.New("mask means network is too large")
	}
	cmd.Usage(fmt.Sprintf("Using network range %v for interface %v\n", addr, iface.Name))

	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return err
	}
	defer handle.Close()

	if err := handle.SetBPFFilter("arp"); err != nil {
		return fmt.Errorf("failed to set BPF filter: %w", err)
	}

	stop := make(chan struct{})
	arpReplies := make(chan *ARPReply)
	go readARP(context.Background(), handle, iface, stop, arpReplies, userService)
	connectedIPs, err := getARPCacheIPs()
    if err != nil {
        cmd.Usage(fmt.Sprintf("Warning: failed to read ARP cache: %v\n", err))
    }

    // Track discovered devices
    found := make(map[string]bool)

    // Process ARP replies
    go func() {
        for result := range arpReplies {
            ipStr := result.IP.String()
            if !found[ipStr] {
                found[ipStr] = true
                cmd.Usage(fmt.Sprintf("%v :: IP %v is at %v\n", time.Now().In(time.Local).Format("2006-01-02 15:04:05"), result.IP, result.MAC))
            }
        }
    }()

    // Periodically send ARP requests to known IPs
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // Refresh ARP cache
            currentIPs, err := getARPCacheIPs()
            if err != nil {
                cmd.Usage(fmt.Sprintf("Warning: failed to refresh ARP cache: %v\n", err))
            } else {
                connectedIPs = currentIPs
            }

            // Send ARP requests only to known connected IPs
            for _, ip := range connectedIPs {
                if !found[ip.String()] {
                    if err := writeARP(handle, iface, addr); err != nil {
                        cmd.Usage(fmt.Sprintf("%v :: error sending ARP request to %v: %v\n", time.Now().In(time.Local).Format("2006-01-02 15:04:05"), ip, err))
                    }
                }
            }
        case <-stop:
            return nil
        }
    }

}

// getARPCacheIPs reads the local ARP cache to get IPs of connected devices
func getARPCacheIPs() ([]net.IP, error) {
    var ips []net.IP

	switch runtime.GOOS {
	case "linux":
		data, err := os.ReadFile("/proc/net/arp")
		if err != nil {
			return nil, fmt.Errorf("failed to read ARP cache on Linux: %w", err)
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines[1:] {
			fields := strings.Fields(line)
			if len(fields) >= 1 {
				ip := net.ParseIP(fields[0])
				if ip != nil && ip.To4() != nil {
					ips = append(ips, ip)
				}
			}
		}

	case "darwin":
		out, err := exec.Command("arp", "-a").Output()
		if err != nil {
			return nil, fmt.Errorf("failed to read ARP cache on macOS: %w", err)
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			// Format: ? (192.168.1.1) at 0:11:22:33:44:55 on en0 ifscope [ethernet]
			if strings.Contains(line, " at ") {
				start := strings.Index(line, "(")
				end := strings.Index(line, ")")
				if start >= 0 && end > start {
					ipStr := line[start+1 : end]
					ip := net.ParseIP(ipStr)
					if ip != nil && ip.To4() != nil {
						ips = append(ips, ip)
					}
				}
			}
		}

	case "windows":
		out, err := exec.Command("arp", "-a").Output()
		if err != nil {
			return nil, fmt.Errorf("failed to read ARP cache on Windows: %w", err)
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 1 {
				ip := net.ParseIP(fields[0])
				if ip != nil && ip.To4() != nil {
					ips = append(ips, ip)
				}
			}
		}

	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	return ips, nil
}

func readARP(ctx context.Context, handle *pcap.Handle, iface *net.Interface, stop chan struct{}, replies chan *ARPReply, userService *service.UserService) {
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()

	for {
		select {
		case <-stop:
			close(replies)
			return
		case packet := <-in:
			if packet == nil {
				continue
			}
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}
			arp := arpLayer.(*layers.ARP)
			if arp.Operation != layers.ARPReply || bytes.Equal([]byte(iface.HardwareAddr), arp.SourceHwAddress) {
				continue
			}
			replies <- &ARPReply{
				IP:  net.IP(arp.SourceProtAddress),
				MAC: net.HardwareAddr(arp.SourceHwAddress),
			}

			if userService != nil {
				hostname, err := os.Hostname()
				if err != nil {
					cmd.Logger.Printf("could not get hostname: %v", err)
					continue
				}

				userService.Register(ctx, net.IP(arp.SourceProtAddress).String(), hostname, net.HardwareAddr(arp.SourceHwAddress).String())
			}
		}
	}
}

func writeARP(handle *pcap.Handle, iface *net.Interface, addr *net.IPNet) error {
	eth := layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(iface.HardwareAddr),
		SourceProtAddress: []byte(addr.IP),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	for _, ip := range ips(addr) {
		arp.DstProtAddress = []byte(ip)
		gopacket.SerializeLayers(buf, opts, &eth, &arp)
		if err := handle.WritePacketData(buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

func ips(n *net.IPNet) (out []net.IP) {
	num := binary.BigEndian.Uint32([]byte(n.IP))
	mask := binary.BigEndian.Uint32([]byte(n.Mask))
	network := num & mask
	broadcast := network | ^mask
	for network++; network < broadcast; network++ {
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[:], network)
		out = append(out, net.IP(buf[:]))
	}
	return
}
