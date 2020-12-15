package main

import (
	"github.com/xanzy/go-gitlab"
	"log"
)

func gitlabCl(surl, stoken, durl, dtoken string) (scr, dest *gitlab.Client) {
	scr, err := gitlab.NewClient(stoken, gitlab.WithBaseURL(surl))
	if err != nil {
		log.Fatalf("创建源gitlab客户端: %v", err)
	}

	dest, err = gitlab.NewClient(dtoken, gitlab.WithBaseURL(durl))
	if err != nil {
		log.Fatalf("创建目标gitlab客户端: %v", err)
	}
	return
}
