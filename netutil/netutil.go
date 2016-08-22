package netutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

type CmdExec interface {
	Command(name string, arg ...string) ([]byte, error)
}

type RealCmdExec struct{}

func (c RealCmdExec) Command(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()
	return out, err
}

func GetDefaultInterface(executor CmdExec) (string, error) {
	// RADAR: Will only work when box has internet. Can we use network connected state instead?
	// TODO: Broken when more than one interface present
	out, err := executor.Command("ip", "route", "get", "8.8.8.8")
	if err != nil {
		return "", err
	}

	reg, err := regexp.Compile("dev e[nt]+[0-9a-z_]+")
	if err != nil {
		return "", err
	}

	found := reg.Find(out)
	if found == nil {
		return "", fmt.Errorf("getDefaultInterface(): error parsing '%v'", out)
	}

	split := strings.Split(string(found), " ")
	if len(split) != 2 {
		return "", fmt.Errorf("getDefaultInterface(): error parsing '%v'", out)
	}

	return split[1], nil
}

func getInterfaceIndex(executor CmdExec, name string) (string, error) {
	out, err := executor.Command("ip", "link", "show", name)
	if err == nil {
		reg, err := regexp.Compile("^\\d+")
		if err == nil {
			result := reg.Find(out)
			if result == nil {
				err = fmt.Errorf("getInterfaceIndex(): error parsing output of `ip link show %v`", name)
			} else {
				return string(result), nil
			}
		}
	}
	return "", err
}

type InterfaceData struct {
	ADMIN_STATE     string
	OPER_STATE      string
	NETWORK_FILE    string
	DNS             []string
	NTP             string
	DOMAINS         []string
	WILDCARD_DOMAIN bool
	LLMNR           bool
	DHCP_LEASE      string
}

func boolify(word string) bool {
	if strings.ToLower(word) == "yes" {
		return true
	}
	return false
}

func parseInterfaceState(data []byte) (InterfaceData, error) {
	result := new(InterfaceData)
	for _, line := range strings.Split(strings.Trim(string(data), "\n"), "\n") {
		splitLine := strings.Split(strings.Trim(line, "="), "=")
		if len(splitLine) == 0 || len(splitLine) > 2 {
			return *result, errors.New("Parser error (1) on: " + line)
		}
		if len(splitLine) == 2 {
			key, value := strings.Trim(splitLine[0], " "), strings.Trim(splitLine[1], " ")
			switch key {
			case "ADMIN_STATE":
				result.ADMIN_STATE = value
			case "OPER_STATE":
				result.OPER_STATE = value
			case "NETWORK_FILE":
				result.NETWORK_FILE = value
			case "DNS":
				result.DNS = strings.Split(value, " ")
			case "NTP":
				result.NTP = value
			case "DOMAINS":
				result.DOMAINS = strings.Split(value, " ")
			case "WILDCARD_DOMAIN":
				result.WILDCARD_DOMAIN = boolify(value)
			case "LLMNR":
				result.LLMNR = boolify(value)
			case "DHCP_LEASE":
				result.DHCP_LEASE = value
			default:
				return *result, errors.New("Parser error (2) on: " + line)
			}
		}
	}
	return *result, nil
}

func GetInterfaceStats(name string) (InterfaceData, error) {
	// RADAR: This will only work on current linux kernels with systemd
	// RADAR: This is currently untested.
	name, err := getInterfaceIndex(RealCmdExec{}, name)
	if err != nil {
		return InterfaceData{}, err
	}
	path := "/run/systemd/netif/links/" + name
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return InterfaceData{}, err
	}
	result, err := parseInterfaceState(data)
	if err != nil {
		return InterfaceData{}, err
	}
	return result, nil
}
