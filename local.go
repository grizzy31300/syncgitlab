package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"log"
	"os"
	"strings"
)

func syncCode(nproject, oproject []pdata, stoken, susername, dtoken, dusername string, newpmap map[string][]gjson.Result) {
	log.Printf("多少个老项目:%d", len(oproject))
	var ncprojects []pdata
	cprojects := confyaml.Get("cproject")
	if cprojects != nil {
		for _, ov := range oproject {
			for _, cproject := range cprojects.([]interface{}) {
				if ov.Name == cproject.(string) {
					ncprojects = append(ncprojects, ov)
				}
			}

		}
		oproject = ncprojects
		if len(oproject) < 1 {
			log.Panicln("没有匹配到配置文件中项目名称")
		}
	}

	for _, v := range oproject {
		if v.Name == "Monitoring" {
			continue
		}
		odir := fmt.Sprintf("./old%s", v.Name)
		log.Println(odir, IsExist(odir))
		if !IsExist(odir) {
			err := os.Mkdir(odir, 0755)
			if err != nil {
				log.Printf("创建目录失败:%v", err)
			}
		}
		ndir := fmt.Sprintf("./new%s", v.Name)
		if !IsExist(ndir) {
			err := os.Mkdir(ndir, 0755)
			if err != nil {
				log.Printf("创建目录失败:%v", err)
			}
		}
		ourl := chageUrl(v.Http_url_to_repo, "src.url")
		log.Printf("源git克隆的url:%s", ourl)
		err := gitclone(ourl, odir, susername, stoken)
		if err != nil {
			if err.Error() == "remote repository is empty" {
				log.Println(IsExist(odir))
				if IsExist(odir) {
					err := os.RemoveAll(odir)
					if err != nil {
						log.Panicf("移除目录失败：%v", err)
					}
				}
				log.Println(IsExist(ndir))
				if IsExist(ndir) {
					err := os.RemoveAll(ndir)
					if err != nil {
						log.Panicf("移除目录失败：%v", err)
					}
				}
				continue
			} else {
				os.Exit(1)
			}
		}
		logs := gitlog(odir)
		for _, np := range nproject {
			if np.Name == "Monitoring" {
				continue
			}
			if v.Name == np.Name && v.Name_with_namespace == np.Name_with_namespace {
				nurl := chageUrl(np.Http_url_to_repo, "dest.url")
				log.Printf("目标git克隆的url:%s", nurl)
				err = gitclone(nurl, ndir, dusername, dtoken)
				if err != nil {
					if err.Error() == "remote repository is empty" {
						log.Println(IsExist(odir))
						if IsExist(odir) {
							err := os.RemoveAll(odir)
							if err != nil {
								log.Panicf("移除目录失败：%v", err)
							}
						}
						log.Println(IsExist(ndir))
						if IsExist(ndir) {
							err := os.RemoveAll(ndir)
							if err != nil {
								log.Panicf("移除目录失败：%v", err)
							}
						}
						continue
					} else {
						os.Exit(1)
					}
				}
				log.Printf("commit次数:%d", len(logs))
				if len(logs) < 1 {
					continue
				}
				log.Println("提交路径:%s", ndir)
				log.Println("提交信息:%s", logs[0])
				for _, branch := range newpmap[v.Name_with_namespace] {
					log.Printf("%s pull的分支%s", v.Name, branch.Str)
					track(odir, branch.Str)
					pull(odir, stoken, susername, branch.Str)
					log.Println("pull完成")
					checkoutBranch(odir, branch.Str)
					log.Println("old checkoutBranch完成")
					checkout(ndir, branch.Str)
					log.Println("new chekout完成")
					delfile(ndir)
					log.Println("删除新目录下文件完成")
					copyfile(odir, ndir)
					log.Println("拷贝文件完成")
					commit(ndir, logs[0])
					log.Println("commit 完成")
					push(ndir, dusername, dtoken, branch.Str)
					log.Println("push 完成")
				}
			}
		}
		log.Println(IsExist(odir))
		if IsExist(odir) {
			err := os.RemoveAll(odir)
			if err != nil {
				log.Panicf("移除目录失败：%v", err)
			}
		}
		log.Println(IsExist(ndir))
		if IsExist(ndir) {
			err := os.RemoveAll(ndir)
			if err != nil {
				log.Panicf("移除目录失败：%v", err)
			}
		}
	}
}

func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

func chageUrl(url, key string) string {
	var newUrl string
	if strings.Contains(url, "http://localhost") {
		change_url := confyaml.Get(key).(string)
		nstr := strings.TrimPrefix(url, "http://localhost")
		newUrl = change_url + nstr
	} else {
		newUrl = url
	}
	return newUrl
}
