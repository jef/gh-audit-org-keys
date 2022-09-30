package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/cli/go-gh"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

//go:embed version.txt
var Version string

func init() {
	cobra.OnInitialize()
}

type member struct {
	Login string `json:"login"`
	Keys  []string
}

type keyTable struct {
	keyDsaSize              uint32
	keyEcdsaSize            uint32
	keyEd25519Size          uint32
	keyRsaSize              uint32
	keyStrongRsaSize        uint32
	keyWeakRsaSize          uint32
	userDsaSize             uint32
	userEcdsaSize           uint32
	userEd25519Size         uint32
	userRsaSize             uint32
	userWithKeySize         uint32
	userWithoutKeySize      uint32
	userWithMultipleKeySize uint32
	keySize                 uint32
	strongKeySize           uint32
	weakKeySize             uint32
	userSize                uint32
	userWithStrongKeySize   uint32
	userWithWeakKeySize     uint32
	userWithWeakKey         []member
}

func printReport(kt keyTable) {
	withKey := [][]string{
		{"users with keys",
			"DSA",
			fmt.Sprintf("%d (%.2f%%)", kt.keyDsaSize, float32(kt.keyDsaSize)/float32(kt.keySize)*100),
			fmt.Sprintf("%d (%.2f%%)", kt.userDsaSize, float32(kt.userDsaSize)/float32(kt.userSize)*100)},
		{"",
			"ECDSA",
			fmt.Sprintf("%d (%.2f%%)", kt.keyEcdsaSize, float32(kt.keyEcdsaSize)/float32(kt.keySize)*100),
			fmt.Sprintf("%d (%.2f%%)", kt.userEcdsaSize, float32(kt.userEcdsaSize)/float32(kt.userSize)*100)},
		{"",
			"Ed25519",
			fmt.Sprintf("%d (%.2f%%)", kt.keyEd25519Size, float32(kt.keyEd25519Size)/float32(kt.keySize)*100),
			fmt.Sprintf("%d (%.2f%%)", kt.userEd25519Size, float32(kt.userEd25519Size)/float32(kt.userSize)*100)},
		{"",
			"RSA",
			fmt.Sprintf("%d (%.2f%%)", kt.keyRsaSize, float32(kt.keyRsaSize)/float32(kt.keySize)*100),
			fmt.Sprintf("%d (%.2f%%)", kt.userRsaSize, float32(kt.userRsaSize)/float32(kt.userSize)*100)},
	}

	withoutKey := [][]string{
		{"users without keys",
			"",
			"",
			fmt.Sprintf("%d (%.2f%%)", kt.userWithoutKeySize, float32(kt.userWithoutKeySize)/float32(kt.userSize)*100)},
	}

	withMultipleKey := [][]string{{"users with multiple keys",
		"",
		"",
		fmt.Sprintf("%d (%.2f%%)", kt.userWithMultipleKeySize, float32(kt.userWithMultipleKeySize)/float32(kt.userSize)*100)},
	}

	strongKey := [][]string{
		{"users with strong keys",
			"",
			fmt.Sprintf("%d (%.2f%%)", kt.strongKeySize, float32(kt.strongKeySize)/float32(kt.keySize)*100),
			fmt.Sprintf("%d (%.2f%%)", kt.userWithStrongKeySize, float32(kt.userWithStrongKeySize)/float32(kt.userWithKeySize)*100)},
	}

	weakKey := [][]string{
		{"users with weak keys",
			"",
			fmt.Sprintf("%d (%.2f%%)", kt.weakKeySize, float32(kt.weakKeySize)/float32(kt.keySize)*100),
			fmt.Sprintf("%d (%.2f%%)", kt.userWithWeakKeySize, float32(kt.userWithWeakKeySize)/float32(kt.userWithKeySize)*100)},
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"description", "key type", "# of keys", "# of users"})
	t.SetHeaderColor(tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgCyanColor})
	t.SetFooter([]string{"", "total", fmt.Sprintf("%d", kt.keySize), fmt.Sprintf("%d", kt.userSize)})
	t.SetFooterColor(tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgCyanColor})
	t.SetRowLine(true)
	t.AppendBulk(withKey)
	t.AppendBulk(withoutKey)
	t.AppendBulk(withMultipleKey)
	t.AppendBulk(strongKey)
	t.AppendBulk(weakKey)
	t.Render()

	if len(kt.userWithWeakKey) > 0 {
		log.Println("users with weak keys:")
		for _, m := range kt.userWithWeakKey {
			log.Printf("%s\n", m.Login)
		}
	}
}

