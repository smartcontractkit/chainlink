package src

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type node struct {
	url       *url.URL
	remoteURL *url.URL
	login     string
	password  string
}

func (n node) IsTerminal() bool {
	return false
}

func (n node) PasswordPrompt(p string) string {
	return n.password
}

func (n node) Prompt(p string) string {
	return n.login
}

func writeNodesList(path string, nodes []*node) error {
	fmt.Println("Writing nodes list to", path)
	var lines []string
	for _, n := range nodes {
		lines = append(lines, fmt.Sprintf("%s %s %s", n.url.String(), n.login, n.password))
	}

	return writeLines(lines, path)
}

func mustReadNodesList(path string) []*node {
	fmt.Println("Reading nodes list from", path)
	nodesList, err := readLines(path)
	helpers.PanicErr(err)

	var nodes []*node
	var hasBoot bool
	for _, r := range nodesList {
		rr := strings.TrimSpace(r)
		if len(rr) == 0 {
			continue
		}
		s := strings.Split(rr, " ")
		if len(s) != 4 {
			helpers.PanicErr(errors.New("wrong nodes list format"))
		}
		if strings.Contains(s[0], "boot") && hasBoot {
			helpers.PanicErr(errors.New("the single boot node must come first"))
		}
		hasBoot = true
		url, err := url.Parse(s[0])
		remoteURL, err := url.Parse(s[1])
		helpers.PanicErr(err)
		nodes = append(nodes, &node{
			url:       url,
			remoteURL: remoteURL,
			login:     s[2],
			password:  s[3],
		})
	}
	return nodes
}
