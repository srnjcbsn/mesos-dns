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

var (
	rfc952Table  = makeDNS952Table()
	rfc1123Table = makeDNS1123Table()
)

// RFC952 mangles a name to conform to the DNS label rules specified in RFC952.
// See http://www.rfc-base.org/txt/rfc-952.txt
func RFC952(name string) string {
	return rfc952Table.label(name, 24)
}

// RFC1123 mangles a name to conform to the DNS label rules specified in RFC1123.
// See http://www.rfc-base.org/txt/rfc-1123.txt
func RFC1123(name string) string {
	return rfc1123Table.label(name, 63)
}