func generateKeyTable(showUsers string, ms []member) keyTable {
	var kt keyTable
	kt.userSize = uint32(len(ms))
	var wg sync.WaitGroup

	for _, m := range ms {
		wg.Add(1)
		m := m

		go func() {
			defer wg.Done()

			var (
				hasDsa           bool
				hasEcdsa         bool
				hasEd25519       bool
				hasRsa           bool
				userHasStrongRsa bool
			)

			for _, key := range m.Keys {
				atomic.AddUint32(&kt.keySize, 1)
				switch {
				case strings.Contains(key, "ssh-dsa"):
					hasDsa = true
					atomic.AddUint32(&kt.keyDsaSize, 1)
					atomic.AddUint32(&kt.weakKeySize, 1)
				case strings.Contains(key, "ssh-eddsa"):
					hasEcdsa = true
					atomic.AddUint32(&kt.keyEcdsaSize, 1)
					atomic.AddUint32(&kt.weakKeySize, 1)
				case strings.Contains(key, "ssh-ed25519"):
					hasEd25519 = true
					atomic.AddUint32(&kt.keyEd25519Size, 1)
					atomic.AddUint32(&kt.strongKeySize, 1)
				case strings.Contains(key, "ssh-rsa"):
					hasRsa = true
					atomic.AddUint32(&kt.keyRsaSize, 1)
					if isRsaStrong(key) {
						userHasStrongRsa = true
						atomic.AddUint32(&kt.keyStrongRsaSize, 1)
						atomic.AddUint32(&kt.strongKeySize, 1)
					} else {
						userHasStrongRsa = false
						atomic.AddUint32(&kt.keyWeakRsaSize, 1)
						atomic.AddUint32(&kt.weakKeySize, 1)
					}
				}
			}

			switch {
			case hasDsa:
				atomic.AddUint32(&kt.userDsaSize, 1)
				atomic.AddUint32(&kt.userWithWeakKeySize, 1)
				kt.userWithWeakKey = append(kt.userWithWeakKey, m)
			case hasEcdsa:
				atomic.AddUint32(&kt.userEcdsaSize, 1)
				atomic.AddUint32(&kt.userWithWeakKeySize, 1)
				kt.userWithWeakKey = append(kt.userWithWeakKey, m)
			case hasEd25519:
				atomic.AddUint32(&kt.userEd25519Size, 1)
				atomic.AddUint32(&kt.userWithStrongKeySize, 1)
			case hasRsa:
				atomic.AddUint32(&kt.userRsaSize, 1)
				if userHasStrongRsa {
					atomic.AddUint32(&kt.userWithStrongKeySize, 1)
				} else {
					atomic.AddUint32(&kt.userWithWeakKeySize, 1)
					kt.userWithWeakKey = append(kt.userWithWeakKey, m)
				}
			}

			if len(m.Keys) == 0 {
				atomic.AddUint32(&kt.userWithoutKeySize, 1)
				if showUsers == "without" || showUsers == "all" {
					log.Printf("retrieved keys for [%s]: %s", m.Login, m.Keys)
				}
			}

			if len(m.Keys) > 0 {
				atomic.AddUint32(&kt.userWithKeySize, 1)
				if showUsers == "with" || showUsers == "all" {
					log.Printf("retrieved keys for [%s]: %s", m.Login, m.Keys)
				}
			}

			if len(m.Keys) > 1 {
				atomic.AddUint32(&kt.userWithMultipleKeySize, 1)
				if showUsers == "multiple" || showUsers == "all" {
					log.Printf("retrieved keys for [%s]: %s", m.Login, m.Keys)
				}
			}
		}()
	}

	wg.Wait()

	return kt
}

func isRsaStrong(key string) bool {
	r := regexp.MustCompile(`(ssh-rsa) (.*)`)
	keyArray := r.FindStringSubmatch(key)

	return len(keyArray[2]) >= 372
}

func getKeys(ms []member) []member {
	client, err := gh.RESTClient(nil)
	if err != nil {
		log.Fatalf("failed to create REST client: %v", err)
	}

	var wg sync.WaitGroup

	for i := 0; i < len(ms); i++ {
		wg.Add(1)
		go func(i int) {
			m := &ms[i]
			defer wg.Done()

			var KeysT []struct {
				Id  int    `json:"id"`
				Key string `json:"key"`
			}

			err = client.Get(fmt.Sprintf("users/%s/keys", m.Login), &KeysT)
			if err != nil {
				log.Fatalf("failed to get keys for [%s]: %v", m.Login, err)
			}

			for _, key := range KeysT {
				m.Keys = append(m.Keys, key.Key)
			}
		}(i)
	}

	wg.Wait()

	return ms
}

func getMembers(gitHubOrg string) []member {
	client, err := gh.RESTClient(nil)
	if err != nil {
		log.Fatalf("failed to create REST client: %v", err)
	}

	p := 1
	var members []member

	for {
		var ms []member

		err = client.Get(fmt.Sprintf("orgs/%s/members?page=%d&per_page=100", gitHubOrg, p), &ms)
		if err != nil {
			log.Fatalf("failed to get members for [%s]: %v", gitHubOrg, err)
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

// NewCommand is the root command of gh-vanity.
func NewCommand() *cobra.Command {
	var (
		gitHubOrg string
		showUsers string
	)

	var command = &cobra.Command{
		Use:               "gh-audit-org-keys",
		Short:             "Audit GitHub organization keys.",
		DisableAutoGenTag: true,
		Version:           Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			members := getMembers(gitHubOrg)
			keys := getKeys(members)
			table := generateKeyTable(showUsers, keys)
			printReport(table)

			return nil
		},
	}

	command.Flags().StringVarP(&gitHubOrg, "organization", "o", "", "GitHub organization provided to inspect.")
	command.Flags().StringVarP(&showUsers, "show-users", "s", "", "Display users with filter (`all`, `with`, `without`, `multiple`).")

	err := command.MarkFlagRequired("organization")
	if err != nil {
		return nil
	}

	return command
}

func main() {
	command := NewCommand()

	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
