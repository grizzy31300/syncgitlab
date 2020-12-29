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
	log.Printf("git clone %s %s --recursive", url, dir)
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
	log.Println(commit)
}

func pull(path, token, username, barch string) {
	// We instantiate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(path)
	CheckIfError(err)
	log.Printf("路径是%s", path)
	log.Printf("token是%s", token)
	log.Printf("username是%s", username)

	// Get the working directory for the repository
	w, err := r.Worktree()
	CheckIfError(err)

	// Pull the latest changes from the origin remote and merge into the current branch
	Info("git pull origin")
	Rbarch := plumbing.NewBranchReferenceName(barch)
	err = w.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: Rbarch,
		Auth: &http.BasicAuth{
			Username: username,
			Password: token,
		},
		Force: true,
	})
	if err != nil {
		if err.Error() != "already up-to-date" {
			log.Printf("执行git push失败:%v", err)
		}
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
	log.Printf("+refs/heads/%s:refs/heads/%s\n", bath, bath)
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

func checkout(path, commit string) {
	bc := false
	if commit != "master" {
		bc = true
	}
	r, err := git.PlainOpen(path)
	CheckIfError(err)
	// ... retrieving the commit being pointed by HEAD
	Info("git show-ref --head HEAD")
	ref, err := r.Head()
	CheckIfError(err)
	fmt.Println(ref.Hash())

	w, err := r.Worktree()
	CheckIfError(err)

	// ... checking out to commit
	Info("git checkout %s", commit)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(commit),
		Create: bc,
		Force:  true,
	})
	CheckIfError(err)

	// ... retrieving the commit being pointed by HEAD, it shows that the
	// repository is pointing to the giving commit in detached mode
	Info("git show-ref --head HEAD")
	ref, err = r.Head()
	CheckIfError(err)
}

func track(path, barch string) {
	r, err := git.PlainOpen(path)
	CheckIfError(err)

	var name, remote, remoteBranch = barch, "origin", barch

	var remoteRef = plumbing.NewRemoteReferenceName(remote, remoteBranch)
	var ref, _ = r.Reference(remoteRef, true)

	var mergeRef = plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", remoteBranch))
	_ = r.CreateBranch(&config.Branch{Name: name, Remote: remote, Merge: mergeRef})

	var localRef = plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", name))
	_ = r.Storer.SetReference(plumbing.NewHashReference(localRef, ref.Hash()))
}

func checkoutBranch(path, commit string) {

	r, err := git.PlainOpen(path)
	CheckIfError(err)

	w, err := r.Worktree()
	CheckIfError(err)

	// ... checking out to commit
	Info("git checkout %s", commit)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(commit),
		Create: false,
		Force:  true,
	})
	CheckIfError(err)

}
