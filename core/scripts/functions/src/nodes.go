package src

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type node struct {
	url      *url.URL
	login    string
	password string
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

func mustReadNodesList(path string) []*node {
	fmt.Println("Reading nodes list from", path)

	nodesList, err := readLines(path)
	if err != nil {
		helpers.PanicErr(err)
	}
	var nodes []*node
	var hasBoot bool
	for _, r := range nodesList {
		rr := strings.TrimSpace(r)
		if len(rr) == 0 {
			continue
		}
		s := strings.Split(rr, " ")
		if len(s) != 3 {
			helpers.PanicErr(errors.New("wrong nodes list format"))
		}
		if !strings.HasPrefix(s[0], "https://") {
			s[0] = "https://" + s[0]
		}
		if strings.Contains(s[0], "boot") && hasBoot {
			helpers.PanicErr(errors.New("the single boot node must come first"))
		}
		hasBoot = true
		url, err := url.Parse(s[0])
		helpers.PanicErr(err)
		nodes = append(nodes, &node{
			url:      url,
			login:    s[1],
			password: s[2],
		})
	}
	return nodes
}
