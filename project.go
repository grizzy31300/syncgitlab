package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type pdata struct {
	Id                                               int                             `json:"id"`
	Description                                      string                          `json:"description"`
	Name                                             string                          `json:"name"`
	Name_with_namespace                              string                          `json:"name_with_namespace"`
	Path                                             string                          `json:"path"`
	Path_with_namespace                              string                          `json:"path_with_namespace"`
	Created_at                                       string                          `json:"created_at"`
	Default_branch                                   string                          `json:"default_branch"`
	Tag_list                                         []string                        `json:"tag_list"`
	Ssh_url_to_repo                                  string                          `json:"ssh_url_to_repo"`
	Http_url_to_repo                                 string                          `json:"http_url_to_repo"`
	Web_url                                          string                          `json:"web_url"`
	Readme_url                                       string                          `json:"readme_url"`
	Avatar_url                                       string                          `json:"avatar_url"`
	Forks_count                                      int                             `json:"forks_count"`
	Star_count                                       int                             `json:"star_count"`
	Last_activity_at                                 string                          `json:"last_activity_at"`
	Namespace                                        namespaceInfo                   `json:"namespace"`
	Links                                            linksInfo                       `json:"_links"`
	Packages_enabled                                 bool                            `json:"packages_enabled"`
	Empty_repo                                       bool                            `json:"empty_repo"`
	Archived                                         bool                            `json:"archived"`
	Visibility                                       string                          `json:"visibility"`
	Resolve_outdated_diff_discussions                bool                            `json:"resolve_outdated_diff_discussions"`
	Container_registry_enabled                       bool                            `json:"container_registry_enabled"`
	Container_expiration_policy                      container_expiration_policyInfo `json:"container_expiration_policy"`
	Issues_enabled                                   bool                            `json:"issues_enabled"`
	Merge_requests_enabled                           bool                            `json:"merge_requests_enabled"`
	Wiki_enabled                                     bool                            `json:"wiki_enabled"`
	Jobs_enabled                                     bool                            `json:"jobs_enabled"`
	Snippets_enabled                                 bool                            `json:"snippets_enabled"`
	Service_desk_enabled                             bool                            `json:"service_desk_enabled"`
	Service_desk_address                             string                          `json:"service_desk_address"`
	Can_create_merge_request_in                      bool                            `json:"can_create_merge_request_in"`
	Issues_access_level                              string                          `json:"issues_access_level"`
	Repository_access_level                          string                          `json:"repository_access_level"`
	Merge_requests_access_level                      string                          `json:"merge_requests_access_level"`
	Forking_access_level                             string                          `json:"forking_access_level"`
	Wiki_access_level                                string                          `json:"wiki_access_level"`
	Builds_access_level                              string                          `json:"builds_access_level"`
	Snippets_access_level                            string                          `json:"snippets_access_level"`
	Pages_access_level                               string                          `json:"pages_access_level"`
	Emails_disabled                                  string                          `json:"emails_disabled"`
	Shared_runners_enabled                           bool                            `json:"shared_runners_enabled"`
	Lfs_enabled                                      bool                            `json:"lfs_enabled"`
	Creator_id                                       int                             `json:"creator_id"`
	Import_status                                    string                          `json:"import_status"`
	Open_issues_count                                int                             `json:"open_issues_count"`
	Ci_default_git_depth                             int                             `json:"ci_default_git_depth"`
	Public_jobs                                      bool                            `json:"public_jobs"`
	Build_timeout                                    int                             `json:"build_timeout"`
	Auto_cancel_pending_pipelines                    string                          `json:"auto_cancel_pending_pipelines"`
	Build_coverage_regex                             string                          `json:"build_coverage_regex"`
	Ci_config_path                                   string                          `json:"ci_config_path"`
	Shared_with_groups                               []string                        `json:"shared_with_groups"`
	Only_allow_merge_if_pipeline_succeeds            bool                            `json:"only_allow_merge_if_pipeline_succeeds"`
	Allow_merge_on_skipped_pipeline                  string                          `json:"allow_merge_on_skipped_pipeline"`
	Request_access_enabled                           bool                            `json:"request_access_enabled"`
	Only_allow_merge_if_all_discussions_are_resolved bool                            `json:"only_allow_merge_if_all_discussions_are_resolved"`
	Remove_source_branch_after_merge                 bool                            `json:"remove_source_branch_after_merge"`
	Printing_merge_request_link_enabled              bool                            `json:"printing_merge_request_link_enabled"`
	Merge_method                                     string                          `json:"merge_method"`
	Suggestion_commit_message                        string                          `json:"suggestion_commit_message"`
	Auto_devops_enabled                              bool                            `json:"auto_devops_enabled"`
	Auto_devops_deploy_strategy                      string                          `json:"auto_devops_deploy_strategy"`
	Autoclose_referenced_issues                      bool                            `json:"autoclose_referenced_issues"`
	Repository_storage                               string                          `json:"repository_storage"`
	Permissions                                      permissionsInfo                 `json:"permissions"`
}

