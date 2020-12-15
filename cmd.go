package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func localcmd(cmdstr string) {
	cmds := exec.Command(cmdstr)
	log.Printf("执行命令:%v", cmds.Args)
	stdout, err := cmds.StdoutPipe()
	if err != nil {
		log.Panicf("读取命令执行结果失败：%v", err)
	}
	err = cmds.Start()
	if err != nil {
		log.Panicf("命令执行报错:%v", err)
	}

	reader := bufio.NewReader(stdout)

	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}
	cmds.Wait()
}
