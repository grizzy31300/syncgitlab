package main

import (
	"github.com/spf13/viper"
	"log"
)

var (
	confyaml *viper.Viper
)

func main() {
	surl := confyaml.Get("src.url").(string)
	stoken := confyaml.Get("src.token").(string)
	susername := confyaml.Get("src.username").(string)
	durl := confyaml.Get("dest.url").(string)
	dtoken := confyaml.Get("dest.token").(string)
	dusername := confyaml.Get("dest.username").(string)
	log.Printf("原地址%s", surl)
	log.Printf("目标地址%s", durl)
	listg := listGroup(surl, stoken)
	createGroup(listg, durl, dtoken)
	listp := listProject(surl, stoken)
	//listp = duplicate(listp)
	oldmap := listGroupProject(listg, surl, stoken)
	nlistg := listGroup(durl, dtoken)
	nlistp := trans(oldmap, nlistg, listp)
	createProject(nlistp, durl, dtoken)
	oldBmap := listBranch(listp, surl, stoken)
	nblistProject := listProject(durl, dtoken)
	trb := trancBranch(oldBmap, nblistProject)
	log.Printf("新branch map:%v", trb)
	createBranch(trb, durl, dtoken)
	log.Printf("分支和project关系%v", oldBmap)
	syncCode(nblistProject, listp, stoken, susername, dtoken, dusername, oldBmap)
}
