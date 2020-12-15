package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"log"
	"os"
)

func syncCode(nproject, oproject []pdata, stoken, susername, dtoken, dusername string, newpmap map[string][]gjson.Result) {
	for _, v := range oproject {
		if v.Name == "Monitoring" {
			continue
		}
		odir := fmt.Sprintf("./old%s", v.Name)
		log.Println(odir, IsExist(odir))
		if !IsExist(odir) {
			os.Mkdir(odir, 0755)
		}
		ndir := fmt.Sprintf("./new%s", v.Name)
		if !IsExist(ndir) {
			os.Mkdir(ndir, 0755)
		}
		log.Printf("克隆的url:%s", v.Http_url_to_repo)
		gitclone(v.Http_url_to_repo, odir, susername, stoken)
		log.Println("###########################################################")
		gitlog(odir)
		log.Println("###########################################################")
		for _, branch := range newpmap[v.Name_with_namespace] {
			log.Printf("%s pull的分支%s", v.Name, branch.Str)
			pull(odir, stoken, susername, branch.Str)
			for _, np := range nproject {
				if np.Name == "Monitoring" {
					continue
				}
				if v.Name == np.Name && v.Name_with_namespace == np.Name_with_namespace {
					/*ngit := fmt.Sprintf("%s/.git", ndir)
					ogit := fmt.Sprintf("%s/.git", odir)
					bgit := fmt.Sprintf("%s/.git_bak", odir)*/
					log.Printf("克隆的url:%s", np.Http_url_to_repo)
					gitclone(np.Http_url_to_repo, ndir, dusername, dtoken)
					/*os.Rename(ogit, bgit)
					os.Rename(ngit, ogit)*/
					log.Println(dtoken)
					log.Println(dusername)
					//push(odir, dtoken, dusername, branch.Str)
					push(odir, branch.Str)
				}
			}
		}
		/*log.Println(IsExist(odir))
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
		}*/
	}
}

func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}