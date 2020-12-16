package main

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"io"
	"io/ioutil"
	"log"
	chttp "net/http"
	"os"
	"strings"
)

const (
	UTF8    = string("UTF-8")
	GB18030 = string("GB18030")
)

func gitclone(url, dir, username, token string) {
	fmt.Printf("git clone %s %s --recursive", url, dir)
	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: username, // yes, this can be anything except an empty string
			Password: token,
		},
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Panicf("克隆仓库失败%v", err)
	}
	res, err := r.Head()
	if err != nil {
		log.Panicf("检索头部指向的分支失败%v", err)
	}
	commit, err := r.CommitObject(res.Hash())
	if err != nil {
		log.Printf("读取命令返回失败%v", err)
	}
	fmt.Println(commit)
}

func pull(path, token, username, barnch string) {

	// We instantiate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(path)
	CheckIfError(err)
	log.Println("路径是%s", path)
	log.Println("token是%s", token)
	log.Println("username是%s", username)

	// Get the working directory for the repository
	w, err := r.Worktree()
	CheckIfError(err)

	// Pull the latest changes from the origin remote and merge into the current branch
	Info("git pull origin")
	barch := plumbing.NewBranchReferenceName(barnch)
	err = w.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: barch,
		Auth: &http.BasicAuth{
			Username: username,
			Password: token,
		},
	})
	log.Println(err)
	if err.Error() != "already up-to-date" {
		CheckIfError(err)
	}

	// Print the latest commit that was just pulled
	ref, err := r.Head()
	CheckIfError(err)
	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)

	fmt.Println(commit)
}

func push(path, token, username, branch string) {
	r, err := git.PlainOpen(path)
	CheckIfError(err)

	Info("git push")
	// push using default options
	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/heads/master:refs/heads/master")},
		Auth: &http.BasicAuth{
			Username: username,
			Password: token,
		},
	})
	log.Println(err)
	if err.Error() != "already up-to-date" {
		CheckIfError(err)
	}

}

func gitlog(path string) []string {
	var logs []string
	r, err := git.PlainOpen(path)
	CheckIfError(err)

	// Gets the HEAD history from HEAD, just like this command:
	Info("git log")

	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	CheckIfError(err)

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), All: true})
	CheckIfError(err)

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c.Message)
		logInfo := bufio.NewReader(strings.NewReader(c.Message))
		for {
			logStr, err := logInfo.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				} else {
					log.Panicln("读取字符串失败:%v", err)
				}
				if len(logStr) == 0 || logStr == "\r\n" {
					continue
				}
			}
			logStr = strings.Trim(logStr, "\n")
			logs = append(logs, logStr)
		}
		return nil
	})
	CheckIfError(err)
	return logs
}

func commit(path, mes string) {
	fmt.Println("git log commit")
	r, err := git.PlainOpen(path)
	CheckIfError(err)
	w, err := r.Worktree()
	CheckIfError(err)

	// commit Since version 5.0.1, we can omit the Author signature, being read
	// from the git config files.
	commit, err := w.Commit(mes, &git.CommitOptions{})

	CheckIfError(err)

	// Prints the current HEAD to verify that all worked well.
	Info("git show -s")
	obj, err := r.CommitObject(commit)
	CheckIfError(err)

	fmt.Println(obj)
}

/*func push(path, branch string) {
	cmd := fmt.Sprintf("cd %s && git push origin master:%s", path, branch)
	execCommand(cmd)
}

func execCommand(commandName string) bool {

	//执行命令
	cmd := exec.Command(commandName)

	//显示运行的命令
	fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()
	errReader, errr := cmd.StderrPipe()

	if errr != nil {
		fmt.Println("err:" + errr.Error())
	}

	//开启错误处理
	go handlerErr(errReader)

	if err != nil {
		fmt.Println(err)
		return false
	}

	cmd.Start()
	in := bufio.NewScanner(stdout)
	for in.Scan() {
		cmdRe := ConvertByte2String(in.Bytes(), "GB18030")
		fmt.Println(cmdRe)
	}

	cmd.Wait()
	cmd.Wait()
	return true
}

//开启一个协程来错误
func handlerErr(errReader io.ReadCloser) {
	in := bufio.NewScanner(errReader)
	for in.Scan() {
		cmdRe := ConvertByte2String(in.Bytes(), "GB18030")
		fmt.Errorf(cmdRe)
	}
}

//对字符进行转码
func ConvertByte2String(byte []byte, charset string) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}*/

func Newbranch(path, username, token string) {
	r, err := git.PlainOpen(path)
	CheckIfError(err)
	Info("git branch my-test")
	headRef, err := r.Head()
	CheckIfError(err)
	ref := plumbing.NewHashReference("refs/heads/my-test", headRef.Hash())
	err = r.Storer.SetReference(ref)
	CheckIfError(err)

	po := &git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/heads/my-test:refs/heads/my-test")},
		Auth: &http.BasicAuth{
			Username: username,
			Password: token,
		},
	}
	err = r.Push(po)
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Print("origin remote was up to date, no push done")
		}
		log.Printf("push to remote origin error: %s", err)
	}
}

func GetOrginBranch(url, token string) {
	client := &chttp.Client{}
	req, err := chttp.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("创建http客户端失败：%s", err)
	}
	req.Header.Set("PRIVATE-TOKEN", token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("调用api失败:%v", err)
	}
	defer resp.Body.Close()
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))

	req2, err := chttp.NewRequest("GET", "http://10.0.0.106/api/v4/projects/1/repository/branches", nil)
	if err != nil {
		log.Fatalf("创建http客户端失败：%s", err)
	}
	req2.Header.Set("PRIVATE-TOKEN", token)
	resp2, err := client.Do(req2)
	if err != nil {
		log.Fatalf("调用api失败:%v", err)
	}
	defer resp2.Body.Close()
	body2, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		log.Fatalf("读取body失败:%v", err)
	}
	fmt.Println(string(body2))
}

func checkout(path, commit string) {
	r, err := git.PlainOpen(path)
	CheckIfError(err)
	Info("git show-ref --head HEAD")
	ref, err := r.Head()
	CheckIfError(err)
	fmt.Println(ref.Hash())

	w, err := r.Worktree()
	CheckIfError(err)

	Info("git checkout %s", commit)
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(commit),
	})
	CheckIfError(err)

	Info("git show-ref --head HEAD")
	ref, err = r.Head()
	CheckIfError(err)
	fmt.Println(ref.Hash())
}
