/*
 *     lmap (LinuxHub's Nmap) is the nmap next generation pro plus max.
 *     Copyright (C) <2021>  <LinuxHub-Group>
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU General Public License as published by
 *     the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU General Public License for more details.
 *
 *     You should have received a copy of the GNU General Public License
 *     along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package lmap

import (
	"log"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// Ping sends an ICMP echo request to the given IP address and waits for a reply
// Returns true if a reply is received within the timeout, false otherwise
func Ping(ip net.IP) bool {
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   12345,
			Seq:  1,
			Data: []byte("lmap ping"),
		},
	}

	sendBytes, err := msg.Marshal(nil)
	if err != nil {
		log.Println("marshal icmp message", err)
		return false
	}

	// Start listening for icmp replies
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Println("dial error:", err)
		return false
	}
	defer conn.Close()

	_, err = conn.WriteTo(sendBytes, &net.IPAddr{
		IP: ip,
	})
	if err != nil {
		log.Println("write error:", err)
		return false
	}

	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 2))

	// Only read a limited number of times to avoid infinite loop
	for i := 0; i < 5; i++ {
		recvBuf := make([]byte, 1500)
		n, addr, err := conn.ReadFrom(recvBuf)
		if err != nil {
			// Timeout or other error, return false
			return false
		}

		if addr.String() == ip.String() {
			// Parse the received message to ensure it's an echo reply
			reply, err := icmp.ParseMessage(1, recvBuf[:n])
			if err != nil {
				continue
			}

			if reply.Type == ipv4.ICMPTypeEchoReply {
				return true
			}
		}
	}

	return false
}
