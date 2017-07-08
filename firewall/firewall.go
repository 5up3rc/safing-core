// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package firewall

import (
	"net"
	"os"
	"sync/atomic"
	"time"

	"github.com/Safing/safing-core/configuration"
	"github.com/Safing/safing-core/firewall/inspection"
	"github.com/Safing/safing-core/firewall/interception"
	"github.com/Safing/safing-core/log"
	"github.com/Safing/safing-core/modules"
	"github.com/Safing/safing-core/network"
	"github.com/Safing/safing-core/network/packet"
	"github.com/Safing/safing-core/process"
	"github.com/Safing/safing-core/sheriff"
)

var (
	firewallModule  *modules.Module
	localNet        net.IPNet
	localhost       net.IP
	dnsServer       net.IPNet
	packetsAccepted *uint64
	packetsBlocked  *uint64
	packetsDropped  *uint64

	config = configuration.Get()
)

func init() {
	localNet = net.IPNet{
		IP:   net.IPv4(127, 0, 0, 0),
		Mask: net.IPv4Mask(255, 0, 0, 0),
	}
	localhost = net.IPv4(127, 0, 0, 1)
	var pA uint64
	packetsAccepted = &pA
	var pB uint64
	packetsBlocked = &pB
	var pD uint64
	packetsDropped = &pD
}

func Start() {
	firewallModule = modules.Register("Firewall", 128)
	defer firewallModule.StopComplete()

	// start interceptor
	go interception.Start()
	go statLogger()

	// go run()
	// go run()
	// go run()
	run()
}

func handlePacket(pkt packet.Packet) {

	// log.Tracef("handling packet")

	// allow anything local, that is not dns
	if pkt.MatchesRemoteIP(packet.Local, localNet) && !(pkt.GetTCPUDPHeader() != nil && pkt.GetTCPUDPHeader().DstPort == 53) {
		pkt.PermanentAccept()
		return
	}

	// allow ICMP and IGMP
	// TODO: actually handle these
	switch pkt.GetIPHeader().Protocol {
	case packet.ICMP:
		pkt.PermanentAccept()
		return
	case packet.ICMPv6:
		pkt.PermanentAccept()
		return
	case packet.IGMP:
		pkt.PermanentAccept()
		return
	}

	// log.Debugf("firewall: pkt %s has ID %s", pkt, pkt.GetConnectionID())

	// use this to time how long it takes process packet
	// timed := time.Now()
	// defer log.Tracef("firewall: took %s to process packet %s", time.Now().Sub(timed).String(), pkt)

	// associate packet to link and handle
	link, created := network.GetOrCreateLinkByPacket(pkt)
	if created {
		link.SetFirewallHandler(initialHandler)
		link.HandlePacket(pkt)
		return
	}
	if link.FirewallHandlerIsSet() {
		link.HandlePacket(pkt)
		return
	}
	verdict(pkt, link.Verdict)

}

func initialHandler(pkt packet.Packet, link *network.Link) {

	// get Connection
	connection, status := network.GetConnectionByFirstPacket(pkt)
	switch status {
	case process.NoSocket:
		// log.Tracef("firewall: unsolicited packet (could not find socket), dropping link: %s", pkt.String())
		link.UpdateVerdict(network.DROP)
		verdict(pkt, network.DROP)
		return
	case process.NoProcess:
		log.Warningf("firewall: could not find process of packet, dropping link: %s", pkt.String())
		link.UpdateVerdict(network.DROP)
		verdict(pkt, network.DROP)
		return
	case process.NoProcessInfo:
		log.Warningf("firewall: could not get process info of packet, dropping link: %s", pkt.String())
		link.UpdateVerdict(network.DROP)
		verdict(pkt, network.DROP)
		return
	}

	// reroute dns requests to nameserver
	if connection.Process().Pid != os.Getpid() && pkt.IsOutbound() && pkt.GetTCPUDPHeader() != nil && !pkt.GetIPHeader().Dst.Equal(localhost) && pkt.GetTCPUDPHeader().DstPort == 53 {
		pkt.RerouteToNameserver()
		return
	}

	// persist connection
	connection.CreateInProcessNamespace()

	// add new Link to Connection
	connection.AddLink(link, pkt)

	// make a decision if not made already
	if connection.Verdict == network.UNDECIDED {
		sheriff.DecideOnConnection(connection, pkt)

	}
	if connection.Verdict != network.CANTSAY {
		link.UpdateVerdict(connection.Verdict)
	} else {
		sheriff.DecideOnLink(connection, link, pkt)
	}

	// log decision
	logInitialVerdict(link)

	if link.Inspect {
		link.SetFirewallHandler(inspectThenVerdict)
		inspectThenVerdict(pkt, link)
	} else {
		link.StopFirewallHandler()
		verdict(pkt, link.Verdict)
	}

}

