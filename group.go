package main

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type data struct {
	Id                                int    `json:"id"`
	Web_url                           string `json:"web_url"`
	Name                              string `json:"name"`
	Path                              string `json:"path"`
	Description                       string `json:"description"`
	Visibility                        string `json:"visibility"`
	Share_with_group_lock             bool   `json:"share_with_group_lock"`
	Require_two_factor_authentication bool   `json:"require_two_factor_authentication"`
	Two_factor_grace_period           int    `json:"two_factor_grace_period"`
	Project_creation_level            string `json:"project_creation_level"`
	Auto_devops_enabled               string `json:"auto_devops_enabled"`
	Subgroup_creation_level           string `json:"subgroup_creation_level"`
	Emails_disabled                   string `json:"emails_disabled"`
	Mentions_disabled                 string `json:"mentions_disabled"`
	Lfs_enabled                       bool   `json:"lfs_enabled"`
	Default_branch_protection         int    `json:"default_branch_protection"`
	Avatar_url                        string `json:"avatar_url"`
	Request_access_enabled            bool   `json:"request_access_enabled"`
	Full_name                         string `json:"full_name"`
	Full_path                         string `json:"full_path"`
	Created_at                        string `json:"created_at"`
	Parent_id                         int    `json:"parent_id"`
}

type cdata struct {
	Path       string `json:"path"`
	Name       string `json:"name"`
	Parent_id  string `json:"parent_id,omitempty"`
	Visibility string `json:"visibility"`
}

type gp struct {
	Pnames     []gjson.Result
	Gname      string
	GFull_name string
}

func listGroup(url, token string) []data {
	url += "/api/v4/groups"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("创建http客户端失败：%s", err)
	}
	req.Header.Set("PRIVATE-TOKEN", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("调用api失败:%v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取返回值失败:%v", err)
	}
	bs := ujson(body)
	return bs
}

func listGroupProject(groups []data, url, token string) map[string]gp {
	var gafterp gp
	var gpinfo map[string]gp
	gpinfo = make(map[string]gp)
	for _, v := range groups {
		if v.Name == "GitLab Instance" {
			continue
		}
		urlg := url + "/api/v4/groups/" + strconv.Itoa(v.Id) + "/projects"
		log.Println(urlg)
		client := &http.Client{}
		req, err := http.NewRequest("GET", urlg, nil)
		if err != nil {
			log.Fatalf("创建http客户端失败：%s", err)
		}
		req.Header.Set("PRIVATE-TOKEN", token)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("调用创建Group获取project api失败:%v", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		names := gjson.GetBytes(body, "#.name")

		gafterp.GFull_name = v.Full_name
		gafterp.Gname = v.Name
		gafterp.Pnames = names.Array()
		gpinfo[v.Full_name] = gafterp
	}
	return gpinfo
}

func ujson(input []byte) []data {
	var dinfo []data
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(input, &dinfo)
	if err != nil {
		log.Printf("josn解析失败%v", err)
	}
	return dinfo
}

func createGroup(grouplist []data, url, token string) {
	var cdinfo cdata
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for _, v := range grouplist {
		cdinfo.Name = v.Name
		cdinfo.Path = v.Path
		cdinfo.Visibility = v.Visibility
		if v.Name == "GitLab Instance" {
			continue
		}
		if v.Parent_id == 0 {
			cdinfo.Parent_id = ""
		} else {
			var fullPath string
			var Parent_id int
			for _, v2 := range grouplist {
				if v2.Id == v.Parent_id {
					fullPath = v2.Full_path
					break
				}
			}
			glist := listGroup(url, token)
			for _, v3 := range glist {
				if v3.Full_path == fullPath {
					Parent_id = v3.Id
				}
			}
			if Parent_id == 0 {
				grouplist = append(grouplist, v)
				continue
			}
			cdinfo.Parent_id = strconv.Itoa(Parent_id)
		}
		hinfo, err := json.Marshal(cdinfo)
		if err != nil {
			log.Panicf("json序列化失败%v", err)
		}
		httpCreate(url, token, hinfo)
	}
}

func httpCreate(url, token string, input []byte) {
	url += "/api/v4/groups"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(input)))
	if err != nil {
		log.Fatalf("创建http客户端失败：%s", err)
	}
	req.Header.Set("PRIVATE-TOKEN", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("调用创建Group api失败:%v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
}

func trans(oldMap map[string]gp, newGroup []data, newProjects []pdata) []pdata {
	var cProjects []pdata
	for _, v := range newGroup {
		for _, pname := range oldMap[v.Full_name].Pnames {
			for _, np := range newProjects {
				if np.Name == pname.Str {
					log.Printf("before NamespaceId:%d", np.Namespace.Id)
					np.Namespace.Id = v.Id
					log.Printf("after NamespaceId:%d", np.Namespace.Id)
					cProjects = append(cProjects, np)
				} else {
					cProjects = append(cProjects, np)
				}
			}
		}
	}
	return cProjects
}
