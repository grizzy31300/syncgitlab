package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type cbarnch struct {
	Branch string `json:"branch"`
	Ref    string `json:"ref"`
}

func listBranch(projecdInfos []pdata, url, token string) map[string][]gjson.Result {
	var pb map[string][]gjson.Result
	pb = make(map[string][]gjson.Result)
	for _, v := range projecdInfos {
		if v.Name == "Monitoring" {
			continue
		}
		purl := fmt.Sprintf("%s/api/v4/projects/%s/repository/branches", url, strconv.Itoa(v.Id))
		log.Printf("分支的url:%s", purl)
		client := &http.Client{}
		req, err := http.NewRequest("GET", purl, nil)
		if err != nil {
			log.Fatalf("创建http客户端失败：%s", err)
		}
		req.Header.Set("PRIVATE-TOKEN", token)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("调用查看分支api失败:%v", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("读取返回值失败:%v", err)
		}
		names := gjson.GetBytes(body, "#.name")
		pb[v.Name_with_namespace] = names.Array()
	}
	return pb
}

func trancBranch(oldmap map[string][]gjson.Result, nplist []pdata) map[int][]gjson.Result {
	var nblist map[int][]gjson.Result
	nblist = make(map[int][]gjson.Result)
	for _, v := range nplist {
		_, ok := oldmap[v.Name_with_namespace]
		if ok {
			nblist[v.Id] = oldmap[v.Name_with_namespace]
		}
	}
	return nblist
}

func createBranch(tranB map[int][]gjson.Result, url, token string) {
	var createBranchOption cbarnch
	for k, v := range tranB {
		for _, bname := range v {
			burl := fmt.Sprintf("%s/api/v4/projects/%s/repository/branches", url, strconv.Itoa(k))
			createBranchOption.Branch = bname.Str
			createBranchOption.Ref = "master"
			hcj, err := json.Marshal(createBranchOption)
			if err != nil {
				log.Panicf("序列化创建branch失败：%v", err)
			}
			client := &http.Client{}
			req, err := http.NewRequest("POST", burl, strings.NewReader(string(hcj)))
			if err != nil {
				log.Fatalf("创建http客户端失败：%s", err)
			}
			req.Header.Set("PRIVATE-TOKEN", token)
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("调用api失败:%v", err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Panicf("读取创建Project返回信息失败:%v", err)
			}
			log.Println(string(body))
		}
	}
}
