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
	"sync/atomic"
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

	var (
		keyDsaSize              uint32
		keyEddsaSize            uint32
		keyEd25519Size          uint32
		keyRsaSize              uint32
		userDsaSize             uint32
		userEddsaSize           uint32
		userEd25519Size         uint32
		userRsaSize             uint32
		userWithKeySize         uint32
		userWithoutKeySize      uint32
		userWithMultipleKeySize uint32
		totalKeySize            uint32
		totalUserSize           = len(ms)
	)

	for _, m := range ms {
		wg.Add(1)
		m := m

		go func() {
			defer wg.Done()

			var (
				hasDsa     bool
				hasRsa     bool
				hasEddsa   bool
				hasEd25519 bool
			)

			for _, key := range m.Keys {
				atomic.AddUint32(&totalKeySize, 1)

				switch {
				case strings.Contains(key, "ssh-dsa"):
					atomic.AddUint32(&keyDsaSize, 1)
					hasDsa = true
				case strings.Contains(key, "ssh-rsa"):
					atomic.AddUint32(&keyRsaSize, 1)
					hasRsa = true
				case strings.Contains(key, "ssh-eddsa"):
					atomic.AddUint32(&keyEddsaSize, 1)
					hasEddsa = true
				case strings.Contains(key, "ssh-ed25519"):
					atomic.AddUint32(&keyEd25519Size, 1)
					hasEd25519 = true
				}
			}

			switch {
			case hasDsa:
				atomic.AddUint32(&userDsaSize, 1)
			case hasRsa:
				atomic.AddUint32(&userRsaSize, 1)
			case hasEddsa:
				atomic.AddUint32(&userEddsaSize, 1)
			case hasEd25519:
				atomic.AddUint32(&userEd25519Size, 1)
			}

			if len(m.Keys) == 0 {
				atomic.AddUint32(&userWithoutKeySize, 1)
				if *showUsers == "without" || *showUsers == "all" {
					zap.S().Infow("retrieved keys",
						"user", m.Login,
						"keys", m.Keys,
					)
				}
			}

			if len(m.Keys) > 0 {
				atomic.AddUint32(&userWithKeySize, 1)
				if *showUsers == "with" || *showUsers == "all" {
					zap.S().Infow("retrieved keys",
						"user", m.Login,
						"keys", m.Keys,
					)
				}
			}

			if len(m.Keys) > 1 {
				atomic.AddUint32(&userWithMultipleKeySize, 1)
				if *showUsers == "multiple" || *showUsers == "all" {
					zap.S().Infow("retrieved keys",
						"user", m.Login,
						"keys", m.Keys,
					)
				}
			}
		}()
	}
	wg.Wait()

	withKey := [][]string{
		{"users with keys", "DSA",
			fmt.Sprintf("%d (%.2f%%)", keyDsaSize, float32(keyDsaSize)/float32(totalKeySize)*100),
			fmt.Sprintf("%d (%.2f%%)", userDsaSize, float32(userDsaSize)/float32(totalUserSize)*100)},
		{"", "RSA",
			fmt.Sprintf("%d (%.2f%%)", keyRsaSize, float32(keyRsaSize)/float32(totalKeySize)*100),
			fmt.Sprintf("%d (%.2f%%)", userRsaSize, float32(userRsaSize)/float32(totalUserSize)*100)},
		{"", "ECDSA",
			fmt.Sprintf("%d (%.2f%%)", keyEddsaSize, float32(keyEddsaSize)/float32(totalKeySize)*100),
			fmt.Sprintf("%d (%.2f%%)", userEddsaSize, float32(userEddsaSize)/float32(totalUserSize)*100)},
		{"", "Ed25519",
			fmt.Sprintf("%d (%.2f%%)", keyEd25519Size, float32(keyEd25519Size)/float32(totalKeySize)*100),
			fmt.Sprintf("%d (%.2f%%)", userEd25519Size, float32(userEd25519Size)/float32(totalUserSize)*100)},
	}

	withoutKey := [][]string{
		{"users without keys", "", "", fmt.Sprintf("%d (%.2f%%)", userWithoutKeySize, float32(userWithoutKeySize)/float32(totalUserSize)*100)},
	}

	withMultipleKey := [][]string{
		{"users with multiple keys", "", "", fmt.Sprintf("%d (%.2f%%)", userWithMultipleKeySize, float32(userWithMultipleKeySize)/float32(totalUserSize)*100)},
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"description", "key type", "# of keys", "# of users"})
	t.SetHeaderColor(tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
	)
	t.SetFooter([]string{"", "total", fmt.Sprintf("%d", totalKeySize), fmt.Sprintf("%d", totalUserSize)})
	t.SetFooterColor(tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
	)
	t.SetRowLine(true)
	t.AppendBulk(withKey)
	t.AppendBulk(withoutKey)
	t.AppendBulk(withMultipleKey)
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
	p := 1

	var members []member

	for {
		client := &http.Client{}

		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("%s/%s/members?filter=all&page=%d", gitHubOrgAPI, gitHubOrg, p),
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
			p++
		} else {
			break
		}
	}

	return members
}
