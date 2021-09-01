package main

import (
	"os"
	"flag"
	"sort"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var icmp ICMP

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func main() {
	isVerbose := false
	flag.BoolVar(&isVerbose, "v", false, "be verbose")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Printf("使用方法：%s [-v] <网络号>/<CIDR>\n", os.Args[0])
		os.Exit(-1)
	}
	CheckIP(args[0], isVerbose)
}

func CheckIP(subnet string, isVerbose bool) {
	checkerGroup := &sync.WaitGroup{}
	t := time.Now()
	hosts, _ := getAllHostsFromCIDR(subnet)
	usedIPChan := make(chan string, len(hosts))
	unusedIPChan := make(chan string, len(hosts))
	for _, ip := range hosts {
		checkerGroup.Add(1)
		go func(ip string) {
			defer checkerGroup.Done()
			isUsed := ping(ip)
			if isUsed {
				usedIPChan <- ip
				if isVerbose {
					fmt.Println("已使用IP：", ip)
				}
			} else {
				unusedIPChan <- ip
				if isVerbose {
					fmt.Println("未使用IP：", ip)
				}
			}
		}(ip)
	}
	checkerGroup.Wait()
	elapsed := time.Since(t)
	fmt.Println("IP扫描完成,耗时", elapsed)
	usedIP := make([]string, 0, len(hosts))
	unusedIP := make([]string, 0, len(hosts))
	close(usedIPChan)
	close(unusedIPChan)
	for ipString := range usedIPChan {
		usedIP = append(usedIP, ipString)
	}
	for ipString := range unusedIPChan {
		unusedIP = append(unusedIP, ipString)
	}
	fmt.Println("已使用IP：", sortIPList(usedIP))
	fmt.Println("未使用IP：", sortIPList(unusedIP))
}

func sortIPList(ipStrings []string) (result []string) {
	var ips []net.IP
	for _, ipString := range ipStrings {
		ips = append(ips, net.ParseIP(ipString))
	}
	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})
	for _, ip := range ips {
		result = append(result, ip.String())
	}
	return
}

func getAllHostsFromCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func ping(ip string) bool {
	//开始填充数据包
	icmp.Type = 8 //8->echo message  0->reply message
	icmp.Code = 0
	icmp.Checksum = 0
	icmp.Identifier = 0
	icmp.SequenceNum = 0

	recvBuf := make([]byte, 32)
	var buffer bytes.Buffer

	//先在buffer中写入icmp数据报求去校验和
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.Checksum = CheckSum(buffer.Bytes())
	//然后清空buffer并把求完校验和的icmp数据报写入其中准备发送
	buffer.Reset()
	binary.Write(&buffer, binary.BigEndian, icmp)

	conn, err := net.DialTimeout("ip4:icmp", ip, 1 * time.Second)
	if err != nil {
		return false
	}
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		log.Println("conn.Write error:", err)
		return false
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	num, err := conn.Read(recvBuf)
	if err != nil {
		return false
	}

	conn.SetReadDeadline(time.Time{})

	return string(recvBuf[0:num]) != ""
}

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}
