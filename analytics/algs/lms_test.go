// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package algs

import "testing"

func TestLmsScoreOfDomain(t *testing.T) {
	testDomain(t, "g.symcd.com.", 100, 100)
	testDomain(t, "www.google.com.", 100, 100)
	testDomain(t, "55ttt5.12abc3.test.com.", 50, 50)
	testDomain(t, "mbtq6opnuodp34gcrma65fxacgxv5ukr7lq6xuhr4mhoibe7.yvqptrozfbnqyemchpovw3q5xwjibuxfsgb72mix3znhpfhc.i2n7jh2gadqaadck3zs3vg3hbv5pkmwzeay4gc75etyettbb.isi5mhmowtfriu33uxzmgvjur5g2p3tloynwohfrggee6fkn.meop7kqyd5gwxxa3.er.spotify.com.", 0, 30)
}

func testDomain(t *testing.T, domain string, min, max float64) {
	score := LmsScoreOfDomain(domain)
	if score < min || score > max {
		t.Errorf("domain %s has scored %.2f, but should be between %.0f and %.0f", domain, score, min, max)
	}
}
