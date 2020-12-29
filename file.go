package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func delfile(file string) {
	r, err := ioutil.ReadDir(file)
	if err != nil {
		log.Panicf("删除目录时，读取目录失败%v", err)
	}
	for _, fileName := range r {
		if fileName.Name() == ".git" {
			continue
		}
		os.RemoveAll(file + "/" + fileName.Name())
	}
}

func copyfile(s, d string) {
	sinfos, err := ioutil.ReadDir(s)
	if err != nil {
		log.Panicf("复制文件时，读取目录失败%v", err)
	}

	for _, sinfo := range sinfos {
		if sinfo.Name() == ".git" {
			continue
		}
		if sinfo.IsDir() {
			err := os.Mkdir(d+"/"+sinfo.Name(), 0766)
			if err != nil {
				log.Panicf("复制时创建目录失败:%v", err)
			}
			copyfile(s+"/"+sinfo.Name(), d+"/"+sinfo.Name())
			continue
		}

		srcfile, err := os.Open(s + "/" + sinfo.Name())
		if err != nil {
			log.Panicf("复制时读取源文件失败:%v", err)
		}
		defer srcfile.Close()

		reader := bufio.NewReader(srcfile)

		dstfile, err := os.OpenFile(d+"/"+sinfo.Name(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0766)
		if err != nil {
			log.Panicf("复制时打开目标文件失败:%v", err)
		}
		defer dstfile.Close()
		writer := bufio.NewWriter(dstfile)

		for {
			buf, err := reader.ReadBytes('\n')
			if err != nil && err != io.EOF {
				log.Printf("读取文件内容失败:%v", err)
			}
			writer.Write(buf)
			writer.Flush()
			if err == io.EOF {
				break
			}

		}
	}
}
