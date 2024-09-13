package linux

import (
	"fmt"
	"os/exec"
	"strings"
)

func CloseWiFi() error {
	q := []string{
		"nmcli",
		"dev", "disconnect", "wlan0",
	}
	str := strings.Join(q, " ")
	cmd := exec.Command("sh", "-c", str)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, out)
	}
	return nil
}
