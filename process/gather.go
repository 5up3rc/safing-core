// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package process

import (
	"net"
	"time"

	"github.com/safing/safing-core/log"
	"github.com/safing/safing-core/network/packet"
)

const (
	Success uint8 = iota
	NoSocket
	NoProcess
	NoProcessInfo
)

func GetPidOfConnection(localIP *net.IP, localPort uint16, protocol uint8) (pid int, status uint8) {
	uid, inode, ok := getConnectionSocket(localIP, localPort, protocol)
	if !ok {
		uid, inode, ok = getListeningSocket(localIP, localPort, protocol)
		// for i := 0; i < 3 && !ok; i++ {
		// 	// give kernel some time, then try again
		// 	log.Tracef("process: giving kernel some time to think")
		// 	time.Sleep(10 * time.Millisecond)
		// 	uid, inode, ok = getConnectionSocket(localIP, localPort, protocol)
		// 	if !ok {
		// 		uid, inode, ok = getListeningSocket(localIP, localPort, protocol)
		// 	}
		// }
		if !ok {
			return -1, NoSocket
		}
	}
	pid, ok = GetPidOfInode(uid, inode)
	for i := 0; i < 3 && !ok; i++ {
		// give kernel some time, then try again
		// log.Tracef("process: giving kernel some time to think")
		time.Sleep(10 * time.Millisecond)
		pid, ok = GetPidOfInode(uid, inode)
	}
	if !ok {
		return -1, NoProcess
	}
	return
}

func GetPidOfIncomingConnection(localIP *net.IP, localPort uint16, protocol uint8) (pid int, status uint8) {
	uid, inode, ok := getListeningSocket(localIP, localPort, protocol)
	if !ok {
		return -1, NoSocket
	}
	pid, ok = GetPidOfInode(uid, inode)
	if !ok {
		return -1, NoProcess
	}
	return
}

func GetPidByPacket(pkt packet.Packet) (pid int, direction bool, status uint8) {
	var protocol uint8
	switch {
	case pkt.GetIPHeader().Protocol == packet.TCP && pkt.IPVersion() == packet.IPv4:
		protocol = TCP4
	case pkt.GetIPHeader().Protocol == packet.UDP && pkt.IPVersion() == packet.IPv4:
		protocol = UDP4
	case pkt.GetIPHeader().Protocol == packet.TCP && pkt.IPVersion() == packet.IPv6:
		protocol = TCP6
	case pkt.GetIPHeader().Protocol == packet.UDP && pkt.IPVersion() == packet.IPv6:
		protocol = UDP6
	default:
		return -1, false, NoSocket
	}

	if pkt.IsOutbound() {
		direction = false
		pid, status = GetPidOfConnection(&pkt.GetIPHeader().Src, pkt.GetTCPUDPHeader().SrcPort, protocol)
		if status == NoSocket {
			pid, status = GetPidOfIncomingConnection(&pkt.GetIPHeader().Src, pkt.GetTCPUDPHeader().SrcPort, protocol)
			if status == Success {
				direction = true
			}
		}
		return
	}
	direction = true
	pid, status = GetPidOfIncomingConnection(&pkt.GetIPHeader().Dst, pkt.GetTCPUDPHeader().DstPort, protocol)
	if status == NoSocket {
		pid, status = GetPidOfConnection(&pkt.GetIPHeader().Dst, pkt.GetTCPUDPHeader().DstPort, protocol)
		if status == Success {
			direction = false
		}
	}
	return

}

func GetProcessByPid(pid int) *Process {
	process, err := GetOrFindProcess(pid)
	if err != nil {
		log.Warningf("process: failed to get process %d: %s", pid, err)
		return nil
	}
	return process
}

func GetProcessOfIncomingConnection(localIP *net.IP, localPort uint16, protocol uint8) (process *Process, status uint8) {
	pid, status := GetPidOfIncomingConnection(localIP, localPort, protocol)
	if status == Success {
		process = GetProcessByPid(pid)
		if process == nil {
			return nil, NoProcessInfo
		}
	}
	return
}

func GetProcessOfConnection(localIP *net.IP, localPort uint16, protocol uint8) (process *Process, status uint8) {
	pid, status := GetPidOfConnection(localIP, localPort, protocol)
	if status == Success {
		process = GetProcessByPid(pid)
		if process == nil {
			return nil, NoProcessInfo
		}
	}
	return
}

func GetProcessByPacket(pkt packet.Packet) (process *Process, direction bool, status uint8) {
	pid, direction, status := GetPidByPacket(pkt)
	if status == Success {
		process = GetProcessByPid(pid)
		if process == nil {
			return nil, direction, NoProcessInfo
		}
	}
	return
}
