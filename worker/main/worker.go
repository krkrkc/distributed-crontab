package main

import (
	"flag"
	"fmt"
	"github.com/krkrkc/distributed-crontab/worker"
	"time"
)

var (
	confFile string
)

func initArgs() {
	flag.StringVar(&confFile, "config", "./worker/main/worker.json", "指定worker.json")
	flag.Parse()
}

func main() {
	var (
		err error
	)

	//初始化命令行参数
	initArgs()

	if err = worker.InitConfig(confFile); err != nil {
		goto ERR
	}

	worker.InitSchduler()

	if err = worker.InitExcutor(); err != nil {
		fmt.Println(err)
		return
	}

	if err = worker.InitJobMgr(); err != nil {
		goto ERR
	}

	for {
		time.Sleep(1 * time.Second)
	}
ERR:
	fmt.Println(err)
}
