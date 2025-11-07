package smartapi

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func publicIP() string {
	services := []string{
		"https://api.ipify.org",
		"https://ifconfig.me",
		"https://icanhazip.com",
	}

	client := &http.Client{Timeout: 5 * time.Second}
	for _, url := range services {
		resp, err := client.Get(url)
		if err != nil {
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err == nil {
			ip := strings.TrimSpace(string(body))
			if ip != "" {
				return ip
			}
		}
	}

	return "127.0.0.1"
}

func macID() string {
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		if i.Flags&net.FlagLoopback != 0 || len(i.HardwareAddr) == 0 {
			continue
		}
		n := strings.ToLower(i.Name)
		if strings.Contains(n, "docker") || strings.Contains(n, "veth") || strings.Contains(n, "virtual") {
			continue
		}
		return strings.ToUpper(strings.ReplaceAll(i.HardwareAddr.String(), ":", "-"))
	}

	// fallback: random 12-byte hex (like MAC)
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "UNKNOWN"
	}
	return strings.ToUpper(hex.EncodeToString(b))
}

func localIP() string {

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
