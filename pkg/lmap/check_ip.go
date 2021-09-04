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

const OUTPUT_IP_PER_LINE = 3

func CheckIP(subnet string, isVerbose bool) {
	checkerGroup := &sync.WaitGroup{}
	t := time.Now()
	hosts, _ := GetAllIPsFromCIDR(subnet)
	for index := range hosts {
		time.Sleep(500)
		checkerGroup.Add(1)

		go func(index int) {
			defer checkerGroup.Done()
			hosts[index].isUsed = Ping(hosts[index].host)
			if isVerbose {
				if hosts[index].isUsed {
					println("已使用IP：", hosts[index].host.String())
				} else {
					println("未使用IP：", hosts[index].host.String())
				}
			}
		}(index)
	}
	checkerGroup.Wait()
	elapsed := time.Since(t)
	_, _ = fmt.Fprintln(os.Stderr, "IP扫描完成,耗时", elapsed)
	fmt.Println("已使用IP：")
	printIPList(hosts, true)
	fmt.Println("未使用IP：")
	printIPList(hosts, false)
}

func printIPList(hosts []HostInfo, boolFilter bool) {

	position := 1

	for _,hostInfo :=range hosts{
		if boolFilter==hostInfo.isUsed {
			fmt.Print(hostInfo.host.String())
			if position%OUTPUT_IP_PER_LINE == 0 {
				fmt.Println()
			} else {
				fmt.Print(", ")
			}
			position++
		}
	}
	fmt.Println("")
}
