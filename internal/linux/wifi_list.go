package linux

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	ghz5 = "5GHz"
	ghz2 = "2GHz"
)

type WiFiPoint struct {
	SSID   string `json:"ssid"`
	Active bool   `json:"active"`
	Freq   string `json:"freq"`
	Signal int    `json:"signal"`
}

type syscalWiFiList interface {
	call() ([]byte, error)
}

type WiFiListCommand struct {
}

func NewWiFiCommand() WiFiListCommand {
	return WiFiListCommand{}
}

func (WiFiListCommand) call() ([]byte, error) {
	cmd := exec.Command("sh", "-c", "nmcli -f ACTIVE,SSID,SIGNAL,FREQ dev wifi list | awk '{print $1, $2, $3, $4}'")
	return cmd.CombinedOutput()
}

func ScanWiFiPoints(sc syscalWiFiList) ([]WiFiPoint, error) {
	raw, err := sc.call()
	if err != nil {
		return nil, fmt.Errorf("failed to scan wifi points: %w - %s", err, raw)
	}
	parsed, err := parseWiFis(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse wifi list: %w", err)
	}
	if len(parsed) == 0 {
		return []WiFiPoint{}, nil
	}
	sort.Slice(parsed, func(i, j int) bool {
		return parsed[i].Signal > parsed[j].Signal
	})
	return parsed, nil
}

var re = regexp.MustCompile(`^(yes|no)\s(.+)\s(\d+)\s(\d+)$`)

func parseWiFis(data []byte) ([]WiFiPoint, error) {
	s := bufio.NewScanner(bytes.NewReader(data))
	list := make([]WiFiPoint, 0, 4)
	for s.Scan() {
		line := s.Text()
		if len(line) == 0 || strings.HasPrefix(line, "ACTIVE") {
			continue
		}
		matches := re.FindStringSubmatch(line)
		if len(matches) != 5 {
			continue
		}
		point := WiFiPoint{}
		point.SSID = matches[2]
		if len(point.SSID) == 0 || point.SSID == "--" {
			continue
		}

		if matches[1] == "yes" {
			point.Active = true
		}

		if len(matches[3]) > 0 {
			sig, err := strconv.Atoi(matches[3])
			if err != nil {
				return nil, fmt.Errorf("%w: %s", err, matches[3])
			}
			point.Signal = sig
		}
		if len(matches[4]) > 0 {
			freq, err := strconv.Atoi(matches[4])
			if err != nil {
				return nil, fmt.Errorf("%w: %s", err, matches[4])
			}
			if freq < 4500 {
				point.Freq = ghz2
			} else if freq > 4500 {
				point.Freq = ghz5
			} else {
				point.Freq = "?"
			}
		}
		list = append(list, point)
	}
	return list, nil
}
