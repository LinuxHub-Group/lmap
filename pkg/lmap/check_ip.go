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
	"net"
	"os"
	"sync"
	"time"
)

const OUTPUT_IP_PER_LINE = 5

func CheckIP(subnet string, isVerbose bool) {
	checkerGroup := &sync.WaitGroup{}
	t := time.Now()
	hosts, _ := GetAllIPsFromCIDR(subnet)
	isUsedList := make([]bool, len(hosts))
	for index, _ := range hosts {
		time.Sleep(500)
		checkerGroup.Add(1)
		go func(index int) {
			defer checkerGroup.Done()
			isUsed := Ping(hosts[index])
			isUsedList[index] = isUsed
			if isVerbose {
				if isUsed {
					println("已使用IP：", hosts[index].String())
				} else {
					println("未使用IP：", hosts[index].String())
				}
			}
		}(index)
	}
	checkerGroup.Wait()
	elapsed := time.Since(t)
	fmt.Fprintln(os.Stderr, "IP扫描完成,耗时", elapsed)
	println("已使用IP：")
	printIPList(hosts, true, isUsedList)
	println("未使用IP：")
	printIPList(hosts, false, isUsedList)
}

func printIPList(hosts []net.IP, boolFilter bool, boolFilterTargetList []bool) {
	firstIndex := 0
	count := 0
	for index, ip := range hosts {
		if boolFilterTargetList[index] == boolFilter {
			fmt.Print(ip.String())
			firstIndex = index
			count = 1
			break
		}
	}
	for index := firstIndex + 1; index < len(hosts); index++ {
		if boolFilterTargetList[index] == boolFilter {
			if count%OUTPUT_IP_PER_LINE == 0 {
				fmt.Println("")
			} else {
				fmt.Print(", ")
			}
			fmt.Print(hosts[index].String())
			count++
		}
	}
	fmt.Println("")
}
