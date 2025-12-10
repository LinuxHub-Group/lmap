/*
 *     lmap (LinuxHub's Nmap) is the nmap next generation pro plus max.
 *     Copyright (C) <2022>  <LinuxHub-Group>
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
	"fmt"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// ListenICMP listens for ICMP packets on the specified address
// This function can be used to listen for ICMP replies when doing network scanning
func ListenICMP(addr string) error {
	conn, err := icmp.ListenPacket("ip4:icmp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on ICMP: %w", err)
	}
	defer conn.Close()

	fmt.Printf("Listening for ICMP packets on %s...\n", addr)

	for {
		buf := make([]byte, 1500)
		n, peer, err := conn.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("failed to read ICMP packet: %w", err)
		}

		msg, err := icmp.ParseMessage(1, buf[:n]) // 1 = ICMPv4
		if err != nil {
			fmt.Printf("Failed to parse ICMP message from %s: %v\n", peer, err)
			continue
		}

		switch msg.Type {
		case ipv4.ICMPTypeEchoReply:
			fmt.Printf("Received ICMP Echo Reply from %s\n", peer)
		case ipv4.ICMPTypeEcho:
			fmt.Printf("Received ICMP Echo Request from %s\n", peer)
		default:
			fmt.Printf("Received ICMP message type %v from %s\n", msg.Type, peer)
		}
	}
}
