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
