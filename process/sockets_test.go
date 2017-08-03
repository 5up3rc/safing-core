// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package process

import (
	"log"
	"net"
	"testing"
)

func TestSockets(t *testing.T) {

	updateListeners(TCP4)
	updateListeners(UDP4)
	updateListeners(TCP6)
	updateListeners(UDP6)
	log.Printf("addressListeningTCP4: %v", addressListeningTCP4)
	log.Printf("globalListeningTCP4: %v", globalListeningTCP4)
	log.Printf("addressListeningUDP4: %v", addressListeningUDP4)
	log.Printf("globalListeningUDP4: %v", globalListeningUDP4)
	log.Printf("addressListeningTCP6: %v", addressListeningTCP6)
	log.Printf("globalListeningTCP6: %v", globalListeningTCP6)
	log.Printf("addressListeningUDP6: %v", addressListeningUDP6)
	log.Printf("globalListeningUDP6: %v", globalListeningUDP6)

	getListeningSocket(&net.IPv4zero, 53, TCP4)
	getListeningSocket(&net.IPv4zero, 53, UDP4)
	getListeningSocket(&net.IPv6zero, 53, TCP6)
	getListeningSocket(&net.IPv6zero, 53, UDP6)

	// spotify: 192.168.0.102:37312     192.121.140.65:80
	localIP := net.IPv4(192, 168, 0, 102)
	uid, inode, ok := getConnectionSocket(&localIP, 37312, TCP4)
	log.Printf("getConnectionSocket: %d %d %v", uid, inode, ok)

	activeConnectionIDs := GetActiveConnectionIDs()
	for _, connID := range activeConnectionIDs {
		log.Printf("active: %s", connID)
	}

}
