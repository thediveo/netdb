package netdb_test

import (
	"fmt"

	"github.com/thediveo/netdb"
)

// Look up a service by its name and protocol name.
func Example_serviceByName() {
	dns := netdb.ServiceByName("domain", "udp")
	fmt.Println(dns.Port)
	// Output: 53
}

// Look up a service by its port number and protocol name.
func Example_serviceByPort() {
	dns := netdb.ServiceByPort(53, "udp")
	fmt.Println(dns.Name)
	// Output: domain
}

// Look up a (TCP/IP subsystem) protocol by its name.
func Example_protocolByName() {
	tcp := netdb.ProtocolByName("tcp")
	fmt.Println(tcp.Number)
	// Output: 6
}

// Looks up a (TCP/IP subsystem) protocol by its (uint8) number.
func Example_protocolByNumber() {
	udp := netdb.ProtocolByNumber(17)
	fmt.Println(udp.Name)
	// Output: udp
}
