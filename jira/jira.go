package jira

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

// client should be a singleton for now
var jiraClient *jira.Client

//GetClient returns an authenticated Jira client
func GetClient(instance string, username string, password string) *jira.Client {
	if jiraClient == nil {
		newClient, err := jira.NewClient(nil, instance)
		if err != nil {
			panic(err)
		} else {
			jiraClient = newClient
		}
	}

	if jiraClient.Authentication.Authenticated() == false {
		res, err := jiraClient.Authentication.AcquireSessionCookie(username, password)
		if err != nil || res == false {
			fmt.Printf("Result: %v\n", res)
			panic(err)
		}
	}

	return jiraClient
}
