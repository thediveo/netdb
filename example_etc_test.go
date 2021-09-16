package netdb_test

import (
	"fmt"

	"github.com/thediveo/netdb"
)

// Where required, merges protocol and service descriptions from /etc/protocols
// and /etc/services with the built-in database, replacing builtin descriptions
// with those found in the files.
func Example_mergeEtc() {
	etcprotocols, _ := netdb.LoadProtocols("/etc/protocols")
	netdb.Protocols.MergeIndex(etcprotocols)
	etcservices, _ := netdb.LoadServices("/etc/services", netdb.Protocols)
	netdb.Services.MergeIndex(etcservices)
	dns := netdb.ServiceByName("domain", "udp")
	fmt.Printf("%s: %d via %s", dns.Name, dns.Port, dns.Protocol.Name)
	// Output: domain: 53 via udp
}

// Where required, uses protocol and service descriptions only from
// /etc/protocols and /etc/services, ignoring the builtin data completely.
func Example_onlyEtc() {
	etcprotocols, _ := netdb.LoadProtocols("/etc/protocols")
	netdb.Protocols = etcprotocols
	etcservices, _ := netdb.LoadServices("/etc/services", netdb.Protocols)
	netdb.Services = etcservices
	dns := netdb.ServiceByName("domain", "udp")
	fmt.Printf("%s: %d via %s", dns.Name, dns.Port, dns.Protocol.Name)
	// Output: domain: 53 via udp
}
