// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package intel

import (
	"testing"
	"time"

	"github.com/miekg/dns"
)

func TestResolve(t *testing.T) {
	Resolve("google.com.", dns.Type(dns.TypeA), 0)
	time.Sleep(200 * time.Millisecond)
}
