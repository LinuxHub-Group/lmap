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
	"net"
	"strings"
)

func GetAllIPsFromCIDR(cidr string, excludes []string) ([]HostInfo, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []HostInfo
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ipStr := ip.String()
		if !isExcluded(ipStr, excludes) {
			ips = append(ips, HostInfo{
				host:   dupIP(ip),
				isUsed: false,
			})
		}
	}

	// For larger networks, remove network and broadcast addresses
	// For very small networks (like /31, /32), keep all addresses
	ones, bits := ipNet.Mask.Size()
	if bits-ones > 2 && len(ips) > 2 {
		return ips[1 : len(ips)-1], nil // Remove network and broadcast addresses
	}
	return ips, nil
}

func isExcluded(ip string, excludes []string) bool {
	for _, exclude := range excludes {
		if strings.Contains(exclude, "/") {
			_, excludeNet, err := net.ParseCIDR(exclude)
			if err != nil {
				//fmt.Printf("解析排除规则 %s 失败: %v\n", exclude, err)
				continue
			}
			if excludeNet.Contains(net.ParseIP(ip)) {
				//fmt.Printf("IP %s 被排除规则 %s 排除\n", ip, exclude)
				return true
			}
		} else {
			if ip == exclude {
				//fmt.Printf("IP %s 被排除规则 %s 排除\n", ip, exclude)
				return true
			}
		}
	}
	return false
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func dupIP(ip net.IP) net.IP {
	if x := ip.To4(); x != nil {
		ip = x
	}
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}
