package main

import (
	"IpProxyPool/cmd"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 检查或设置命令行参数
	cmd.Execute()
}
