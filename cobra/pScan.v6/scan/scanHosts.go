package scan

import (
	"fmt"
	"net"
	"time"
)

// Portstate represent the state of a single TCP port
type PortState struct {
	Port int
	Open state
}

type state bool

// String converts the boolean value od state to a human readable string
func (s state) String() string {
	if s {
		return "open"
	}

	return "closed"
}

// scanport performs a port scan on a single port
func scanPort(host string, port int) PortState {
	p := PortState{
		Port: port,
	}

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	scanConn, err := net.DialTimeout("tcp", address, 1*time.Second)

	if err != nil {
		return p
	}

	scanConn.Close()
	p.Open = true
	return p
}

// Result represents the scan results for a single host
type Results struct {
	Host      string
	NotFound  bool
	PortState []PortState
}

// Run performs a prot port scan on the hosts list
func Run(hl *HostsList, ports []int) []Results {

	res := make([]Results, 0, len(hl.Hosts))

	for _, h := range hl.Hosts {
		r := Results{
			Host: h,
		}

		if _, err := net.LookupHost(h); err != nil {
			r.NotFound = true
			res = append(res, r)
			continue
		}

		for _, p := range ports {
			r.PortState = append(r.PortState, scanPort(h, p))
		}

		res = append(res, r)
	}

	return res
}
