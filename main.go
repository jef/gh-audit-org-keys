package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

const (
	githubURL     = "https://github.com"
	githubOrgsAPI = "https://api.github.com/orgs"
)

var (
	githubOrg = os.Getenv("GITHUB_ORGANIZATION")
	githubPat = os.Getenv("GITHUB_PAT")
)

type member struct {
	Login string `json:"login"`
}

func main() {
	fmt.Println("getting members")
	members := getMembers()

	fmt.Println("getting keys")

	getKeys(members)
}

func getKeys(members []member) {
	var membersWithNoKey []member
	var wg sync.WaitGroup
	client := &http.Client{}
	for _, member := range members {
		wg.Add(1)
		member := member

		go func() {
			defer wg.Done()

			req, err := http.NewRequest(
				"GET",
				fmt.Sprintf("%s/%s.keys", githubURL, member.Login),
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
			} else {
				membersWithNoKey = append(membersWithNoKey, member)
			}
		}()
	}
	wg.Wait()

	fmt.Println(fmt.Sprintf("users with no key (%d):", len(membersWithNoKey)))
	for _, member := range membersWithNoKey {
		wg.Add(1)
		member := member
		go func() {
			defer wg.Done()
			fmt.Println(member.Login)
		}()
	}
	wg.Wait()
}

func getMembers() []member {
	page := 1

	var members []member

	for {
		client := &http.Client{}

		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("%s/%s/members?filter=all&page=%d", githubOrgsAPI, githubOrg, page),
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

		var ms []member

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
