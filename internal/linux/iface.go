package linux

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

type IfaceInfo struct {
	Name    string     `json:"name"`
	Addrs   []net.Addr `json:"addrs"`
	Default bool       `json:"default"`
	Active  bool       `json:"active"`
}

func GetInterfaces() ([]IfaceInfo, error) {
	result := make([]IfaceInfo, 0, 2)
	defaultIface := getDefaultGatewayInterface()
	interfaces, err := getNetworkInterfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		ifaceInfo := IfaceInfo{
			Name:   iface.Name,
			Addrs:  make([]net.Addr, 0),
			Active: iface.Flags&net.FlagUp != 0,
		}
		if iface.Name == defaultIface {
			ifaceInfo.Default = true
		}
		addrs, _ := iface.Addrs()
		ifaceInfo.Addrs = addrs
		result = append(result, ifaceInfo)
	}
	return result, nil
}

func getDefaultGatewayInterface() string {
	out, err := exec.Command("ip", "route").Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "default via") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				return fields[4]
			}
		}
	}
	return ""
}

func getNetworkInterfaces() ([]net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %v", err)
	}
	return interfaces, nil
}
