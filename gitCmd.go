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
	"log"
	"os"
	"strings"
)

func gitclone(url, dir, username, token string) {
	fmt.Printf("git clone %s %s --recursive", url, dir)
	r, err := git.PlainClone(dir, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: username,
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
func commit(path, msg string) {
	r, err := git.PlainOpen(path)
	CheckIfError(err)
	tree, enderr := r.Worktree()
	if enderr != nil {
		log.Printf("获取git tree失败:%v", enderr)
	}
	enderr = tree.AddWithOptions(&git.AddOptions{
		All: true,
	})
	if enderr != nil {
		log.Printf("git add失败:%v", enderr)
	}
	treeCommit, endError := tree.Commit(msg, &git.CommitOptions{All: true})
	if endError != nil {
		log.Panicf("设置commitc参数失败:%v", endError)
	}

	_, endError = r.CommitObject(treeCommit)
	if endError != nil {
		log.Panicf("执行git commit失败%v", endError)
	}

}

func push(path, username, token, bath string) {
	r, err := git.PlainOpen(path)
	CheckIfError(err)
	barchInfo := fmt.Sprintf("+refs/heads/%s:refs/heads/%s", bath, bath)
	po := &git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs:   []config.RefSpec{config.RefSpec(barchInfo)},
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
