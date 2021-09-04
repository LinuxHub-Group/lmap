/*
 *
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
	"github.com/LinuxHub-Group/lmap/pkg/lmap"
	"os"
)

func main() {
	isVerbose := false
	flag.BoolVar(&isVerbose, "v", false, "be verbose")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "使用方法：%s [-v] <网络号>/<CIDR>\n", os.Args[0])
		os.Exit(-1)
	}
	lmap.CheckIP(args[0], isVerbose)
}
