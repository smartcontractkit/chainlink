package src

import (
	"net/url"
	"sort"
)

type NodeWthCreds struct {
	URL            *url.URL
	RemoteURL      *url.URL
	ServiceName    string
	DeploymentName string
	Login          string
	Password       string
}

func (n NodeWthCreds) IsTerminal() bool {
	return false
}

func (n NodeWthCreds) PasswordPrompt(p string) string {
	return n.Password
}

func (n NodeWthCreds) Prompt(p string) string {
	return n.Login
}

// clNodesWithCredsToNodes converts CLNodeCredentials to a slice of nodes.
func clNodesWithCredsToNodes(clNodesWithCreds []CLNodeCredentials) []*NodeWthCreds {
	nodes := []*NodeWthCreds{}
	for _, cl := range clNodesWithCreds {
		n := NodeWthCreds{
			URL:            cl.URL,
			RemoteURL:      cl.URL,
			ServiceName:    cl.ServiceName,
			DeploymentName: cl.DeploymentName,
			Password:       cl.Password,
			Login:          cl.Username,
		}
		nodes = append(nodes, &n)
	}

	// Sort nodes by URL
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].URL.String() < nodes[j].URL.String()
	})
	return nodes
}
