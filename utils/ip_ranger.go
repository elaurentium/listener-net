package utils

import (
	"fmt"
	"net"
	"sync"
	"time"
)


func localRangeIP() ([]string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return nil, fmt.Errorf("failed to get local IP addresses: %w", err)
	}

	var localIP net.IP

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				localIP = ipNet.IP.To4()
				break
			}
 		}
	} 

	if localIP == nil {
		return nil, fmt.Errorf("no valid local IP address found")
	}

	ipBase := fmt.Sprintf("%d.%d.%d.", localIP[0], localIP[1], localIP[2])
	ips := []string{}

	for i := 0; i < 256; i++ {
		ips = append(ips, fmt.Sprintf("%s.%d", ipBase, i))
	}

	return ips, nil
}

func isHostUp(ip string) bool {
	conn, err := net.DialTimeout("tcp", ip+":80", 500*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func DetectIPRange() {
	ips, err := localRangeIP()
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	activeIPs := []string{}

	for _, ip := range ips {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			if isHostUp(ip) {
				mu.Lock()
				activeIPs = append(activeIPs, ip)
				mu.Unlock()
			}
		}(ip)
	}

	wg.Wait()

	fmt.Println("Dispositivos ativos:")
	for _, ip := range activeIPs {
		fmt.Println(" -", ip)
	}
}