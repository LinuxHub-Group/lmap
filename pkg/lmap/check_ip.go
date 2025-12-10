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
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	OUTPUT_IP_PER_LINE   = 3
	MAX_CONCURRENT_PINGS = 100
)

// CheckIP scans the given subnets for active hosts
// subnets: list of CIDR notation networks to scan
// excludes: list of IPs or subnets to exclude from scanning
// isVerbose: if true, print detailed information during scanning
func CheckIP(subnets, excludes []string, isVerbose bool) {
	checkerGroup := &sync.WaitGroup{}
	t := time.Now()

	// Use a semaphore to limit concurrent pings
	semaphore := make(chan struct{}, MAX_CONCURRENT_PINGS)

	var allHosts []HostInfo
	for _, subnet := range subnets {
		hosts, err := GetAllIPsFromCIDR(subnet, excludes)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error parsing subnet %s: %v\n", subnet, err)
			continue
		}
		allHosts = append(allHosts, hosts...)
	}

	for index := range allHosts {
		checkerGroup.Add(1)
		semaphore <- struct{}{} // Acquire semaphore

		go func(index int) {
			defer checkerGroup.Done()
			defer func() { <-semaphore }() // Release semaphore

			allHosts[index].isUsed = Ping(allHosts[index].host)
			if isVerbose {
				if allHosts[index].isUsed {
					fmt.Println("已使用IP：", allHosts[index].host.String())
				} else {
					fmt.Println("未使用IP：", allHosts[index].host.String())
				}
			}
		}(index)
	}

	checkerGroup.Wait()
	elapsed := time.Since(t)
	_, _ = fmt.Fprintln(os.Stderr, "IP扫描完成,耗时", elapsed)
	fmt.Println("已使用IP：")
	printIPList(allHosts, true)
	fmt.Println("未使用IP：")
	printIPList(allHosts, false)
}

// printIPList prints the list of IPs based on the filter
// hosts: list of HostInfo to print
// boolFilter: if true, print used IPs; if false, print unused IPs
func printIPList(hosts []HostInfo, boolFilter bool) {
	position := 1

	for _, hostInfo := range hosts {
		if boolFilter == hostInfo.isUsed {
			fmt.Print(hostInfo.host.String())
			if position%OUTPUT_IP_PER_LINE == 0 {
				fmt.Println()
			} else {
				fmt.Print(", ")
			}
			position++
		}
	}

	// Only add a newline if we've printed something and it didn't end with a newline
	if position > 1 && (position-1)%OUTPUT_IP_PER_LINE != 0 {
		fmt.Println()
	} else if position == 1 {
		// No IPs to print
		fmt.Println("(无)")
	}
}
