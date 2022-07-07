package main

import (
	"flag"
	"fmt"
	"github.com/krkrkc/distributed-crontab/master"
)

var (
	confFile string
)

func initArgs() {
	flag.StringVar(&confFile, "config", "./master/main/master.json", "指定master.json")
	flag.Parse()
}

func main() {
	var (
		err error
	)

	//初始化命令行参数
	initArgs()

	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	//启动Api HTTP服务
	master.InitApiServer()

ERR:
	fmt.Println(err)
}
