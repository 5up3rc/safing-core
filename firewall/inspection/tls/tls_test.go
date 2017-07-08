// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package tls

import (
	"fmt"
	"testing"

	"github.com/Safing/safing-core/firewall/inspection/tls/tlslib"
)

var clientHelloSample = []byte{
	// 0x16, 0x03, 0x01, 0x01, 0x01,
	0x01, 0x00, 0x00, 0xfd, 0x03, 0x03, 0x58, 0xf3, 0xbf, 0x00, 0xa1, 0x23, 0xff, 0x99, 0xca, 0x9a, 0x3c, 0x16, 0xe3, 0x34, 0xc6, 0xc1, 0xe6, 0xf0, 0xe5, 0x84, 0xff, 0x87, 0x8a, 0x88, 0x04, 0x17, 0xf2, 0xa2, 0xc2, 0x2c, 0x4a, 0x32, 0x00, 0x00, 0x6c, 0xc0, 0x2b, 0xc0, 0x2c, 0xc0, 0x86, 0xc0, 0x87, 0xc0, 0x09, 0xc0, 0x23, 0xc0, 0x0a, 0xc0, 0x24, 0xc0, 0x72, 0xc0, 0x73, 0xc0, 0xac, 0xc0, 0xad, 0xc0, 0x08, 0xc0, 0x2f, 0xc0, 0x30, 0xc0, 0x8a, 0xc0, 0x8b, 0xc0, 0x13, 0xc0, 0x27, 0xc0, 0x14, 0xc0, 0x28, 0xc0, 0x76, 0xc0, 0x77, 0xc0, 0x12, 0x00, 0x9c, 0x00, 0x9d, 0xc0, 0x7a, 0xc0, 0x7b, 0x00, 0x2f, 0x00, 0x3c, 0x00, 0x35, 0x00, 0x3d, 0x00, 0x41, 0x00, 0xba, 0x00, 0x84, 0x00, 0xc0, 0xc0, 0x9c, 0xc0, 0x9d, 0x00, 0x0a, 0x00, 0x9e, 0x00, 0x9f, 0xc0, 0x7c, 0xc0, 0x7d, 0x00, 0x33, 0x00, 0x67, 0x00, 0x39, 0x00, 0x6b, 0x00, 0x45, 0x00, 0xbe, 0x00, 0x88, 0x00, 0xc4, 0xc0, 0x9e, 0xc0, 0x9f, 0x00, 0x16, 0x01, 0x00, 0x00, 0x68, 0x00, 0x17, 0x00, 0x00, 0x00, 0x16, 0x00, 0x00, 0x00, 0x05, 0x00, 0x05, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0b, 0x00, 0x09, 0x00, 0x00, 0x06, 0x6f, 0x72, 0x66, 0x2e, 0x61, 0x74, 0xff, 0x01, 0x00, 0x01, 0x00, 0x00, 0x23, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x0c, 0x00, 0x0a, 0x00, 0x17, 0x00, 0x18, 0x00, 0x19, 0x00, 0x15, 0x00, 0x13, 0x00, 0x0b, 0x00, 0x02, 0x01, 0x00, 0x00, 0x0d, 0x00, 0x16, 0x00, 0x14, 0x04, 0x01, 0x04, 0x03, 0x05, 0x01, 0x05, 0x03, 0x06, 0x01, 0x06, 0x03, 0x03, 0x01, 0x03, 0x03, 0x02, 0x01, 0x02, 0x03, 0x00, 0x10, 0x00, 0x0b, 0x00, 0x09, 0x08, 0x68, 0x74, 0x74, 0x70, 0x2f, 0x31, 0x2e, 0x31,
}

func TestClientHelloParsing(t *testing.T) {
	var msg tlslib.ClientHelloMsg
	if msg.Unmarshal(clientHelloSample) {
		fmt.Printf("ClientHello: %v", msg)
	} else {
		t.Error("Could not parse ClientHello")
	}
}
