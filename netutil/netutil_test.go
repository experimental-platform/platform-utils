package netutil

import (
	"reflect"
	"strings"
	"testing"
)

type mocCommandExecutor struct {
	expected []string
	data     []byte
	err      error
}

func (c mocCommandExecutor) Command(name string, arg ...string) ([]byte, error) {
	// TODO: check name + arg against c.expected
	return c.data, c.err
}

// make sure the moc satisfies the interface
var _ CmdExec = (*mocCommandExecutor)(nil)

func TestGetDefaultInterface(t *testing.T) {
	moc := mocCommandExecutor{data: []byte(
		"8.8.8.8 via 172.16.0.1 dev eno1  src 172.16.10.239\n    cache",
	)}
	result, err := GetDefaultInterface(moc)
	if err != nil {
		t.Errorf("Static mode failure: %v", err)
	}
	if !strings.Contains(result, "eno1") {
		t.Errorf("Result should contain 'eno1' but is '%v'.", result)
	}
}

func TestGetInterfaceIndex(t *testing.T) {
	moc := mocCommandExecutor{data: []byte(`4: eno1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP mode DEFAULT group default qlen 1000
link/ether 54:be:f7:66:2c:49 brd ff:ff:ff:ff:ff:ff`)}
	result, err := getInterfaceIndex(moc, "eno1")
	if err != nil {
		t.Errorf("Static mode failure: %v", err)
	}
	if result != "4" {
		t.Errorf("Expected '4', got '%v'.", result)
	}
}

func TestGetInterfaceStateUnconfigured(t *testing.T) {
	// /run/systemd/netif/links/3
	data := []byte(`# This is private data. Do not parse.
ADMIN_STATE=configuring
OPER_STATE=no-carrier
NETWORK_FILE=/usr/lib64/systemd/network/zz-default.network
DNS=
NTP=
DOMAINS=
WILDCARD_DOMAIN=no
LLMNR=yes
`)
	result, err := parseInterfaceState(data)
	if err != nil {
		t.Errorf("Failure: %v", err)
	}
	expected := InterfaceData{
		ADMIN_STATE:     "configuring",
		OPER_STATE:      "no-carrier",
		NETWORK_FILE:    "/usr/lib64/systemd/network/zz-default.network",
		WILDCARD_DOMAIN: false,
		LLMNR:           true,
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected '%+v', got '%+v'.", expected, result)
	}
}

func TestGetInterfaceStateDHCP(t *testing.T) {
	//
	data := []byte(`# This is private data. Do not parse.
ADMIN_STATE=configured
OPER_STATE=routable
NETWORK_FILE=/usr/lib64/systemd/network/zz-default.network
DNS=8.8.8.8 10.11.0.2 62.220.18.8
NTP=
DOMAINS=office.protorz.net
WILDCARD_DOMAIN=no
LLMNR=yes
DHCP_LEASE=/run/systemd/netif/leases/4
`)
	result, err := parseInterfaceState(data)
	if err != nil {
		t.Errorf("Failure: %v", err)
	}
	expected := InterfaceData{
		ADMIN_STATE:     "configured",
		OPER_STATE:      "routable",
		NETWORK_FILE:    "/usr/lib64/systemd/network/zz-default.network",
		DNS:             []string{"8.8.8.8", "10.11.0.2", "62.220.18.8"},
		DOMAINS:         []string{"office.protorz.net"},
		WILDCARD_DOMAIN: false,
		LLMNR:           true,
		DHCP_LEASE:      "/run/systemd/netif/leases/4",
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected '%v', got '%v'.", expected, result)
	}
}
