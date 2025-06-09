package sub

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
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

	stop := make(chan struct{})
	arpReplies := make(chan *ARPReply)
	go readARP(context.Background(), handle, iface, stop, arpReplies, userService)
	for {
		found := make(map[string]bool)

		go func() {
			for result := range arpReplies {
				if !found[result.IP.String()] {
					found[result.IP.String()] = true
					cmd.Usage(fmt.Sprintf("%v :: IP %v is at %v\n", time.Now().In(time.Local).Format("2006-01-02 15:04:05"), result.IP, result.MAC))
				}
			}
		}()

		if err := writeARP(handle, iface, addr); err != nil {
			cmd.Usage(fmt.Sprintf("%v :: error writing packets on %v: %v\n", time.Now().In(time.Local).Format("2006-01-02 15:04:05"), iface.Name, err))
			return err
		}

		time.Sleep(3 * time.Second)
	}

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
				_, err := userService.Register(ctx, net.IP(arp.SourceProtAddress).String(), "", net.HardwareAddr(arp.SourceHwAddress).String())
				if err != nil && err.Error() != "ip already registred" {
					cmd.Logger.Printf("Failed to register IP %v: %v\n", net.IP(arp.SourceProtAddress), err)
				}
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
