// Copyright 2016 Muhammad Shulhan <ms@kilabit.info>. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package j2p

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

const (
	// JiraAPIProject define project API for JIRA
	JiraAPIProject = "/rest/api/2/project"

	// JiraQIssueAll define query for all issue based on a project.
	JiraQIssueAll = "project = \"%s\" ORDER BY created ASC"
)

//
// NewJiraClient will create HTTP client to connect to JIRA server and do
// authentication using basic method (user,password).
//
func NewJiraClient(cfg *Config) (jiraCl *jira.Client, e error) {
	var res bool

	if DEBUG >= 1 {
		fmt.Printf("[j2p] NewJiraClient >> %s:%s@%s\n", cfg.Jira.User,
			cfg.Jira.Pass, cfg.Jira.URL)
	}

	jiraCl, e = jira.NewClient(nil, cfg.Jira.URL)

	if e != nil {
		return nil, e
	}

	res, e = jiraCl.Authentication.AcquireSessionCookie(cfg.Jira.User,
		cfg.Jira.Pass)

	if e != nil || !res {
		return nil, e
	}

	return jiraCl, nil
}

//
// JiraGetProjects will query all project in JIRA filtered by command parameters
// `-projects` and return it.
//
func (cmd *Cmd) JiraGetProjects() (
	jiraProjects *[]jira.Project,
	e error,
) {
	req, _ := cmd.JiraCl.NewRequest("GET", JiraAPIProject, nil)

	allProjects := new([]jira.Project)

	_, e = cmd.JiraCl.Do(req, allProjects)
	if e != nil {
		return nil, e
	}

	if len(cmd.Args.Projects) == 0 {
		jiraProjects = allProjects
		goto out
	}

	jiraProjects = new([]jira.Project)

	for _, v := range cmd.Args.Projects {
		for _, project := range *allProjects {
			if project.Name == v {
				*jiraProjects = append(*jiraProjects, project)
				break
			}
		}
	}

out:
	if DEBUG >= 2 {
		for x, project := range *jiraProjects {
			fmt.Printf("[j2p] JiraGetProjects %d >> %s: %s\n", x,
				project.Key, project.Name)
		}
	}

	return jiraProjects, nil
}