type namespaceInfo struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Kind       string `json:"kind"`
	Full_path  string `json:"full_path"`
	Parent_id  int    `json:"parent_id"`
	Avatar_url string `json:"avatar_url"`
	Web_url    string `json:"web_url"`
}

type linksInfo struct {
	Self           string `json:"self"`
	Issues         string `json:"issues"`
	Merge_requests string `json:"merge_requests"`
	Repo_branches  string `json:"repo_branches"`
	Labels         string `json:"labels"`
	Events         string `json:"events"`
	Members        string `json:"members"`
}

type container_expiration_policyInfo struct {
	Cadence         string `json:"cadence"`
	Enabled         bool   `json:"enabled"`
	Keep_n          int    `json:"keep_n"`
	Older_than      string `json:"older_than"`
	Name_regex      string `json:"name_regex"`
	Name_regex_keep string `json:"name_regex_keep"`
	Next_run_at     string `json:"next_run_at"`
}

type permissionsInfo struct {
	Project_access project_accessInfo `json:"project_access"`
	Group_access   group_accessInfo   `json:"group_access"`
}

type group_accessInfo struct {
	Access_level       int `json:"access_level"`
	Notification_level int `json:"notification_level"`
}

type project_accessInfo struct {
	Access_level       int `json:"access_level"`
	Notification_level int `json:"notification_level"`
}

type createpOptions struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Visibility   string `json:"visibility"`
	Namespace_id string `json:"namespace_id,omitempty"`
}

type newProject struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func listProject(url, token string) []pdata {
	url = url + "/api/v4/projects"
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
	bs := pujson(body)
	return bs
}

func duplicate(bs []pdata) []pdata {
	log.Println("#################################")
	log.Println(len(bs))

	for _, v := range bs {
		log.Printf("主ID:%d,子ID:%d\n", v.Id, v.Namespace.Parent_id)
		if v.Namespace.Parent_id != 0 {
			for nk, nv := range bs {
				log.Printf("主ID:%d\n", nv.Id)
				if nv.Id == v.Namespace.Parent_id {
					bs = append(bs[:nk], bs[nk+1:]...)
				}
			}
		}
	}
	log.Println(len(bs))
	log.Println("#################################")
	return bs
}

func pujson(input []byte) []pdata {
	var dinfo []pdata
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(input, &dinfo)
	if err != nil {
		log.Printf("josn解析失败%v", err)
	}
	return dinfo
}

func createProject(listProject []pdata, url, token string) {
	var Options createpOptions
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for _, v := range listProject {
		if v.Name == "Monitoring" {
			continue
		}
		log.Printf("name:%s,fullname:%s", v.Name, v.Namespace.Full_path)
		Options.Name = v.Name
		Options.Visibility = v.Visibility
		Options.Description = v.Description
		Options.Namespace_id = strconv.Itoa(v.Namespace.Id)
		log.Printf("OptionsNamespace_id:%s", Options.Namespace_id)
		createOption, err := json.Marshal(Options)
		if err != nil {
			log.Panicf("序列化请求数据失败：%v", err)
		}
		log.Printf("创建项目的参数：createOption:%v", Options)
		httpCreateP(url, token, createOption)
	}
}

func httpCreateP(url, token string, input []byte) newProject {
	var newp newProject
	url = url + "/api/v4/projects"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(input)))
	if err != nil {
		log.Fatalf("创建http客户端失败：%s", err)
	}
	req.Header.Set("PRIVATE-TOKEN", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("调用创建Project api失败:%v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("读取创建Project返回信息失败:%v", err)
	}
	name := gjson.GetBytes(body, "message.name")
	if name.String() != "[\"has already been taken\"]" {
		newp.Id = gjson.GetBytes(body, "id").String()
		newp.Name = gjson.GetBytes(body, "name").String()
	}
	log.Printf("NewProject的参数%v", newp)
	return newp
}

func Transfer(url, token string, newmap map[string][]string) {
	var urlp string
	for k, newpids := range newmap {
		if len(newpids) == 0 {
			continue
		}
		log.Println("key的值:%s", k)
		log.Printf("获取grop信息%s", newmap[k])
		for _, pid := range newpids {
			urlp = fmt.Sprintf("%s/api/v4/groups/%s/projects/%s", url, k, pid)
			log.Println(urlp)
			client := &http.Client{}
			req, err := http.NewRequest("POST", urlp, nil)
			if err != nil {
				log.Fatalf("创建http客户端失败：%s", err)
			}
			req.Header.Set("PRIVATE-TOKEN", token)
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("调用关联组和项目api失败:%v", err)
			}
			defer resp.Body.Close()
			_, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Panicf("读取创建Project返回信息失败:%v", err)
			}
		}
	}
}
