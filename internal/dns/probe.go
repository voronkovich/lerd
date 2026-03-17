package dns

import (
	"net"
)

// Check resolves test-lerd-probe.{tld} and checks if 127.0.0.1 is in the results.
// Returns (true, nil) if DNS is working correctly for the given TLD.
func Check(tld string) (bool, error) {
	host := "test-lerd-probe." + tld
	addrs, err := net.LookupHost(host)
	if err != nil {
		return false, nil //nolint:nilerr // DNS failure is a probe negative, not an error
	}

	for _, addr := range addrs {
		if addr == "127.0.0.1" {
			return true, nil
		}
	}
	return false, nil
}
