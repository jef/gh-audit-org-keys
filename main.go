package main

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

type member struct {
	Login string `json:"login"`
	Keys  []string
}

func main() {
	ms := getMembers()
	ms = getKeys(ms)
	printReport(ms)
}

func printReport(ms []member) {
	var wg sync.WaitGroup

	withSize := 0
	withoutSize := 0
	multipleSize := 0

	for _, m := range ms {
		wg.Add(1)
		m := m
		go func() {
			defer wg.Done()
			if len(m.Keys) == 0 {
				withoutSize++
				if *showUsers == "without" || *showUsers == "all" {
					zap.S().Infow("retrieved keys",
						"user", m.Login,
						"keys", m.Keys,
					)
				}
			}

			if len(m.Keys) != 0 {
				withSize++
				if *showUsers == "with" || *showUsers == "all" {
					zap.S().Infow("retrieved keys",
						"user", m.Login,
						"keys", m.Keys,
					)
				}
			}

			if len(m.Keys) > 1 {
				multipleSize++
				if *showUsers == "multiple" || *showUsers == "all" {
					zap.S().Infow("retrieved keys",
						"user", m.Login,
						"keys", m.Keys,
					)
				}
			}
			// todo strong and weak keys
		}()
	}
	wg.Wait()

	d := [][]string{
		{"users with keys", fmt.Sprintf("%d (%.2f%%)", withSize,
			float32(withSize)/float32(len(ms)) * 100)},
		{"users without keys", fmt.Sprintf("%d (%.2f%%)", withoutSize,
			float32(withoutSize)/float32(len(ms)) * 100)},
		{"users with multiple keys", fmt.Sprintf("%d (%.2f%%)", multipleSize,
			float32(multipleSize)/float32(len(ms)) * 100)},
		// todo: calculate bit length of keys
		//{"users with strong keys", fmt.Sprintf("%d", 0)},
		//{"users with weak keys", fmt.Sprintf("%d", 0)},
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"description", "# of users"})
	t.SetHeaderColor(tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
	)
	t.SetFooter([]string{"total users", fmt.Sprintf("%d", len(ms))})

	t.AppendBulk(d)
	t.Render()
}

func getKeys(ms []member) []member {
	var wg sync.WaitGroup
	client := &http.Client{}
	for i := 0; i < len(ms); i++ {
		wg.Add(1)
		i := i

		go func() {
			m := &ms[i]
			defer wg.Done()

			req, err := http.NewRequest(
				"GET",
				fmt.Sprintf("%s/%s.keys", gitHubURL, m.Login),
				nil,
			)
			if err != nil {
				zap.S().Fatal(err)
			}

			res, err := client.Do(req)
			if err != nil {
				zap.S().Fatal(err)
			}

			defer res.Body.Close()

			key, err := ioutil.ReadAll(res.Body)
			if err != nil {
				zap.S().Fatal(err)
			} else if (len(key)) != 0 {
				m.Keys = strings.Split(strings.TrimSpace(string(key)), "\n")
			}
		}()
	}
	wg.Wait()

	return ms
}

func getMembers() []member {
	page := 1

	var members []member

	for {
		client := &http.Client{}

		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("%s/%s/members?filter=all&page=%d", gitHubOrgAPI, gitHubOrg, page),
			nil,
		)
		if err != nil {
			zap.S().Fatal(err)
		}
		req.Header.Add("authorization", fmt.Sprintf("token %s", gitHubPAT))

		res, err := client.Do(req)
		if err != nil {
			zap.S().Fatal(err)
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			zap.S().Fatal(err)
		}

		var ms []member

		err = json.Unmarshal(body, &ms)
		if err != nil {
			zap.S().Fatal(err)
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
