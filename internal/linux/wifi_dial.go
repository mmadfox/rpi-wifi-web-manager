package linux

import (
	"fmt"
	"os/exec"
	"strings"
)

func DialWiFi(ssid string, password string, save bool, iface string) error {
	if len(ssid) == 0 {
		return fmt.Errorf("ssid is empty")
	}
	q := []string{
		"nmcli",
		"--colors", "no",
		"dev", "wifi",
		"connect", fmt.Sprintf("'%s'", ssid),
		"ifname", iface,
	}
	if len(password) > 0 {
		q = append(q, "password", fmt.Sprintf("'%s'", password))
	}
	if save {
		q = append(q, "--ask")
	}
	str := strings.Join(q, " ")
	cmd := exec.Command("sh", "-c", str)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, out)
	}
	return nil
}