func inspectThenVerdict(pkt packet.Packet, link *network.Link) {
	pktVerdict, continueInspection := inspection.RunInspectors(pkt, link)
	if continueInspection {
		// do not allow to circumvent link decision: e.g. to ACCEPT packets from a DROP-ed link
		if pktVerdict > link.Verdict {
			verdict(pkt, pktVerdict)
		} else {
			verdict(pkt, link.Verdict)
		}
		return
	}

	// we are done with inspecting
	link.StopFirewallHandler()

	config.Changed()
	config.RLock()
	link.VerdictPermanent = config.PermanentVerdicts
	config.RUnlock()

	link.Save()
	permanentVerdict(pkt, link.Verdict)
}

func permanentVerdict(pkt packet.Packet, action network.Verdict) {
	switch action {
	case network.ACCEPT:
		atomic.AddUint64(packetsAccepted, 1)
		pkt.PermanentAccept()
		return
	case network.BLOCK:
		atomic.AddUint64(packetsBlocked, 1)
		pkt.PermanentBlock()
		return
	case network.DROP:
		atomic.AddUint64(packetsDropped, 1)
		pkt.PermanentDrop()
		return
	}
	pkt.Drop()
}

func verdict(pkt packet.Packet, action network.Verdict) {
	switch action {
	case network.ACCEPT:
		atomic.AddUint64(packetsAccepted, 1)
		pkt.Accept()
		return
	case network.BLOCK:
		atomic.AddUint64(packetsBlocked, 1)
		pkt.Block()
		return
	case network.DROP:
		atomic.AddUint64(packetsDropped, 1)
		pkt.Drop()
		return
	}
	pkt.Drop()
}

func logInitialVerdict(link *network.Link) {
	// switch link.Verdict {
	// case network.ACCEPT:
	// 	log.Infof("firewall: accepting new link: %s", link.String())
	// case network.BLOCK:
	// 	log.Infof("firewall: blocking new link: %s", link.String())
	// case network.DROP:
	// 	log.Infof("firewall: dropping new link: %s", link.String())
	// }
}

func logChangedVerdict(link *network.Link) {
	// switch link.Verdict {
	// case network.ACCEPT:
	// 	log.Infof("firewall: change! - now accepting link: %s", link.String())
	// case network.BLOCK:
	// 	log.Infof("firewall: change! - now blocking link: %s", link.String())
	// case network.DROP:
	// 	log.Infof("firewall: change! - now dropping link: %s", link.String())
	// }
}

func run() {

packetProcessingLoop:
	for {
		select {
		case <-firewallModule.Stop:
			break packetProcessingLoop
		case pkt := <-interception.Packets:
			handlePacket(pkt)
		}
	}

}

func statLogger() {
	for {
		time.Sleep(10 * time.Second)
		log.Tracef("firewall: packets accepted %d, blocked %d, dropped %d", atomic.LoadUint64(packetsAccepted), atomic.LoadUint64(packetsBlocked), atomic.LoadUint64(packetsDropped))
		atomic.StoreUint64(packetsAccepted, 0)
		atomic.StoreUint64(packetsBlocked, 0)
		atomic.StoreUint64(packetsDropped, 0)
	}
}
