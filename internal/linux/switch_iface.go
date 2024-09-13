package linux

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func SwitchInterface(interfaceName string, metric int) error {
	if len(interfaceName) == 0 || len(interfaceName) > 12 {
		return fmt.Errorf("invalid iface name")
	}
	if metric <= 0 {
		metric = 50
	}
	newGatewayIP, err := getDevIP(interfaceName)
	if err != nil {
		return fmt.Errorf("failed to retrieve gateway IP: %v", err)
	}
	delRouteCmd := exec.Command("ip", "route", "del", "default")
	if err := delRouteCmd.Run(); err != nil {
		return fmt.Errorf("failed to delete default route: %v", err)
	}
	addRouteCmd := exec.Command("ip", "route", "add", "default", "via", newGatewayIP, "dev", interfaceName, "metric", fmt.Sprintf("%d", metric))
	if err := addRouteCmd.Run(); err != nil {
		return fmt.Errorf("failed to add default route: %v", err)
	}
	return nil
}

func getDevIP(interfaceName string) (string, error) {
	cmd := exec.Command("ip", "-4", "addr", "show", "dev", interfaceName)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run ip command: %v", err)
	}
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "inet") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				ip := fields[1]
				ip = strings.Split(ip, "/")[0]
				octets := strings.Split(ip, ".")
				if len(octets) == 4 {
					return fmt.Sprintf("%s.%s.%s.1", octets[0], octets[1], octets[2]), nil
				}
			}
		}
	}
	return "", fmt.Errorf("no IP address found for interface %s", interfaceName)
}
