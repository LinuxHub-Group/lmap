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
	"encoding/binary"
	"golang.org/x/net/icmp"
	"log"
	"net"
	"time"
)

func Ping(ip net.IP) bool {
	recvBuf := make([]byte, 8)

	// Start listening for icmp replies
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")

	if err != nil {
		log.Println("dial error:", err)
		return false
	}
	defer conn.Close()

	sendBytes := []byte{8, 0, 247, 255, 0, 0, 0, 0}
	//expectCheckSum := int(checkSum([]byte{0, 0, 0, 0, 0, 0, 0, 0}))
	expectCheckSum := 65535

	_, err = conn.WriteTo(sendBytes, &net.IPAddr{
		IP: ip,
	})
	if err != nil {
		return false
	}
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 2))

	_, _, err = conn.ReadFrom(recvBuf)
	if err != nil {
		return false
	}

	recvType := recvBuf[0]
	recvCheckSum := int(binary.BigEndian.Uint16(recvBuf[2:4]))

	_ = conn.SetReadDeadline(time.Time{})

	if recvType != 0 {
		return false
	}

	if recvCheckSum != expectCheckSum {
		return false
	}

	return true
}

//func checkSum(data []byte) uint16 {
//	var (
//		sum    uint32
//		length = len(data)
//		index  int
//	)
//	for length > 1 {
//		sum += uint32(data[index])<<8 + uint32(data[index+1])
//		index += 2
//		length -= 2
//	}
//	if length > 0 {
//		sum += uint32(data[index])
//	}
//	sum += sum >> 16
//
//	return uint16(^sum)
//}
