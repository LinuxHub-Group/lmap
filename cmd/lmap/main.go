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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/LinuxHub-Group/lmap/pkg/lmap"
)

func main() {
	isVerbose := false
	flag.BoolVar(&isVerbose, "v", false, "详细输出")

	subnets := multiArg{}
	flag.Var(&subnets, "subnet", "要扫描的网络段，CIDR格式 (可多次指定)")

	excludes := multiArg{}
	flag.Var(&excludes, "exclude", "要排除的IP或子网 (可多次指定)")

	flag.Parse()

	if len(subnets) < 1 {
		printUsage()
		os.Exit(1)
	}

	lmap.CheckIP(subnets, excludes, isVerbose)
}

// printUsage prints the usage information for the program
func printUsage() {
	progName := os.Args[0]
	if progName == "" {
		progName = "lmap"
	}

	fmt.Fprintf(os.Stderr, "使用方法: %s [选项] -subnet <网络号>/<CIDR> [-subnet ...]\n", progName)
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "选项:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "示例:")
	fmt.Fprintf(os.Stderr, "  %s -subnet 192.168.1.0/24\n", progName)
	fmt.Fprintf(os.Stderr, "  %s -subnet 192.168.1.0/24 -subnet 10.0.0.0/16\n", progName)
	fmt.Fprintf(os.Stderr, "  %s -subnet 192.168.1.0/24 -exclude 192.168.1.1\n", progName)
	fmt.Fprintf(os.Stderr, "  %s -subnet 192.168.1.0/24 -exclude 192.168.1.10/32 -v\n", progName)
}

// multiArg implements the flag.Value interface for collecting multiple string values
type multiArg []string

func (m *multiArg) String() string {
	return fmt.Sprintf("%v", *m)
}

func (m *multiArg) Set(value string) error {
	*m = append(*m, value)
	return nil
}
