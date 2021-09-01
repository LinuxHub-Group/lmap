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
