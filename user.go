package main

type UserInfo struct {
	Id         int
	Name       string
	State      string
	Avatar_url string
}

/*func listuser(url, token string) {
	url = url + "/api/v4/user"
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

func uujson() {

}*/
