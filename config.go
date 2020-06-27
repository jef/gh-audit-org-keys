package main

import (
	"flag"
	"os"
)

const (
	gitHubURL    = "https://github.com"
	gitHubOrgAPI = "https://api.github.com/orgs"
)

var (
	showUsers = flag.String("show-users", "", "display users with filter (`all`, `with`, `without`, `multiple`)")
	gitHubOrg = os.Getenv("GITHUB_ORGANIZATION")
	gitHubPAT = os.Getenv("GITHUB_PAT")
)

func init() {
	flag.Parse()
	initLogger()
}
