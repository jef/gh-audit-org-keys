package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	githubUrl     = "https://github.com"
	githubOrgsApi = "https://api.github.com/orgs"
)

var (
	githubOrg = os.Getenv("GITHUB_ORGANIZATION")
	githubPat = os.Getenv("GITHUB_PAT")
)

type Member struct {
	Login string `json:"login"`
}

func main() {
	fmt.Println("getting members")
	members := getMembers()

	fmt.Println("getting keys")
	getKeys(members)
}

func getKeys(members []Member) {
	client := &http.Client{}

	var membersWithNoKey []Member

	for _, member := range members {
		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("%s/%s.keys", githubUrl, member.Login),
			nil,
		)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Add("authorization", fmt.Sprintf("token %s", githubPat))

		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()

		key, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		if len(key) != 0 {
			fmt.Println(fmt.Sprintf("%s:\n%s", member.Login, key))
			fmt.Println("-------------------------------------------------------------------------------------")
			fmt.Println()
		} else {
			membersWithNoKey = append(membersWithNoKey, member)
		}
	}

	fmt.Println(fmt.Sprintf("members with no keys (%d):", len(membersWithNoKey)))
	for _, member := range membersWithNoKey {
		fmt.Println(fmt.Sprintf("%s", member.Login))
	}
}

func getMembers() []Member {
	page := 1

	var members []Member

	for {
		client := &http.Client{}

		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("%s/%s/members?filter=all&page=%d", githubOrgsApi, githubOrg, page),
			nil,
		)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Add("authorization", fmt.Sprintf("token %s", githubPat))

		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var ms []Member

		err = json.Unmarshal(body, &ms)
		if err != nil {
			log.Fatal(err)
		}

		if len(ms) != 0 {
			members = append(members, ms...)
			page++
		} else {
			break
		}
	}

	return members
}
