package snowflake

import (
	"errors"
	"fmt"
	"net"
)

func privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if ipnet.IP.IsLoopback() {
			// uncomment to test w/o internet
			// return ipnet.IP.To4(), nil
		}
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func lower16BitPrivateIP() (uint16, error) {
	ip, err := privateIPv4()
	if err != nil {
		return 0, err
	}

	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}

var defaultNode *Node

func Init() error {
	nodeID, err := lower16BitPrivateIP()
	if err != nil {
		return err
	}

	defaultNode, err = NewNode(int64(nodeID))
	if err != nil {
		return fmt.Errorf("snowflake error creating NewNode, %s", err)
	}
	return nil
}

// New .
func New() ID {
	return defaultNode.Generate()
}

// NewString generate a string id eg. "17054961268909056"
func NewString() string {
	return defaultNode.Generate().String()
}

// NewBase64  generate a base64 id eg. "MTcwNTQ5NjEyNjg5MDkwNTk="
func NewBase64() string {
	return defaultNode.Generate().Base64()
}

// NewBase58 generate a base58 id eg. "3ibdUVn1Z3"
func NewBase58() string {
	return defaultNode.Generate().Base58()
}

// NewBase36 generate a base36 id eg. "4nxh91lgh6t"
func NewBase36() string {
	return defaultNode.Generate().Base36()
}

// NewBase32 generate a base32 id eg. "xrzcsqbq3yg"
func NewBase32() string {
	return defaultNode.Generate().Base32()
}

// NewInt64 generate a Int64 id eg. 17054961268909056
func NewInt64() int64 {
	return defaultNode.Generate().Int64()
}
