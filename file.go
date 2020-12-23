package main

import (
	"io/ioutil"
	"log"
	"os"
)

func mvfile(s, d string) {
	r, err := ioutil.ReadDir(s)
	if err != nil {
		log.Panicf("移动文件时，读取目录失败%v", err)
	}

	for _, fileName := range r {
		if fileName.Name() == ".git" {
			continue
		}
		os.Rename(s+"/"+fileName.Name(), d+"/"+fileName.Name())
	}
}

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
