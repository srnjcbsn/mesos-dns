package labels

import (
	"strings"
)

// Sep is the default domain fragment separator.
const Sep = "."

// DomainFrag mangles the given name in order to produce a valid domain fragment.
// A valid domain fragment will consist of one or more host name labels
// concatenated by the given separator.
func DomainFrag(name, sep string, label Func) string {
	var labels []string
	for _, part := range strings.Split(name, sep) {
		if lab := label(part); lab != "" {
			labels = append(labels, lab)
		}
	}
	return strings.Join(labels, sep)
}

// Func is a function type representing label functions.
type Func func(string) string

// RFC952 mangles a name to conform to the DNS label rules specified in RFC952.
// See http://www.rfc-base.org/txt/rfc-952.txt
func RFC952(name string) string {
	return newState([]byte(name), 24, "-0123456789", "-").run()
}

// RFC1123 mangles a name to conform to the DNS label rules specified in RFC1123.
// See http://www.rfc-base.org/txt/rfc-1123.txt
func RFC1123(name string) string {
	return newState([]byte(name), 63, "-", "-").run()
}

// min returns the minimum of two ints
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
