package records

import (
	"fmt"
	"net"
)

func validateEnabledServices(c *Config) error {
	if !c.DNSOn && !c.HTTPOn {
		return fmt.Errorf("Either DNS or HTTP server should be on")
	}
	if len(c.Masters) == 0 && c.Zk == "" {
		return fmt.Errorf("Specify mesos masters or zookeeper in config.json")
	}
	return nil
}

// validateMasters checks that each master in the list is a properly formatted host|IP:port pair.
// duplicate masters in the list are not allowed.
// returns nil if the masters list is empty, or else all masters in the list are valid.
func validateMasters(ms []string) error {
	if err := validateHostPorts(ms, false, "5050", true); err != nil {
		return fmt.Errorf("Error validating masters: %v", err)
	}
	return nil
}

// validateResolvers checks that each resolver in the list is a properly formatted IP address || IP:port pair.
// duplicate resolvers in the list are not allowed.
// returns nil if the resolver list is empty, or else all resolvers in the list are valid.
func validateResolvers(rs []string) error {
	if err := validateHostPorts(rs, true, "53", false); err != nil {
		return fmt.Errorf("Error validating resolvers: %v", err)
	}
	return nil
}

func validateHostPorts(hostPorts []string, ipRequired bool, defaultPort string, portRequired bool) error {
	if len(hostPorts) == 0 {
		return nil
	}
	valid := make(map[string]struct{}, len(hostPorts))
	for _, hp := range hostPorts {
		normalized, err := normalizeValidateHostPort(hp, ipRequired, defaultPort, portRequired)
		if err != nil {
			return err
		}
		if _, found := valid[normalized]; found {
			return fmt.Errorf("Duplicate host specified: %v", normalized)
		}
		valid[normalized] = struct{}{}
	}
	return nil
}

func normalizeValidateHostPort(hostPort string, ipRequired bool, defaultPort string, portRequired bool) (string, error) {
	host, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		if portRequired {
			return "", fmt.Errorf("Illegal host:port specified: %v. Error: %v", hostPort, err)
		}
		host = hostPort
		port = defaultPort
	}
	ip := net.ParseIP(host)
	if ip == nil {
		if ipRequired {
			return "", fmt.Errorf("Illegal ip specified: %v", host)
		}
	} else {
		//TODO(jdef) distinguish between intended hostnames and invalid ip addresses
		host = ip.String()
	}

	return host + "_" + port, nil
}

// validateIPSources checks validity of ip sources
func validateIPSources(srcs []string) error {
	if len(srcs) == 0 {
		return fmt.Errorf("empty ip sources")
	}
	if len(srcs) != len(unique(srcs)) {
		return fmt.Errorf("duplicate ip source specified")
	}
	for _, src := range srcs {
		switch src {
		case "host", "docker", "mesos", "netinfo":
		default:
			return fmt.Errorf("invalid ip source %q", src)
		}
	}

	return nil
}
