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

func Ping(ip net.IP) bool {
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.RawBody{
			Data: ip,
		},
	}

	sendBytes, err := msg.Marshal(nil)
	if err != nil {
		log.Println("marshal icmp message", err)
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
		return false
	}
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 2))

	for {
		recvBuf := make([]byte, 20)
		_, addr, err := conn.ReadFrom(recvBuf)
		if err != nil {
			return false
		}

		if addr.String() == ip.String() {
			return true
		}
	}
}
