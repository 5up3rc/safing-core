// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package netutils

import "net"

// IPIsLocal determines wheter the given IP is a site-local or link-local address
func IPIsLocal(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil { // check if valid IPv4
		return (ip4[0] == 10) || // site local: 10/8
			(ip4[0] == 172 && ip4[1]&0xf0 == 16) || // site local: 172.16/12
			(ip4[0] == 192 && ip4[1] == 168) || // site local: 192.168/16
			(ip4[0] == 169 && ip4[1] == 254) // link local: 169.254/16
	}
	return len(ip) == net.IPv6len && // check if valid IPv6
		(ip[0]&0xfe == 0xfc || // site local: fc00::/7
			(ip[0] == 0xfe && ip[1]&0xc0 == 0x80)) // link local: fe80::/10
}
