package main

import "github.com/xanzy/go-gitlab"

func diffg(sgroups, dgroups []*gitlab.Group) []*gitlab.Group {
	var group []*gitlab.Group
	for _, sv := range sgroups {
		sf := *sv
		for _, dv := range dgroups {
			df := *dv
			if sf.FullPath == df.FullPath {
				continue
			} else {
				group = append(group, sv)
			}
		}
	}
	return group
}

func diffp(sprojects, dprojects []*gitlab.Project) []*gitlab.Project {
	var project []*gitlab.Project
	for _, sv := range sprojects {
		sf := *sv
		for _, dv := range dprojects {
			df := *dv
			if sf.Path == df.Path {
				continue
			} else {
				project = append(project, sv)
			}
		}
	}
	return project
}
