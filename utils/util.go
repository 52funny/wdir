package utils

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// ConvertSize is Converting bytes into individual units
func ConvertSize(n int64) (ans string) {
	if n < 1024 {
		// B
		ans = fmt.Sprintf("%dB", n)
	} else if n < 1024*1024 {
		// KB
		ans = fmt.Sprintf("%.2fKB", float64(n)/1024)
	} else if n < 1024*1024*1024 {
		// M
		ans = fmt.Sprintf("%.2fM", float64(n)/1024/1024)
	} else {
		// G
		ans = fmt.Sprintf("%.2fG", float64(n)/1024/1024/1024)
	}
	return
}

// PathExists is determine if a folder exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// Check for hidden directories in the path
func PathHidden(path string) bool {
	fileList := strings.Split(path, "/")
	for _, d := range fileList {
		if len(d) == 0 {
			continue
		}
		if FileHidden(d) {
			return true
		}
	}
	return false
}

func isLocalIPv6(ip net.IP) bool {
	// remove link-local and private addresses
	if ip.IsLinkLocalUnicast() || ip.IsPrivate() {
		return false
	}
	return ip.IsGlobalUnicast()
}

func GetNetIPv6Address() []string {
	netS := make([]string, 0)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() == nil && ipNet.IP.To16() != nil && isLocalIPv6(ipNet.IP) {
				netS = append(netS, fmt.Sprintf("[%s]", ipNet.IP.String()))
			}
		}
	}
	return netS
}

// Get all network interface
func GetNetIPv4Address() []string {
	netS := make([]string, 0)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				netS = append(netS, ipNet.IP.String())
			}
		}
	}
	return netS
}
