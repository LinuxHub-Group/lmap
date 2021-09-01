package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

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
	usedIP := make([]string, len(hosts))
	unusedIP := make([]string, len(hosts))
	var (
		usedIndex   int64 = 0
		unusedIndex int64 = 0
	)
	for _, ip := range hosts {
		time.Sleep(500)
		checkerGroup.Add(1)
		go func(data string) {
			defer checkerGroup.Done()
			isUsed := ping(data)
			if isUsed {
				old := atomic.LoadInt64(&usedIndex)
				for !atomic.CompareAndSwapInt64(&usedIndex, old, old+1) {
					old = atomic.LoadInt64(&usedIndex)
				}
				usedIP[old] = data
				if isVerbose {
					fmt.Println("已使用IP：", data)
				}
			} else {
				old := atomic.LoadInt64(&unusedIndex)
				for !atomic.CompareAndSwapInt64(&unusedIndex, old, old+1) {
					old = atomic.LoadInt64(&unusedIndex)
				}
				unusedIP[old] = data
				if isVerbose {
					fmt.Println("未使用IP：", data)
				}
			}
		}(ip)
	}
	checkerGroup.Wait()
	elapsed := time.Since(t)
	fmt.Println("IP扫描完成,耗时", elapsed)
	fmt.Println("已使用IP：", sortIPList(usedIP[:usedIndex]))
	fmt.Println("未使用IP：", sortIPList(unusedIP[:unusedIndex]))
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
	var icmp ICMP
	//开始填充数据包
	icmp.Type = 8 //8->echo message  0->reply message

	recvBuf := make([]byte, 32)
	var buffer bytes.Buffer

	//先在buffer中写入icmp数据报求去校验和
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.Checksum = CheckSum(buffer.Bytes())
	//然后清空buffer并把求完校验和的icmp数据报写入其中准备发送
	buffer.Reset()
	binary.Write(&buffer, binary.BigEndian, icmp)

	conn, err := net.DialTimeout("ip4:icmp", ip, 1*time.Second)
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
